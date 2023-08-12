package source

import (
	"context"
	"testing"

	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/source/scanner"

	configv2 "github.com/aws/aws-sdk-go-v2/config"
	credentialsv2 "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"

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
			"source3": {
				Type: pkg.SourceTypeAWSECR,
				Configuration: map[string]string{
					"access_key":    "xxx",
					"secret_key":    "xxx",
					"session_token": "xxx",
					"region":        "me-south-1",
					"account_id":    "xxxx",
				},
			},
		},
	}

	sources := BuildSources(appConfig)

	fsScanner := scanner.NewFsScanner("source1", "/path/to/directory", []string{"rootDirectory"})

	gc := scanner.NewGitHubRepositoryClient(github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "token"},
	))))

	gh := scanner.NewGitHubRepositoryScanner("source2", gc, "user_or_org")

	awsConfig, _ := configv2.LoadDefaultConfig(context.TODO())
	awsCredentials := credentialsv2.NewStaticCredentialsProvider("xxx", "xxx", "xxx")

	awsConfig.Credentials = awsCredentials
	awsConfig.Region = "me-south-1"

	awsEcr := scanner.NewAWSECR(
		"source3",
		"me-south-1",
		"xxxx",
		ecr.NewFromConfig(awsConfig),
		resourcegroupstaggingapi.NewFromConfig(awsConfig),
	)

	expectedSources := []Source{
		{
			Name:    "source1",
			Scanner: &fsScanner,
		},
		{
			Name:    "source2",
			Scanner: gh,
		},
		{
			Name:    "source3",
			Scanner: awsEcr,
		},
	}

	assert.Equal(t, len(sources), len(expectedSources))
}
