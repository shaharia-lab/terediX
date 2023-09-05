// Package scanner scans targets
package scanner

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/resource"
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
}

// NewAWSS3 construct AWS S3 source
func NewAWSS3(sourceName string, region string, s3Client AWSS3Client, fields []string) *AWSS3 {
	return &AWSS3{
		SourceName: sourceName,
		S3Client:   s3Client,
		Region:     region,
		Fields:     fields,
	}
}

func (a *AWSS3) GetKind() string {
	return pkg.ResourceKindAWSS3
}

// Scan discover resource and send to resource channel
func (a *AWSS3) Scan(resourceChannel chan resource.Resource, nextResourceVersion int) error {
	// List all S3 buckets
	output, err := a.S3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("unable to list buckets: %w", err)
	}

	for _, bucket := range output.Buckets {
		res := resource.NewResource(pkg.ResourceKindAWSS3, aws.ToString(bucket.Name), aws.ToString(bucket.Name), a.SourceName, 1)
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
