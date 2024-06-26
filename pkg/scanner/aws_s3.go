// Package scanner scans targets
package scanner

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/metrics"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/storage"
	"github.com/shaharia-lab/teredix/pkg/util"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	s3fieldBucketName = "bucketName"
	s3fieldRegion     = "region"
	s3fieldARN        = "arn"
	s3fieldTags       = "tags"

	s3ARNFormat = "arn:aws:s3:::%s"
)

// AWSS3Client build aws client
type AWSS3Client interface {
	ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
	GetBucketTagging(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error)
}

// AWSS3 AWS S3 source
type AWSS3 struct {
	SourceName string
	S3Client   AWSS3Client
	Region     string
	Fields     []string
	schedule   string
	storage    storage.Storage
	logger     *logrus.Logger
	metrics    *metrics.Collector
}

// Setup AWS S3 source
func (a *AWSS3) Setup(name string, cfg config.Source, dependencies *Dependencies) error {
	a.storage = dependencies.GetStorage()
	a.logger = dependencies.GetLogger()
	a.schedule = cfg.Schedule
	a.S3Client = s3.NewFromConfig(buildAWSConfig(cfg))
	a.Region = cfg.Configuration["region"]
	a.Fields = cfg.Fields
	a.SourceName = name
	a.metrics = dependencies.GetMetrics()

	a.logger.WithFields(logrus.Fields{
		"scanner_name": a.SourceName,
		"scanner_kind": a.GetKind(),
	}).Info("Scanner has been setup")

	return nil
}

// GetName return source name
func (a *AWSS3) GetName() string {
	return a.SourceName
}

// GetSchedule return schedule
func (a *AWSS3) GetSchedule() string {
	return a.schedule
}

// GetKind return resource kind
func (a *AWSS3) GetKind() string {
	return pkg.ResourceKindAWSS3
}

// Scan discover resource and send to resource channel
func (a *AWSS3) Scan(resourceChannel chan resource.Resource) error {
	nextResourceVersion, err := a.storage.GetNextVersionForResource(a.SourceName, pkg.ResourceKindAWSS3)
	if err != nil {
		a.logger.WithFields(logrus.Fields{
			"scanner_name": a.SourceName,
			"scanner_kind": a.GetKind(),
		}).WithError(err).Error("Unable to get next version for resource")

		return fmt.Errorf("unable to get next version for resource: %w", err)
	}

	// List all S3 buckets
	output, err := a.S3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		a.logger.WithFields(logrus.Fields{
			"scanner_name": a.SourceName,
			"scanner_kind": a.GetKind(),
		}).WithError(err).Error("Unable to get bucket list from aws s3 endpoint")

		return fmt.Errorf("unable to get bucket list from aws s3 endpoint: %w", err)
	}

	totalResourceDiscovered := 0

	for _, bucket := range output.Buckets {
		res := resource.NewResource(pkg.ResourceKindAWSS3, aws.ToString(bucket.Name), aws.ToString(bucket.Name), a.SourceName, nextResourceVersion)
		res.AddMetaData(a.getMetaData(bucket))
		resourceChannel <- res

		totalResourceDiscovered++
	}

	a.logger.WithFields(logrus.Fields{
		"scanner_name":              a.SourceName,
		"scanner_kind":              a.GetKind(),
		"total_resource_discovered": totalResourceDiscovered,
	}).Info("scan completed")

	a.metrics.CollectTotalResourceDiscoveredByScanner(a.SourceName, a.GetKind(), float64(totalResourceDiscovered))
	return nil
}

func (a *AWSS3) getMetaData(bucket types.Bucket) map[string]string {
	mappings := map[string]func() string{
		s3fieldBucketName: func() string { return aws.ToString(bucket.Name) },
		s3fieldARN: func() string {
			return fmt.Sprintf(s3ARNFormat, aws.ToString(bucket.Name))
		},
		s3fieldRegion: func() string { return a.Region },
	}

	getTags := func() []ResourceTag {
		var tt []ResourceTag

		if util.IsFieldExistsInConfig(s3fieldTags, a.Fields) == false {
			return tt
		}

		tagResult, _ := a.S3Client.GetBucketTagging(context.TODO(), &s3.GetBucketTaggingInput{
			Bucket: aws.String(aws.ToString(bucket.Name)),
		})

		for _, tag := range tagResult.TagSet {
			tt = append(tt, ResourceTag{
				Key:   aws.ToString(tag.Key),
				Value: aws.ToString(tag.Value),
			})
		}

		return tt
	}

	return NewFieldMapper(mappings, getTags, a.Fields).getResourceMetaData()
}
