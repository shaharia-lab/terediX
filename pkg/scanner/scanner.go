// Package scanner scans targets
package scanner

import (
	"fmt"
	"strings"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/go-github/v50/github"
	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/util"
	"golang.org/x/oauth2"

	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
)

const (
	fieldTags = "tags"
)

// Scanner interface to build different scanner
type Scanner interface {
	Scan(resourceChannel chan resource.Resource, nextResourceVersion int) error
	GetKind() string
}

// MetaDataMapper map the fields
type MetaDataMapper struct {
	field string
	value func() string
}

// ResourceTag represents a tag on any resource.
type ResourceTag struct {
	Key   string
	Value string
}

// RunScannerForTests initiates a scan using the provided scanner and collects
// the resources it discovers into a slice. This function is specifically
// designed to help with testing, allowing you to run a scanner and easily
// gather its results for verification.
func RunScannerForTests(scanner Scanner) []resource.Resource {
	resourceChannel := make(chan resource.Resource)

	var res []resource.Resource

	go func() {
		scanner.Scan(resourceChannel, 1)
		close(resourceChannel)
	}()

	for r := range resourceChannel {
		res = append(res, r)
	}
	return res
}

func safeDereference(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func stringValueOrDefault(s string) string {
	if s != "" {
		return s
	}
	return ""
}

// FieldMapper is a structure that helps in mapping various fields
// and tags to resource.MetaData structures.
type FieldMapper struct {
	mappings map[string]func() string // Map of field names to their corresponding value functions.
	tags     func() []ResourceTag     // Function that retrieves a list of tags.
	fields   []string                 // List of fields to consider during the mapping.
}

// NewFieldMapper initializes and returns a new instance of FieldMapper.
func NewFieldMapper(mappings map[string]func() string, tags func() []ResourceTag, fields []string) *FieldMapper {
	return &FieldMapper{
		mappings: mappings,
		tags:     tags,
		fields:   fields,
	}
}

// getResourceMetaData constructs and returns a list of resource.MetaData based on
// the FieldMapper's mappings and tags. Only fields specified in the FieldMapper's
// 'fields' slice or having the "tag_" prefix are considered.
//
// For each field in mappings, the associated function is called to retrieve its value.
// Additionally, if tags are specified in the configuration, they are appended with
// the "tag_" prefix and included in the final resource.MetaData list.
func (f *FieldMapper) getResourceMetaData() map[string]string {
	md := make(map[string]string)

	var fieldMapper []MetaDataMapper
	for field, fn := range f.mappings {
		fieldMapper = append(fieldMapper, MetaDataMapper{field: field, value: fn})
	}

	if util.IsFieldExistsInConfig(fieldTags, f.fields) {
		for _, tag := range f.tags() {
			fieldMapper = append(fieldMapper, MetaDataMapper{
				field: fmt.Sprintf("tag_%s", tag.Key),
				value: func() string { return stringValueOrDefault(tag.Value) },
			})
		}
	}

	for _, mapper := range fieldMapper {
		if util.IsFieldExistsInConfig(mapper.field, f.fields) || strings.Contains(mapper.field, "tag_") {
			val := mapper.value()
			if val != "" && mapper.field != "" {
				md[mapper.field] = val
			}
		}
	}

	return md
}

// Source represent source configuration
type Source struct {
	Name    string
	Scanner Scanner
	Kind    string
}

// BuildSources build source based on configuration
func BuildSources(appConfig *config.AppConfig) []Source {
	var finalSources []Source
	for sourceKey, s := range appConfig.Sources {
		if s.Type == pkg.SourceTypeFileSystem {
			fs := NewFsScanner(sourceKey, s.Configuration["root_directory"], s.Fields)
			finalSources = append(finalSources, Source{
				Name:    sourceKey,
				Scanner: fs,
			})
		}

		if s.Type == pkg.SourceTypeGitHubRepository {

			ctx := context.Background()
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: s.Configuration["token"]},
			)
			tc := oauth2.NewClient(ctx, ts)
			client := github.NewClient(tc)
			gc := NewGitHubRepositoryClient(client)

			gh := NewGitHubRepositoryScanner(sourceKey, gc, s.Configuration["user_or_org"], s.Fields)
			finalSources = append(finalSources, Source{
				Name:    sourceKey,
				Scanner: gh,
			})
		}

		if s.Type == pkg.SourceTypeAWSS3 {
			s3Client := s3.NewFromConfig(buildAWSConfig(s))

			awsS3 := NewAWSS3(sourceKey, s.Configuration["region"], s3Client, s.Fields)
			finalSources = append(finalSources, Source{
				Name:    sourceKey,
				Scanner: awsS3,
			})
		}

		if s.Type == pkg.SourceTypeAWSRDS {
			rdsClient := rds.NewFromConfig(buildAWSConfig(s))

			awsS3 := NewAWSRDS(sourceKey, s.Configuration["region"], s.Configuration["account_id"], rdsClient, s.Fields)
			finalSources = append(finalSources, Source{
				Name:    sourceKey,
				Scanner: awsS3,
			})
		}

		if s.Type == pkg.SourceTypeAWSEC2 {
			finalSources = append(finalSources, Source{
				Name:    sourceKey,
				Scanner: NewAWSEC2(sourceKey, s.Configuration["region"], s.Configuration["account_id"], ec2.NewFromConfig(buildAWSConfig(s)), s.Fields),
			})
		}

		if s.Type == pkg.SourceTypeAWSECR {
			finalSources = append(finalSources, Source{
				Name: sourceKey,
				Scanner: NewAWSECR(
					sourceKey,
					s.Configuration["region"],
					s.Configuration["account_id"],
					ecr.NewFromConfig(buildAWSConfig(s)),
					resourcegroupstaggingapi.NewFromConfig(buildAWSConfig(s)),
					s.Fields,
				),
			})
		}
	}
	return finalSources
}

// BuildScanners build source based on configuration
func BuildScanners(appConfig *config.AppConfig) []Scanner {
	var scanners []Scanner
	for sourceKey, s := range appConfig.Sources {
		if s.Type == pkg.SourceTypeFileSystem {
			fs := NewFsScanner(sourceKey, s.Configuration["root_directory"], s.Fields)
			scanners = append(scanners, fs)
		}

		if s.Type == pkg.SourceTypeGitHubRepository {

			ctx := context.Background()
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: s.Configuration["token"]},
			)
			tc := oauth2.NewClient(ctx, ts)
			client := github.NewClient(tc)
			gc := NewGitHubRepositoryClient(client)

			gh := NewGitHubRepositoryScanner(sourceKey, gc, s.Configuration["user_or_org"], s.Fields)
			scanners = append(scanners, gh)
		}

		if s.Type == pkg.SourceTypeAWSS3 {
			s3Client := s3.NewFromConfig(buildAWSConfig(s))

			awsS3 := NewAWSS3(sourceKey, s.Configuration["region"], s3Client, s.Fields)
			scanners = append(scanners, awsS3)
		}

		if s.Type == pkg.SourceTypeAWSRDS {
			rdsClient := rds.NewFromConfig(buildAWSConfig(s))

			awsS3 := NewAWSRDS(sourceKey, s.Configuration["region"], s.Configuration["account_id"], rdsClient, s.Fields)
			scanners = append(scanners, awsS3)
		}

		if s.Type == pkg.SourceTypeAWSEC2 {
			scanners = append(scanners, NewAWSEC2(sourceKey, s.Configuration["region"], s.Configuration["account_id"], ec2.NewFromConfig(buildAWSConfig(s)), s.Fields))
		}

		if s.Type == pkg.SourceTypeAWSECR {
			scanners = append(scanners, NewAWSECR(
				sourceKey,
				s.Configuration["region"],
				s.Configuration["account_id"],
				ecr.NewFromConfig(buildAWSConfig(s)),
				resourcegroupstaggingapi.NewFromConfig(buildAWSConfig(s)),
				s.Fields,
			))
		}
	}
	return scanners
}

func buildAWSConfig(s config.Source) aws.Config {
	cfg, _ := awsConfig.LoadDefaultConfig(context.TODO())
	awsCredentials := credentials.NewStaticCredentialsProvider(s.Configuration["access_key"], s.Configuration["secret_key"], s.Configuration["session_token"])

	cfg.Credentials = awsCredentials
	cfg.Region = s.Configuration["region"]
	return cfg
}
