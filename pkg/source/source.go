package source

import (
	"teredix/pkg/config"
	"teredix/pkg/source/scanner"
)

type Source struct {
	Name    string
	Scanner scanner.Scanner
}

func BuildSources(appConfig *config.AppConfig) []Source {
	var finalSources []Source
	for sourceKey, s := range appConfig.Sources {
		if s.Type == "file_system" {
			fs := scanner.NewFsScanner("fs-scanner_1", s.Configuration["root_directory"], map[string]string{})
			finalSources = append(finalSources, Source{
				Name:    sourceKey,
				Scanner: &fs,
			})
		}
	}
	return finalSources
}
