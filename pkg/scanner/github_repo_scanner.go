// Package scanner scans targets
package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-co-op/gocron"
	"github.com/google/go-github/v50/github"
	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/storage"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const (
	fieldCompany    = "company"
	fieldHomepage   = "homepage"
	fieldLanguage   = "language"
	fieldOrg        = "organization"
	fieldStars      = "stars"
	fieldGitURL     = "git_url"
	fieldOwnerName  = "owner_name"
	fieldOwnerLogin = "owner_login"
	fieldTopics     = "topics"
)

// GitHubClient present interface to build GitHub client
type GitHubClient interface {
	ListRepositories(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, error)
}

// GitHubRepositoryClient GitHub repository client
type GitHubRepositoryClient struct {
	client *github.Client
}

// NewGitHubRepositoryClient construct new GitHub repository client
func NewGitHubRepositoryClient(client *github.Client) *GitHubRepositoryClient {
	return &GitHubRepositoryClient{client: client}
}

// ListRepositories provide list of repositories from GitHub
func (c *GitHubRepositoryClient) ListRepositories(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, error) {
	var rs []*github.Repository

	for {
		repos, resp, err := c.client.Repositories.List(ctx, user, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list repositories for user %s: %w", user, err)
		}

		for _, repo := range repos {
			rs = append(rs, repo)
		}

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return rs, nil
}

// GitHubRepositoryScanner GitHub repository scanner
type GitHubRepositoryScanner struct {
	ghClient  GitHubClient
	user      string
	name      string
	fields    []string
	scheduler *gocron.Scheduler
	storage   storage.Storage
	logger    *logrus.Logger
}

// NewGitHubRepositoryScanner construct a new GitHub repository scanner
func NewGitHubRepositoryScanner(name string, ghClient GitHubClient, user string, fields []string) *GitHubRepositoryScanner {
	return &GitHubRepositoryScanner{ghClient: ghClient, user: user, name: name, fields: fields}
}

func (r *GitHubRepositoryScanner) setGitHubClient(ghClient GitHubClient) {
	r.ghClient = ghClient
}

// GetKind return resource kind
func (r *GitHubRepositoryScanner) GetKind() string {
	return pkg.ResourceKindGitHubRepository
}

// Build GitHub repository scanner
func (r *GitHubRepositoryScanner) Build(sourceKey string, cfg config.Source, storage storage.Storage, scheduler *gocron.Scheduler, logger *logrus.Logger) Scanner {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Configuration["token"]},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	r.name = sourceKey
	r.ghClient = NewGitHubRepositoryClient(client)
	r.user = cfg.Configuration["user"]
	r.fields = cfg.Fields
	r.scheduler = scheduler
	r.storage = storage
	r.logger = logger

	return r
}

// Scan scans GitHub to get the list of repositories as resources
func (r *GitHubRepositoryScanner) Scan(resourceChannel chan resource.Resource) error {
	nextVersion, err := r.storage.GetNextVersionForResource(r.name, pkg.ResourceKindGitHubRepository)
	if err != nil {
		return err
	}

	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	repos, err := r.ghClient.ListRepositories(context.Background(), r.user, opt)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		res := resource.NewResource(pkg.ResourceKindGitHubRepository, repo.GetFullName(), repo.GetFullName(), r.name, nextVersion)
		res.AddMetaData(r.getMetaData(repo))
		resourceChannel <- res
	}

	return nil
}

func (r *GitHubRepositoryScanner) getMetaData(repo *github.Repository) map[string]string {
	mappings := map[string]func() string{
		fieldCompany: func() string {
			if repo.GetOwner() != nil {
				return repo.GetOwner().GetCompany()
			}
			return ""
		},
		fieldLanguage: repo.GetLanguage,
		fieldHomepage: repo.GetHomepage,
		fieldOrg: func() string {
			if repo.GetOrganization() != nil {
				return repo.GetOrganization().GetName()
			}
			return ""
		},
		fieldStars:  func() string { return strconv.Itoa(repo.GetStargazersCount()) },
		fieldGitURL: repo.GetGitURL,
		fieldOwnerName: func() string {
			if repo.GetOwner() != nil {
				return repo.GetOwner().GetName()
			}
			return ""
		},
		fieldOwnerLogin: func() string {
			if repo.GetOwner() != nil {
				return repo.GetOwner().GetLogin()
			}
			return ""
		},
		fieldTopics: func() string {
			topics, err := json.Marshal(repo.Topics)
			if err == nil {
				return string(topics)
			}

			return ""
		},
	}

	fm := NewFieldMapper(mappings, func() []ResourceTag {
		return []ResourceTag{}
	}, r.fields)
	return fm.getResourceMetaData()
}
