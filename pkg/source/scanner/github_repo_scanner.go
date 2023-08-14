// Package scanner scans targets
package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/util"

	"github.com/google/go-github/v50/github"
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
		resourceChannel <- r.mapToResource(repo)
	}

	return nil
}

func (r *GitHubRepositoryScanner) mapToResource(repo *github.Repository) resource.Resource {
	var repoMeta []resource.MetaData

	if util.IsFieldExistsInConfig(fieldCompany, r.fields) && repo.GetOwner() != nil && repo.GetOwner().GetCompany() != "" {
		repoMeta = append(repoMeta, resource.MetaData{Key: fieldCompany, Value: repo.GetOwner().GetCompany()})
	}

	if util.IsFieldExistsInConfig(fieldLanguage, r.fields) && repo.GetLanguage() != "" {
		repoMeta = append(repoMeta, resource.MetaData{Key: fieldLanguage, Value: repo.GetLanguage()})
	}

	if util.IsFieldExistsInConfig(fieldHomepage, r.fields) && repo.GetHomepage() != "" {
		repoMeta = append(repoMeta, resource.MetaData{Key: fieldHomepage, Value: repo.GetHomepage()})
	}

	if util.IsFieldExistsInConfig(fieldOrg, r.fields) && repo.GetOrganization() != nil && repo.GetOrganization().GetName() != "" {
		repoMeta = append(repoMeta, resource.MetaData{Key: fieldOrg, Value: repo.GetOrganization().GetName()})
	}

	if util.IsFieldExistsInConfig(fieldStars, r.fields) {
		repoMeta = append(repoMeta, resource.MetaData{Key: fieldStars, Value: strconv.Itoa(repo.GetStargazersCount())})
	}

	if util.IsFieldExistsInConfig(fieldGitURL, r.fields) && repo.GetGitURL() != "" {
		repoMeta = append(repoMeta, resource.MetaData{Key: fieldGitURL, Value: repo.GetGitURL()})
	}

	if util.IsFieldExistsInConfig(fieldOwnerName, r.fields) && repo.GetOwner() != nil && repo.GetOwner().GetName() != "" {
		repoMeta = append(repoMeta, resource.MetaData{Key: fieldOwnerName, Value: repo.GetOwner().GetName()})
	}

	if util.IsFieldExistsInConfig(fieldOwnerLogin, r.fields) && repo.GetOwner() != nil && repo.GetOwner().GetLogin() != "" {
		repoMeta = append(repoMeta, resource.MetaData{Key: fieldOwnerLogin, Value: repo.GetOwner().GetLogin()})
	}

	if util.IsFieldExistsInConfig(fieldTopics, r.fields) && len(repo.Topics) > 0 {
		topics, err := json.Marshal(repo.Topics)
		if err == nil {
			repoMeta = append(repoMeta, resource.MetaData{Key: fieldTopics, Value: string(topics)})
		}
	}

	re := resource.Resource{
		Kind:       pkg.ResourceKindGitHubRepository,
		UUID:       util.GenerateUUID(),
		Name:       repo.GetFullName(),
		ExternalID: repo.GetFullName(),
		MetaData:   repoMeta,
	}
	return re
}
