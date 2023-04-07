// Package scanner scans targets
package scanner

import (
	"context"
	"fmt"

	"github.com/shahariaazam/teredix/pkg"
	"github.com/shahariaazam/teredix/pkg/resource"
	"github.com/shahariaazam/teredix/pkg/util"

	"github.com/google/go-github/v50/github"
)

// GitHubClient present interface to build GitHub client
type GitHubClient interface {
	ListRepositories(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, error)
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
}

// NewGitHubRepositoryScanner construct a new GitHub repository scanner
func NewGitHubRepositoryScanner(name string, ghClient GitHubClient, user string) *GitHubRepositoryScanner {
	return &GitHubRepositoryScanner{ghClient: ghClient, user: user, name: name}
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
	repoMeta := []resource.MetaData{
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
	}

	for _, t := range repo.Topics {
		repoMeta = append(repoMeta, resource.MetaData{
			Key:   "GitHub-Repo-Topic",
			Value: t,
		})
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
