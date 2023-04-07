// Package source represent source
package source

import (
	"context"

	"github.com/shahariaazam/teredix/pkg"
	"github.com/shahariaazam/teredix/pkg/config"
	"github.com/shahariaazam/teredix/pkg/source/scanner"

	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"

	aws2 "github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/ecr"

	ec2v2 "github.com/aws/aws-sdk-go-v2/service/ec2"

	configv2 "github.com/aws/aws-sdk-go-v2/config"
	credentialsv2 "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

// Source represent source configuration
type Source struct {
	Name    string
	Scanner scanner.Scanner
}

// BuildSources build source based on configuration
func BuildSources(appConfig *config.AppConfig) []Source {
	var finalSources []Source
	for sourceKey, s := range appConfig.Sources {
		if s.Type == pkg.SourceTypeFileSystem {
			fs := scanner.NewFsScanner(sourceKey, s.Configuration["root_directory"], map[string]string{})
			finalSources = append(finalSources, Source{
				Name:    sourceKey,
				Scanner: &fs,
			})
		}

		if s.Type == pkg.SourceTypeGitHubRepository {

			ctx := context.Background()
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: s.Configuration["token"]},
			)
			tc := oauth2.NewClient(ctx, ts)
			client := github.NewClient(tc)
			gc := scanner.NewGitHubRepositoryClient(client)

			gh := scanner.NewGitHubRepositoryScanner(sourceKey, gc, s.Configuration["user_or_org"])
			finalSources = append(finalSources, Source{
				Name:    sourceKey,
				Scanner: gh,
			})
		}

		if s.Type == pkg.SourceTypeAWSS3 {
			awsCnf := aws.NewConfig().WithRegion(s.Configuration["region"]).WithCredentials(credentials.NewStaticCredentials(s.Configuration["access_key"], s.Configuration["secret_key"], s.Configuration["session_token"]))
			newSession, _ := session.NewSession(awsCnf)
			s3Client := s3.New(newSession)

			awsS3 := scanner.NewAWSS3(sourceKey, s.Configuration["region"], s3Client)
			finalSources = append(finalSources, Source{
				Name:    sourceKey,
				Scanner: awsS3,
			})
		}

		if s.Type == pkg.SourceTypeAWSRDS {
			awsCnf := aws.NewConfig().WithRegion(s.Configuration["region"]).WithCredentials(credentials.NewStaticCredentials(s.Configuration["access_key"], s.Configuration["secret_key"], s.Configuration["session_token"]))
			newSession, _ := session.NewSession(awsCnf)
			rdsClient := rds.New(newSession)

			awsS3 := scanner.NewAWSRDS(sourceKey, s.Configuration["region"], s.Configuration["account_id"], rdsClient)
			finalSources = append(finalSources, Source{
				Name:    sourceKey,
				Scanner: awsS3,
			})
		}

		if s.Type == pkg.SourceTypeAWSEC2 {
			finalSources = append(finalSources, Source{
				Name:    sourceKey,
				Scanner: scanner.NewAWSEC2(sourceKey, s.Configuration["region"], s.Configuration["account_id"], ec2v2.NewFromConfig(buildAWSConfig(s))),
			})
		}

		if s.Type == pkg.SourceTypeAWSECR {
			finalSources = append(finalSources, Source{
				Name: sourceKey,
				Scanner: scanner.NewAWSECR(
					sourceKey,
					s.Configuration["region"],
					s.Configuration["account_id"],
					ecr.NewFromConfig(buildAWSConfig(s)),
					resourcegroupstaggingapi.NewFromConfig(buildAWSConfig(s)),
				),
			})
		}
	}
	return finalSources
}

func buildAWSConfig(s config.Source) aws2.Config {
	cfg, _ := configv2.LoadDefaultConfig(context.TODO())
	awsCredentials := credentialsv2.NewStaticCredentialsProvider(s.Configuration["access_key"], s.Configuration["secret_key"], s.Configuration["session_token"])

	cfg.Credentials = awsCredentials
	cfg.Region = s.Configuration["region"]
	return cfg
}
