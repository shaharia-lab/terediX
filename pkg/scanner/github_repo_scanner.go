// Package scanner scans targets
package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/go-github/v50/github"
	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/metrics"
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
	logger *logrus.Logger
}

// NewGitHubRepositoryClient construct new GitHub repository client
func NewGitHubRepositoryClient(client *github.Client, logger *logrus.Logger) *GitHubRepositoryClient {
	return &GitHubRepositoryClient{client: client, logger: logger}
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
	ghClient GitHubClient
	user     string
	name     string
	fields   []string
	schedule string
	storage  storage.Storage
	logger   *logrus.Logger
	metrics  *metrics.Collector
}

// Setup GitHub repository scanner
func (r *GitHubRepositoryScanner) Setup(name string, cfg config.Source, dependencies *Dependencies) error {
	r.storage = dependencies.GetStorage()
	r.logger = dependencies.GetLogger()
	r.schedule = cfg.Schedule
	r.name = name
	r.user = cfg.Configuration["user_or_org"]
	r.fields = cfg.Fields
	r.metrics = dependencies.GetMetrics()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Configuration["token"]},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	gc := NewGitHubRepositoryClient(client, r.logger)
	r.ghClient = gc

	r.logger.WithFields(logrus.Fields{
		"scanner_name": r.name,
		"scanner_kind": r.GetKind(),
	}).Info("Scanner has been setup")

	return nil
}

func (r *GitHubRepositoryScanner) GetName() string {
	return r.name
}

func (r *GitHubRepositoryScanner) GetSchedule() string {
	return r.schedule
}

// GetKind return resource kind
func (r *GitHubRepositoryScanner) GetKind() string {
	return pkg.ResourceKindGitHubRepository
}

// Scan scans GitHub to get the list of repositories as resources
func (r *GitHubRepositoryScanner) Scan(resourceChannel chan resource.Resource) error {
	nextResourceVersion, err := r.storage.GetNextVersionForResource(r.name, pkg.ResourceKindGitHubRepository)
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"scanner_name": r.name,
			"scanner_kind": r.GetKind(),
		}).WithError(err).Error("Unable to get next version for resource")

		return fmt.Errorf("unable to get next version for resource: %w", err)
	}

	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	repos, err := r.ghClient.ListRepositories(context.Background(), r.user, opt)
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"scanner_name": r.name,
			"scanner_kind": r.GetKind(),
		}).WithError(err).Error("Unable to get repository list from GitHub")

		return err
	}

	totalResourceDiscovered := 0

	for _, repo := range repos {
		res := resource.NewResource(pkg.ResourceKindGitHubRepository, repo.GetFullName(), repo.GetFullName(), r.name, nextResourceVersion)
		res.AddMetaData(r.getMetaData(repo))
		resourceChannel <- res

		totalResourceDiscovered++
	}

	r.logger.WithFields(logrus.Fields{
		"scanner_name":              r.name,
		"scanner_kind":              r.GetKind(),
		"total_resource_discovered": totalResourceDiscovered,
	}).Info("scan completed")

	r.metrics.CollectTotalResourceDiscoveredByScanner(r.name, r.GetKind(), strconv.Itoa(nextResourceVersion), float64(totalResourceDiscovered))
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
