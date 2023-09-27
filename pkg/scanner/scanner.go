// Package scanner scans targets
package scanner

import (
	"fmt"
	"strings"

	"context"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/metrics"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/scheduler"
	"github.com/shaharia-lab/teredix/pkg/storage"
	"github.com/shaharia-lab/teredix/pkg/util"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
)

const (
	fieldTags = "tags"
)

// Factory is a function that returns a new instance of a Scanner.
type Factory func() Scanner

var scannerFactories = map[string]Factory{
	pkg.ResourceKindAWSEC2:         func() Scanner { return &AWSEC2{} },
	pkg.ResourceKindAWSECR:         func() Scanner { return &AWSECR{} },
	pkg.ResourceKindAWSRDS:         func() Scanner { return &AWSRDS{} },
	pkg.ResourceKindAWSS3:          func() Scanner { return &AWSS3{} },
	pkg.SourceTypeFileSystem:       func() Scanner { return &FsScanner{} },
	pkg.SourceTypeGitHubRepository: func() Scanner { return &GitHubRepositoryScanner{} },
}

func getScanner(sType string) Scanner {
	if factory, found := scannerFactories[sType]; found {
		return factory()
	}
	return nil
}

// BuildScanners build source based on configuration
func BuildScanners(appConfig *config.AppConfig, dependencies *Dependencies) []Scanner {
	var scanners []Scanner
	totalScannerSetup := 0
	for sourceKey, s := range appConfig.Sources {
		scanner := getScanner(s.Type)
		if scanner == nil {
			// Handle unsupported scanner type
			continue
		}

		err := scanner.Setup(sourceKey, s, dependencies)
		if err != nil {
			return nil
		}
		scanners = append(scanners, scanner)

		dependencies.GetMetrics().CollectTotalScannerBuildByName(scanner.GetName(), scanner.GetKind())

		totalScannerSetup++
	}

	dependencies.GetMetrics().CollectTotalScannerBuild(float64(totalScannerSetup))
	return scanners
}

// Dependencies scanner dependencies
type Dependencies struct {
	scheduler scheduler.Scheduler
	storage   storage.Storage
	logger    *logrus.Logger
	metrics   *metrics.Collector
}

// NewScannerDependencies construct new scanner dependencies
func NewScannerDependencies(scheduler scheduler.Scheduler, storage storage.Storage, logger *logrus.Logger, metricsCollector *metrics.Collector) *Dependencies {
	return &Dependencies{scheduler: scheduler, storage: storage, logger: logger, metrics: metricsCollector}
}

// GetScheduler return scheduler
func (d *Dependencies) GetScheduler() scheduler.Scheduler {
	return d.scheduler
}

// GetStorage return storage
func (d *Dependencies) GetStorage() storage.Storage {
	return d.storage
}

// GetLogger return logger
func (d *Dependencies) GetLogger() *logrus.Logger {
	return d.logger
}

// GetMetrics return metrics
func (d *Dependencies) GetMetrics() *metrics.Collector {
	return d.metrics
}

// Scanner interface to build different scanner
type Scanner interface {
	Setup(name string, cfg config.Source, dependencies *Dependencies) error
	Scan(resourceChannel chan resource.Resource) error
	GetName() string
	GetKind() string
	GetSchedule() string
}

// Scanners list of scanners
type Scanners struct {
	Scanners []Scanner
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

func buildAWSConfig(s config.Source) aws.Config {
	cfg, _ := awsConfig.LoadDefaultConfig(context.TODO())
	awsCredentials := credentials.NewStaticCredentialsProvider(s.Configuration["access_key"], s.Configuration["secret_key"], s.Configuration["session_token"])

	cfg.Credentials = awsCredentials
	cfg.Region = s.Configuration["region"]
	return cfg
}
