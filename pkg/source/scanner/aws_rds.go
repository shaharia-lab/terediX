// Package scanner scans targets
package scanner

import (
	"fmt"
	"teredix/pkg"
	"teredix/pkg/resource"
	"teredix/pkg/util"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
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
}

// NewAWSRDS construct AWS S3 source
func NewAWSRDS(sourceName string, region string, accountID string, rdsClient RdsClient) *AWSRDS {
	return &AWSRDS{
		SourceName: sourceName,
		RdsClient:  rdsClient,
		Region:     region,
		AccountID:  accountID,
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

		// Get the tags for the instance
		tagResult, err := a.RdsClient.ListTagsForResource(&rds.ListTagsForResourceInput{
			ResourceName: aws.String(fmt.Sprintf("arn:aws:rds:%s:%s:db:%s", a.Region, a.AccountID, instanceID)),
		})
		if err != nil {
			return fmt.Errorf("failed to get tags for RDS instance %s. error: %w", instanceID, err)
		}

		r := resource.Resource{
			Kind:        pkg.ResourceKindAWSRDS,
			UUID:        util.GenerateUUID(),
			Name:        instanceID,
			ExternalID:  instanceID,
			RelatedWith: nil,
			MetaData: []resource.MetaData{
				{
					Key:   "AWS-RDS-Instance-ID",
					Value: instanceID,
				},
				{
					Key:   "AWS-RDS-Region",
					Value: a.Region,
				},
				{
					Key:   "AWS-ARN",
					Value: fmt.Sprintf("arn:aws:rds:%s:%s:db:%s", a.Region, a.AccountID, instanceID),
				},
			},
		}

		for _, tag := range tagResult.TagList {
			r.MetaData = append(r.MetaData, resource.MetaData{
				Key:   fmt.Sprintf("AWS-RDS-Tag-%s", aws.StringValue(tag.Key)),
				Value: aws.StringValue(tag.Value),
			})
		}

		resourceChannel <- r
	}

	return nil
}
