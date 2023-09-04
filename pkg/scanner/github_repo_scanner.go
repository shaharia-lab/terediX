// Package scanner scans targets
package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/go-github/v50/github"
	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/resource"
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
	ghClient GitHubClient
	user     string
	name     string
	fields   []string
}

// NewGitHubRepositoryScanner construct a new GitHub repository scanner
func NewGitHubRepositoryScanner(name string, ghClient GitHubClient, user string, fields []string) *GitHubRepositoryScanner {
	return &GitHubRepositoryScanner{ghClient: ghClient, user: user, name: name, fields: fields}
}

// Scan scans GitHub to get the list of repositories as resources
func (r *GitHubRepositoryScanner) Scan(resourceChannel chan resource.Resource) error {
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	repos, err := r.ghClient.ListRepositories(context.Background(), r.user, opt)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		res := resource.NewResource(pkg.ResourceKindGitHubRepository, repo.GetFullName(), repo.GetFullName(), r.name, "")
		res.MetaData = r.getMetaData(repo)
		resourceChannel <- res
	}

	return nil
}

func (r *GitHubRepositoryScanner) getMetaData(repo *github.Repository) []resource.MetaData {
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
