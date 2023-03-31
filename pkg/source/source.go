// Package source represent source
package source

import (
	"context"
	"teredix/pkg"
	"teredix/pkg/config"
	"teredix/pkg/source/scanner"

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

		if s.Type == pkg.SourceTypeGitHubRepository {
			awsS3 := scanner.NewAWSS3(sourceKey, s.Configuration["access_key"], s.Configuration["secret_key"], s.Configuration["session_token"], s.Configuration["zone"])
			finalSources = append(finalSources, Source{
				Name:    sourceKey,
				Scanner: awsS3,
			})
		}
	}
	return finalSources
}
