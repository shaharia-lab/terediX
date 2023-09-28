// Package scanner scans targets
package scanner

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/metrics"
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
	storage    storage.Storage
	logger     *logrus.Logger
	schedule   string
	metrics    *metrics.Collector
}

// GetName return source name
func (a *AWSRDS) GetName() string {
	return a.SourceName
}

// GetSchedule return schedule
func (a *AWSRDS) GetSchedule() string {
	return a.schedule
}

// Setup AWS S3 source
func (a *AWSRDS) Setup(name string, cfg config.Source, dependencies *Dependencies) error {
	a.SourceName = name
	a.storage = dependencies.GetStorage()
	a.logger = dependencies.GetLogger()
	a.schedule = cfg.Schedule
	a.Region = cfg.Configuration["region"]
	a.AccountID = cfg.Configuration["accountID"]
	a.RdsClient = rds.NewFromConfig(buildAWSConfig(cfg))
	a.Fields = cfg.Fields
	a.metrics = dependencies.GetMetrics()

	a.logger.WithFields(logrus.Fields{
		"scanner_name": a.SourceName,
		"scanner_kind": a.GetKind(),
	}).Info("Scanner has been setup")

	return nil
}

// GetKind return resource kind
func (a *AWSRDS) GetKind() string {
	return pkg.ResourceKindAWSRDS
}

// Scan discover resource and send to resource channel
func (a *AWSRDS) Scan(resourceChannel chan resource.Resource) error {
	nextResourceVersion, err := a.storage.GetNextVersionForResource(a.SourceName, pkg.ResourceKindAWSRDS)
	if err != nil {
		a.logger.WithFields(logrus.Fields{
			"scanner_name": a.SourceName,
			"scanner_kind": a.GetKind(),
		}).WithError(err).Error("Unable to get next version for resource")

		return fmt.Errorf("unable to get next version for resource: %w", err)
	}

	totalResourceDiscovered := 0

	rdsInstances, err := a.listAllRDSInstances(a.RdsClient)

	for _, rdsInstance := range rdsInstances {
		instanceID := aws.StringValue(rdsInstance.DBInstanceIdentifier)

		if err != nil {
			a.logger.WithFields(logrus.Fields{
				"scanner_name": a.SourceName,
				"scanner_kind": a.GetKind(),
			}).WithError(err).Error("Unable to get tags for RDS instance")

			return fmt.Errorf("failed to get tags for RDS instance %s. error: %w", instanceID, err)
		}

		r := resource.NewResource(pkg.ResourceKindAWSRDS, instanceID, instanceID, a.SourceName, nextResourceVersion)
		r.AddMetaData(a.getMetaData(rdsInstance))

		resourceChannel <- r

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

func (a *AWSRDS) listAllRDSInstances(client RdsClient) ([]types.DBInstance, error) {
	var allInstances []types.DBInstance

	// Define Pagination logic
	paginator := rds.NewDescribeDBInstancesPaginator(client, &rds.DescribeDBInstancesInput{})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			a.logger.WithFields(logrus.Fields{
				"scanner_name": a.SourceName,
				"scanner_kind": a.GetKind(),
			}).WithError(err).Error("Unable to make api call to aws rds endpoint")

			return nil, fmt.Errorf("unable to make api call to aws rds endpoint: %w", err)
		}

		allInstances = append(allInstances, output.DBInstances...)
	}

	return allInstances, nil
}
