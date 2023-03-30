// Package scanner scans targets
package scanner

import (
	"context"
	"fmt"
	"teredix/pkg"
	"teredix/pkg/resource"
	"teredix/pkg/util"

	"github.com/google/go-github/v50/github"
)

// GitHubClient present interface to build GitHub client
type GitHubClient interface {
	ListRepositories(ctx context.Context, user string) ([]*github.Repository, error)
}

// GitHubRepositoryClient github repository client
type GitHubRepositoryClient struct {
	client *github.Client
}

// NewGitHubRepositoryClient construct new GitHub repository client
func NewGitHubRepositoryClient(client *github.Client) *GitHubRepositoryClient {
	return &GitHubRepositoryClient{client: client}
}

// ListRepositories provide list of repositories from GitHub
func (c *GitHubRepositoryClient) ListRepositories(ctx context.Context, user string) ([]*github.Repository, error) {
	repos, _, err := c.client.Repositories.List(ctx, user, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories for user %s: %w", user, err)
	}

	return repos, nil
}

// GitHubRepositoryScanner GitHub repository scanner
type GitHubRepositoryScanner struct {
	ghClient GitHubClient
	user     string
	name     string
}

// NewGitHubRepositoryScanner construct a new GitHub repository scanner
func NewGitHubRepositoryScanner(name string, ghClient GitHubClient, user string) *GitHubRepositoryScanner {
	return &GitHubRepositoryScanner{ghClient: ghClient, user: user, name: name}
}

// Scan scans GitHub to get the list of repositories as resources
func (r *GitHubRepositoryScanner) Scan() []resource.Resource {
	var resources []resource.Resource

	repos, err := r.ghClient.ListRepositories(context.Background(), r.user)
	if err != nil {
		return resources
	}

	for _, repo := range repos {
		re := resource.Resource{
			Kind:       pkg.ResourceKindGitHubRepository,
			UUID:       util.GenerateUUID(),
			Name:       repo.GetFullName(),
			ExternalID: repo.GetFullName(),
			MetaData: []resource.MetaData{
				{
					Key:   "GitHub-Repo-Language",
					Value: repo.GetLanguage(),
				},
				{
					Key:   "GitHub-Repo-Stars",
					Value: fmt.Sprintf("%d", repo.GetStargazersCount()),
				},
				{
					Key:   pkg.MetaKeyScannerLabel,
					Value: r.name,
				},
				{
					Key:   "GitHub-Repo-Git-URL",
					Value: repo.GetGitURL(),
				},
				{
					Key:   "GitHub-Repo-Homepage",
					Value: repo.GetHomepage(),
				},
				{
					Key:   "GitHub-Repo-Organization",
					Value: repo.GetOrganization().GetLogin(),
				},
				{
					Key:   "GitHub-Owner",
					Value: repo.GetOwner().GetLogin(),
				},
				{
					Key:   "GitHub-Company",
					Value: repo.GetOwner().GetCompany(),
				},
			},
		}

		resources = append(resources, re)
	}

	return resources
}
