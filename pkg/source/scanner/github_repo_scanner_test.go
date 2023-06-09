package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shahariaazam/teredix/pkg/resource"

	"github.com/google/go-github/v50/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// GitHubClientMock is an autogenerated mock type for the GitHubClient type
type GitHubClientMock struct {
	mock.Mock
}

// ListRepositories provides a mock function with given fields: ctx, user, opts
func (_m *GitHubClientMock) ListRepositories(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, error) {
	ret := _m.Called(ctx, user, opts)

	var r0 []*github.Repository
	if rf, ok := ret.Get(0).(func(context.Context, string, *github.RepositoryListOptions) []*github.Repository); ok {
		r0 = rf(ctx, user, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*github.Repository)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *github.RepositoryListOptions) error); ok {
		r1 = rf(ctx, user, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func TestGitHubRepositoryScanner_Scan(t *testing.T) {
	testCases := []struct {
		name                 string
		user                 string
		ghRepositories       []*github.Repository
		want                 []resource.Resource
		expectedMetaDataKeys []string
	}{
		{
			name: "returns resources",
			user: "testuser",
			ghRepositories: []*github.Repository{
				{
					ID:              github.Int64(123),
					Name:            github.String("testrepo"),
					FullName:        github.String("testuser/testrepo"),
					Language:        github.String("Go"),
					StargazersCount: github.Int(42),
					GitURL:          github.String("https://git_url"),
				},
			},
			want: []resource.Resource{
				{
					Kind:       "GitHubRepository",
					UUID:       "123",
					Name:       "testrepo",
					ExternalID: "testuser/testrepo",
					MetaData: []resource.MetaData{
						{Key: "Language", Value: "Go"},
						{Key: "Stars", Value: "42"},
					},
				},
			},
			expectedMetaDataKeys: []string{
				"GitHub-Repo-Language",
				"GitHub-Repo-Stars",
				"GitHub-Repo-Homepage",
				"GitHub-Repo-Organization",
				"GitHub-Owner",
				"GitHub-Company",
				"GitHub-Repo-Topic",
				"Scanner-Label",
				"GitHub-Repo-Git-URL",
			},
		},
		{
			name:                 "returns empty resource list on error",
			user:                 "testuser",
			ghRepositories:       []*github.Repository{},
			want:                 []resource.Resource{},
			expectedMetaDataKeys: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := new(GitHubClientMock)
			mockClient.On("ListRepositories", mock.Anything, mock.Anything, mock.Anything).Return(tc.ghRepositories, nil)

			resourceChannel := make(chan resource.Resource, len(tc.ghRepositories))
			var res []resource.Resource

			go func() {
				s := NewGitHubRepositoryScanner("test", mockClient, tc.user)
				s.Scan(resourceChannel)
				close(resourceChannel)
			}()

			for r := range resourceChannel {
				res = append(res, r)
			}

			for _, r := range res {
				for _, md := range r.MetaData {
					assert.True(t, containsValue(tc.expectedMetaDataKeys, md.Key))
				}
			}

			assert.Equal(t, len(tc.ghRepositories), len(res))
		})
	}
}

func TestNewGitHubRepositoryClient_ListRepositories_Return_Data(t *testing.T) {
	ctx := context.Background()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return a list of repositories
		repos := []*github.Repository{
			{Name: github.String("repo1")},
			{Name: github.String("repo2")},
			{Name: github.String("repo3")},
		}
		jsonBytes, _ := json.Marshal(repos)
		fmt.Fprintln(w, string(jsonBytes))
	}))
	defer ts.Close()

	client, _ := github.NewEnterpriseClient(ts.URL, "", ts.Client())
	gc := NewGitHubRepositoryClient(client)
	repositories, err := gc.ListRepositories(ctx, "HI", &github.RepositoryListOptions{})

	assert.NoError(t, err)
	assert.Equal(t, 3, len(repositories))
}

func TestNewGitHubRepositoryClient_ListRepositories_Bad_Response_Code(t *testing.T) {
	ctx := context.Background()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	client, _ := github.NewEnterpriseClient(ts.URL, "", ts.Client())
	gc := NewGitHubRepositoryClient(client)
	_, err := gc.ListRepositories(ctx, "HI", &github.RepositoryListOptions{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to list repositories for user HI")
}

func containsValue(values []string, value string) bool {
	for _, v := range values {
		if strings.Contains(v, value) {
			return true
		}
	}
	return false
}
