// Package scanner scans targets
package scanner

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/util"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
)

const (
	rdsFieldInstanceID = "instanceID"
	rdsFieldRegion     = "region"
	rdsFieldARN        = "arn"
	rdsFieldTags       = "tags"

	rdsARNFormat = "arn:aws:rds:%s:%s:db:%s"
)

// RdsClient build aws client
type RdsClient interface {
	DescribeDBInstancesPages(*rds.DescribeDBInstancesInput, func(*rds.DescribeDBInstancesOutput, bool) bool) error
	ListTagsForResource(*rds.ListTagsForResourceInput) (*rds.ListTagsForResourceOutput, error)
}

// AWSRDS AWS S3 source
type AWSRDS struct {
	SourceName string
	RdsClient  RdsClient
	Region     string
	AccountID  string
	Fields     []string
}

// NewAWSRDS construct AWS S3 source
func NewAWSRDS(sourceName string, region string, accountID string, rdsClient RdsClient, fields []string) *AWSRDS {
	return &AWSRDS{
		SourceName: sourceName,
		RdsClient:  rdsClient,
		Region:     region,
		AccountID:  accountID,
		Fields:     fields,
	}
}

// Scan discover resource and send to resource channel
func (a *AWSRDS) Scan(resourceChannel chan resource.Resource) error {
	// Get a list of all RDS instances
	var rdsInstances []*rds.DBInstance
	err := a.RdsClient.DescribeDBInstancesPages(&rds.DescribeDBInstancesInput{}, func(page *rds.DescribeDBInstancesOutput, lastPage bool) bool {
		rdsInstances = append(rdsInstances, page.DBInstances...)
		return !lastPage
	})
	if err != nil {
		return fmt.Errorf("failed to list RDS instances. error: %w", err)
	}

	// Loop through each instance and get its tags
	for _, rdsInstance := range rdsInstances {
		instanceID := aws.StringValue(rdsInstance.DBInstanceIdentifier)

		if err != nil {
			return fmt.Errorf("failed to get tags for RDS instance %s. error: %w", instanceID, err)
		}

		r := resource.Resource{
			Kind:        pkg.ResourceKindAWSRDS,
			UUID:        util.GenerateUUID(),
			Name:        instanceID,
			ExternalID:  instanceID,
			RelatedWith: nil,
			MetaData:    a.getMetaData(rdsInstance),
		}

		resourceChannel <- r
	}

	return nil
}

func (a *AWSRDS) getMetaData(rdsInstance *rds.DBInstance) []resource.MetaData {
	mappings := map[string]func() string{
		rdsFieldInstanceID: func() string { return aws.StringValue(rdsInstance.DBInstanceIdentifier) },
		rdsFieldARN: func() string {
			return fmt.Sprintf(rdsARNFormat, a.Region, a.AccountID, aws.StringValue(rdsInstance.DBInstanceIdentifier))
		},
		rdsFieldRegion: func() string { return a.Region },
	}

	getTags := func() []types.Tag {
		var tt []types.Tag
		for _, tag := range rdsInstance.TagList {
			tt = append(tt, types.Tag{
				Key:   aws.String(aws.StringValue(tag.Key)),
				Value: aws.String(aws.StringValue(tag.Value)),
			})
		}

		return tt
	}

	return NewFieldMapper(mappings, getTags, a.Fields).getResourceMetaData()
}
