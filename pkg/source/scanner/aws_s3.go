// Package scanner scans targets
package scanner

import (
	"fmt"
	"teredix/pkg"
	"teredix/pkg/resource"
	"teredix/pkg/util"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/s3"
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
}

// NewAWSS3 construct AWS S3 source
func NewAWSS3(sourceName string, region string, s3Client AWSS3Client) *AWSS3 {
	return &AWSS3{
		SourceName: sourceName,
		S3Client:   s3Client,
		Region:     region,
	}
}

// Scan discover resource and send to resource channel
func (a *AWSS3) Scan(resourceChannel chan resource.Resource) error {
	result, err := a.S3Client.ListBuckets(nil)
	if err != nil {
		return fmt.Errorf("failed to list buckets. error: %w", err)
	}

	for _, bucket := range result.Buckets {
		resourceChannel <- a.mapToResource(bucket)
	}

	return nil
}

func (a *AWSS3) mapToResource(bucket *s3.Bucket) resource.Resource {
	res := resource.Resource{
		Kind:        pkg.ResourceKindAWSS3,
		UUID:        util.GenerateUUID(),
		Name:        aws.StringValue(bucket.Name),
		ExternalID:  aws.StringValue(bucket.Name),
		RelatedWith: nil,
		MetaData: []resource.MetaData{
			{
				Key:   "AWS-S3-Bucket-Name",
				Value: aws.StringValue(bucket.Name),
			},
			{
				Key:   pkg.MetaKeyScannerLabel,
				Value: a.SourceName,
			},
			{
				Key:   "AWS-S3-Region",
				Value: a.Region,
			},
			{
				Key:   "AWS-ARN",
				Value: fmt.Sprintf("arn:aws:s3:::%s", aws.StringValue(bucket.Name)),
			},
		},
	}

	bucketName := aws.StringValue(bucket.Name)

	tagResult, _ := a.S3Client.GetBucketTagging(&s3.GetBucketTaggingInput{
		Bucket: aws.String(bucketName),
	})

	for _, tag := range tagResult.TagSet {
		res.MetaData = append(res.MetaData, resource.MetaData{
			Key:   fmt.Sprintf("AWS-S3-Tag-%s", aws.StringValue(tag.Key)),
			Value: aws.StringValue(tag.Value),
		})
	}

	return res
}
