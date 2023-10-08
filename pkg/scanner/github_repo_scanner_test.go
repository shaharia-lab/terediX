package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-github/v50/github"
	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/metrics"
	"github.com/shaharia-lab/teredix/pkg/scheduler"
	"github.com/shaharia-lab/teredix/pkg/storage"
	"github.com/sirupsen/logrus"
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
		name                  string
		sourceFields          []string
		ghRepositories        []*github.Repository
		expectedTotalResource int
		expectedTotalMetaData int
		expectedMetaDataKeys  []string
	}{
		{
			name: "returns resources",
			sourceFields: []string{
				fieldCompany,
				fieldHomepage,
				fieldLanguage,
				fieldOrg,
				fieldStars,
				fieldGitURL,
				fieldOwnerName,
				fieldOwnerLogin,
				fieldTopics,
			},
			ghRepositories: []*github.Repository{
				{
					Name: github.String("repo1"),
					Owner: &github.User{
						Login:   github.String("shaharia-lab"),
						Name:    github.String("Shaharia Lab"),
						Company: github.String("Shaharia Lab"),
					},
					GitURL:          github.String("https://github.com/shaharia-lab/teredix"),
					Description:     github.String("This is a test repository"),
					Homepage:        github.String("https://github.com/shaharia-lab/teredix"),
					Language:        github.String("go"),
					Topics:          []string{"teredix", "go", "github", "repository"},
					StargazersCount: github.Int(1),
					Organization: &github.Organization{
						Login: github.String("shaharia-lab"),
						Name:  github.String("Shaharia Lab"),
					},
				},
			},
			expectedTotalResource: 1,
			expectedTotalMetaData: 9,
			expectedMetaDataKeys: []string{
				fieldCompany,
				fieldHomepage,
				fieldLanguage,
				fieldOrg,
				fieldStars,
				fieldGitURL,
				fieldOwnerName,
				fieldOwnerLogin,
				fieldTopics,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := new(GitHubClientMock)
			mockClient.On("ListRepositories", mock.Anything, mock.Anything, mock.Anything).Return(tc.ghRepositories, nil)

			mockStorage := new(storage.Mock)
			mockStorage.On("GetNextVersionForResource", mock.Anything, mock.Anything).Return(1, nil)

			sc := config.Source{
				Type: pkg.SourceTypeGitHubRepository,
				Configuration: map[string]string{
					"token":       "test",
					"user_or_org": "shaharia-lab",
				},
				Fields:   tc.sourceFields,
				Schedule: "",
			}

			gh := GitHubRepositoryScanner{}
			gh.Setup("test", sc, NewScannerDependencies(scheduler.NewGoCron(), mockStorage, &logrus.Logger{}, metrics.NewCollector()))
			gh.ghClient = mockClient

			RunCommonScannerAssertionTest(t, &gh, tc.expectedTotalResource, tc.expectedTotalMetaData, tc.expectedMetaDataKeys)
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
	gc := NewGitHubRepositoryClient(client, &logrus.Logger{})
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
	gc := NewGitHubRepositoryClient(client, &logrus.Logger{})
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
