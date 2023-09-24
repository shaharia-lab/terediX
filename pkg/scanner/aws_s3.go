// Package scanner scans targets
package scanner

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/go-co-op/gocron"
	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/storage"
	"github.com/sirupsen/logrus"

	"github.com/shaharia-lab/teredix/pkg/util"

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
	Schedule   config.Schedule
	scheduler  *gocron.Scheduler
	storage    storage.Storage
	logger     *logrus.Logger
}

// Build AWS S3 source
func (a *AWSS3) Build(sourceKey string, cfg config.Source, storage storage.Storage, scheduler *gocron.Scheduler, logger *logrus.Logger) Scanner {
	a.SourceName = sourceKey
	a.S3Client = s3.NewFromConfig(BuildAWSConfig(cfg))
	a.Region = cfg.Configuration["region"]
	a.Fields = cfg.Fields
	a.Schedule = cfg.Schedule
	a.storage = storage
	a.scheduler = scheduler
	a.logger = logger

	return a
}

func (a *AWSS3) setS3Client(s3Client AWSS3Client) {
	a.S3Client = s3Client
}

// GetKind return resource kind
func (a *AWSS3) GetKind() string {
	return pkg.ResourceKindAWSS3
}

// Scan discover resource and send to resource channel
func (a *AWSS3) Scan(resourceChannel chan resource.Resource) error {
	nextVersion, err := a.storage.GetNextVersionForResource(a.SourceName, pkg.ResourceKindAWSS3)
	if err != nil {
		return err
	}
	// List all S3 buckets
	output, err := a.S3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("unable to list buckets: %w", err)
	}

	for _, bucket := range output.Buckets {
		res := resource.NewResource(pkg.ResourceKindAWSS3, aws.ToString(bucket.Name), aws.ToString(bucket.Name), a.SourceName, nextVersion)
		res.AddMetaData(a.getMetaData(bucket))
		resourceChannel <- res
	}

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
