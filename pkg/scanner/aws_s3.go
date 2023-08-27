// Package scanner scans targets
package scanner

import (
	"fmt"

	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/util"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/s3"
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
	ListBuckets(listBucketInput *s3.ListBucketsInput) (*s3.ListBucketsOutput, error)
	GetBucketTagging(bucketTaggingInput *s3.GetBucketTaggingInput) (*s3.GetBucketTaggingOutput, error)
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

// Scan discover resource and send to resource channel
func (a *AWSS3) Scan(resourceChannel chan resource.Resource) error {
	result, err := a.S3Client.ListBuckets(nil)
	if err != nil {
		return fmt.Errorf("failed to list buckets. error: %w", err)
	}

	for _, bucket := range result.Buckets {
		resourceChannel <- resource.Resource{
			Kind:        pkg.ResourceKindAWSS3,
			UUID:        util.GenerateUUID(),
			Name:        aws.StringValue(bucket.Name),
			ExternalID:  aws.StringValue(bucket.Name),
			RelatedWith: nil,
			MetaData:    a.getMetaData(bucket),
		}
	}

	return nil
}

func (a *AWSS3) getMetaData(bucket *s3.Bucket) []resource.MetaData {
	mappings := map[string]func() string{
		s3fieldBucketName: func() string { return aws.StringValue(bucket.Name) },
		s3fieldARN: func() string {
			return fmt.Sprintf(s3ARNFormat, aws.StringValue(bucket.Name))
		},
		s3fieldRegion: func() string { return a.Region },
	}

	getTags := func() []ResourceTag {
		var tt []ResourceTag

		if util.IsFieldExistsInConfig(s3fieldTags, a.Fields) == false {
			return tt
		}

		tagResult, _ := a.S3Client.GetBucketTagging(&s3.GetBucketTaggingInput{
			Bucket: aws.String(aws.StringValue(bucket.Name)),
		})

		for _, tag := range tagResult.TagSet {
			tt = append(tt, ResourceTag{
				Key:   aws.StringValue(tag.Key),
				Value: aws.StringValue(tag.Value),
			})
		}

		return tt
	}

	return NewFieldMapper(mappings, getTags, a.Fields).getResourceMetaData()
}
