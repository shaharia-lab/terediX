package scanner

import (
	"context"
	"fmt"
	"teredix/pkg"
	"teredix/pkg/resource"
	"teredix/pkg/util"

	"github.com/google/go-github/v50/github"
)

type GitHubClient interface {
	ListRepositories(ctx context.Context, user string) ([]*github.Repository, error)
}

type GitHubRepositoryClient struct {
	client *github.Client
}

func NewGitHubRepositoryClient(client *github.Client) *GitHubRepositoryClient {
	return &GitHubRepositoryClient{client: client}
}

func (c *GitHubRepositoryClient) ListRepositories(ctx context.Context, user string) ([]*github.Repository, error) {
	repos, _, err := c.client.Repositories.List(ctx, user, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories for user %s: %w", user, err)
	}

	return repos, nil
}

type GitHubRepositoryScanner struct {
	ghClient GitHubClient
	user     string
	name     string
}

func NewGitHubRepositoryScanner(name string, ghClient GitHubClient, user string) *GitHubRepositoryScanner {
	return &GitHubRepositoryScanner{ghClient: ghClient, user: user, name: name}
}

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
