package source

import (
	"context"
	"teredix/pkg"
	"teredix/pkg/config"
	"teredix/pkg/source/scanner"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type Source struct {
	Name    string
	Scanner scanner.Scanner
}

func BuildSources(appConfig *config.AppConfig) []Source {
	var finalSources []Source
	for sourceKey, s := range appConfig.Sources {
		if s.Type == pkg.SourceTypeFileSystem {
			fs := scanner.NewFsScanner("fs-scanner_1", s.Configuration["root_directory"], map[string]string{})
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

			gh := scanner.NewGitHubRepositoryScanner(gc, s.Configuration["user_or_org"])
			finalSources = append(finalSources, Source{
				Name:    sourceKey,
				Scanner: gh,
			})
		}
	}
	return finalSources
}
