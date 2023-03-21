package source

import (
	"infrastructure-discovery/pkg/config"
	"infrastructure-discovery/pkg/source/scanner"
)

type Source struct {
	Name    string
	Scanner scanner.Scanner
}

func BuildSources(appConfig *config.AppConfig) []Source {
	var finalSources []Source
	for sourceKey, s := range appConfig.Sources {
		if s.Type == "file_system" {
			fs := scanner.NewFsScanner("fs-scanner_1", s.Configuration["root_directory"], map[string]string{"key1": "value1", "key2": "value2"})
			finalSources = append(finalSources, Source{
				Name:    sourceKey,
				Scanner: &fs,
			})
		}
	}
	return finalSources
}
