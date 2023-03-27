package scanner

import (
	"context"
	"fmt"
	"teredix/pkg/resource"

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
	repos, _, err := c.client.Repositories.List(context.Background(), user, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories for user %s: %w", user, err)
	}

	return repos, nil
}

type GitHubRepositoryScanner struct {
	ghClient GitHubClient
	user     string
}

func NewGitHubRepositoryScanner(ghClient GitHubClient, user string) *GitHubRepositoryScanner {
	return &GitHubRepositoryScanner{ghClient: ghClient, user: user}
}

func (r *GitHubRepositoryScanner) Scan() []resource.Resource {
	var resources []resource.Resource

	repos, err := r.ghClient.ListRepositories(context.Background(), r.user)
	if err != nil {
		return resources
	}

	for _, repo := range repos {
		re := resource.Resource{
			Kind:       "GitHubRepository",
			UUID:       fmt.Sprintf("%v", repo.GetID()),
			Name:       repo.GetName(),
			ExternalID: repo.GetFullName(),
			MetaData: []resource.MetaData{
				{
					Key:   "Language",
					Value: repo.GetLanguage(),
				},
				{
					Key:   "Stars",
					Value: fmt.Sprintf("%d", repo.GetStargazersCount()),
				},
			},
		}

		resources = append(resources, re)
	}

	return resources
}
