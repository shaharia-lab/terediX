package source

import (
	"context"
	"teredix/pkg"
	"teredix/pkg/config"
	"teredix/pkg/source/scanner"
	"testing"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"

	"github.com/stretchr/testify/assert"
)

func TestBuildSources(t *testing.T) {
	appConfig := &config.AppConfig{
		Sources: map[string]config.Source{
			"source1": {
				Type: pkg.SourceTypeFileSystem,
				Configuration: map[string]string{
					"root_directory": "/path/to/directory",
				},
			},
			"source2": {
				Type: pkg.SourceTypeGitHubRepository,
				Configuration: map[string]string{
					"token":       "token",
					"user_or_org": "user_or_org",
				},
			},
		},
	}

	sources := BuildSources(appConfig)

	fsScanner := scanner.NewFsScanner("source1", "/path/to/directory", map[string]string{})

	gc := scanner.NewGitHubRepositoryClient(github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "token"},
	))))

	gh := scanner.NewGitHubRepositoryScanner("source2", gc, "user_or_org")

	expectedSources := []Source{
		{
			Name:    "source1",
			Scanner: &fsScanner,
		},
		{
			Name:    "source2",
			Scanner: gh,
		},
	}

	assert.Equal(t, expectedSources, sources)
}
