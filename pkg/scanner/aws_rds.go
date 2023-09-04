// Package scanner scans targets
package scanner

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	types "github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/resource"

	"github.com/aws/aws-sdk-go/aws"
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
	DescribeDBInstances(context.Context, *rds.DescribeDBInstancesInput, ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error)
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
	rdsInstances, err := listAllRDSInstances(a.RdsClient)

	for _, rdsInstance := range rdsInstances {
		instanceID := aws.StringValue(rdsInstance.DBInstanceIdentifier)

		if err != nil {
			return fmt.Errorf("failed to get tags for RDS instance %s. error: %w", instanceID, err)
		}

		r := resource.NewResource(pkg.ResourceKindAWSRDS, instanceID, instanceID, a.SourceName, "")
		r.AddMetaData(a.getMetaData(rdsInstance))

		resourceChannel <- r
	}

	return nil
}

func (a *AWSRDS) getMetaData(rdsInstance types.DBInstance) map[string]string {
	mappings := map[string]func() string{
		rdsFieldInstanceID: func() string { return aws.StringValue(rdsInstance.DBInstanceIdentifier) },
		rdsFieldARN: func() string {
			return fmt.Sprintf(rdsARNFormat, a.Region, a.AccountID, aws.StringValue(rdsInstance.DBInstanceIdentifier))
		},
		rdsFieldRegion: func() string { return a.Region },
	}

	getTags := func() []ResourceTag {
		var tt []ResourceTag
		for _, tag := range rdsInstance.TagList {
			tt = append(tt, ResourceTag{
				Key:   aws.StringValue(tag.Key),
				Value: aws.StringValue(tag.Value),
			})
		}

		return tt
	}

	return NewFieldMapper(mappings, getTags, a.Fields).getResourceMetaData()
}

func listAllRDSInstances(client RdsClient) ([]types.DBInstance, error) {
	var allInstances []types.DBInstance

	// Define Pagination logic
	paginator := rds.NewDescribeDBInstancesPaginator(client, &rds.DescribeDBInstancesInput{})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}
		allInstances = append(allInstances, output.DBInstances...)
	}

	return allInstances, nil
}
