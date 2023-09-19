// Package scanner scans targets
package scanner

import (
	"context"
	"fmt"
	"strings"

	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/util"

	"github.com/aws/aws-sdk-go-v2/aws"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

const (
	fieldTags = "tags"
)

// Scanner interface to build different scanner
type Scanner interface {
	Build(sourceKey string, source config.Source) Scanner
	Scan(resourceChannel chan resource.Resource, nextResourceVersion int) error
	GetKind() string
}

// Sources registry
type Sources struct {
	Scanners map[string]Scanner
}

// NewSourceRegistry build source registry
func NewSourceRegistry(scanners map[string]Scanner) *Sources {
	return &Sources{Scanners: scanners}
}

// Get source
func (s *Sources) Get(sourceKey string) Scanner {
	return s.Scanners[sourceKey]
}

// Add new source
func (s *Sources) Add(sourceKey string, scanner Scanner) {
	s.Scanners[sourceKey] = scanner
}

// BuildFromAppConfig build source based on configuration
func (s *Sources) BuildFromAppConfig(appConfig config.AppConfig) []Scanner {
	var scanners []Scanner
	for sourceKey, sc := range appConfig.Sources {
		scanners = append(scanners, s.Scanners[sourceKey].Build(sourceKey, sc))
	}
	return scanners
}

// GetScannerRegistries get all scanner registries
func GetScannerRegistries() map[string]Scanner {
	return map[string]Scanner{
		pkg.SourceTypeFileSystem:       &FsScanner{},
		pkg.SourceTypeGitHubRepository: &GitHubRepositoryScanner{},
		pkg.SourceTypeAWSS3:            &AWSS3{},
		pkg.SourceTypeAWSRDS:           &AWSRDS{},
		pkg.SourceTypeAWSEC2:           &AWSEC2{},
		pkg.SourceTypeAWSECR:           &AWSECR{},
	}
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

// BuildAWSConfig build aws config
func BuildAWSConfig(s config.Source) aws.Config {
	cfg, _ := awsConfig.LoadDefaultConfig(context.TODO())
	awsCredentials := credentials.NewStaticCredentialsProvider(s.Configuration["access_key"], s.Configuration["secret_key"], s.Configuration["session_token"])

	cfg.Credentials = awsCredentials
	cfg.Region = s.Configuration["region"]
	return cfg
}
