// Package scanner scans targets
package scanner

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	types "github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/go-co-op/gocron"
	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/storage"
	"github.com/sirupsen/logrus"

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
	Schedule   config.Schedule
	scheduler  *gocron.Scheduler
	storage    storage.Storage
	logger     *logrus.Logger
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

// Build AWS S3 source
func (a *AWSRDS) Build(sourceKey string, cfg config.Source, storage storage.Storage, scheduler *gocron.Scheduler, logger *logrus.Logger) Scanner {
	a.SourceName = sourceKey
	a.RdsClient = rds.NewFromConfig(BuildAWSConfig(cfg))
	a.Region = cfg.Configuration["region"]
	a.AccountID = cfg.Configuration["account_id"]
	a.Fields = cfg.Fields
	a.Schedule = cfg.Schedule
	a.scheduler = scheduler
	a.storage = storage
	a.logger = logger

	return a
}

// GetKind return resource kind
func (a *AWSRDS) GetKind() string {
	return pkg.ResourceKindAWSRDS
}

// Scan discover resource and send to resource channel
func (a *AWSRDS) Scan(resourceChannel chan resource.Resource) error {
	nextVersion, err := a.storage.GetNextVersionForResource(a.SourceName, pkg.ResourceKindAWSEC2)
	if err != nil {
		return err
	}
	rdsInstances, err := listAllRDSInstances(a.RdsClient)

	for _, rdsInstance := range rdsInstances {
		instanceID := aws.StringValue(rdsInstance.DBInstanceIdentifier)

		if err != nil {
			return fmt.Errorf("failed to get tags for RDS instance %s. error: %w", instanceID, err)
		}

		r := resource.NewResource(pkg.ResourceKindAWSRDS, instanceID, instanceID, a.SourceName, nextVersion)
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

func (a *AWSRDS) setRDSClient(rdsClient RdsClient) {
	a.RdsClient = rdsClient
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
