// Package scanner scans targets
package scanner

import (
	"context"
	"fmt"

	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/storage"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

const (
	perPage = 100

	fieldInstanceID        = "instanceId"
	fieldImageID           = "imageId"
	fieldPrivateDNSName    = "privateDNSName"
	fieldInstanceType      = "instanceType"
	fieldArchitecture      = "architecture"
	fieldInstanceLifecycle = "instanceLifecycle"
	fieldInstanceState     = "instanceState"
	fieldVpcID             = "vpcId"
)

// Ec2Client build aws client
type Ec2Client interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

// AWSEC2 AWS Ec2 source
type AWSEC2 struct {
	SourceName string
	Ec2Client  Ec2Client
	Region     string
	AccountID  string
	Fields     []string
	storage    storage.Storage
	logger     *logrus.Logger
	schedule   string
}

// GetName return source name
func (a *AWSEC2) GetName() string {
	return a.SourceName
}

// GetSchedule return schedule
func (a *AWSEC2) GetSchedule() string {
	return a.schedule
}

// Setup AWS EC2 source
func (a *AWSEC2) Setup(name string, cfg config.Source, dependencies *Dependencies) error {
	a.SourceName = name
	a.Ec2Client = ec2.NewFromConfig(buildAWSConfig(cfg))
	a.Region = cfg.Configuration["region"]
	a.AccountID = cfg.Configuration["account_id"]
	a.Fields = cfg.Fields
	a.schedule = cfg.Schedule
	a.storage = dependencies.GetStorage()
	a.logger = dependencies.GetLogger()

	a.logger.WithFields(logrus.Fields{
		"scanner_name": a.SourceName,
		"scanner_kind": a.GetKind(),
	}).Info("Scanner has been setup")

	return nil
}

// GetKind return resource kind
func (a *AWSEC2) GetKind() string {
	return pkg.ResourceKindAWSEC2
}

// Scan discover resource and send to resource channel
func (a *AWSEC2) Scan(resourceChannel chan resource.Resource) error {
	nextVersion, err := a.storage.GetNextVersionForResource(a.SourceName, pkg.ResourceKindAWSEC2)
	if err != nil {
		a.logger.WithFields(logrus.Fields{
			"scanner_name": a.SourceName,
			"scanner_kind": a.GetKind(),
		}).WithError(err).Error("Unable to get next version for resource")

		return fmt.Errorf("unable to get next version for resource: %w", err)
	}

	totalResourceDiscovered := 0

	pageNum := 0
	nextToken := ""

	for {
		resp, err := a.makeAPICallToAWS(nextToken)
		if err != nil {
			a.logger.WithFields(logrus.Fields{
				"scanner_name": a.SourceName,
				"scanner_kind": a.GetKind(),
			}).WithError(err).Error("Unable to make api call to aws ec2 endpoint")

			return fmt.Errorf("unable to make api call to aws ec2 endpoint: %w", err)
		}

		// Loop through instances and their tags
		for _, reservation := range resp.Reservations {
			for _, instance := range reservation.Instances {
				res := resource.NewResource(pkg.ResourceKindAWSEC2, *instance.InstanceId, *instance.InstanceId, a.SourceName, nextVersion)
				res.AddMetaData(a.getMetaData(instance))
				resourceChannel <- res

				totalResourceDiscovered++
			}
		}

		if resp.NextToken == nil {
			break
		}
		nextToken = *resp.NextToken
		pageNum++
	}

	a.logger.WithFields(logrus.Fields{
		"scanner_name":              a.SourceName,
		"scanner_kind":              a.GetKind(),
		"total_resource_discovered": totalResourceDiscovered,
	}).Info("scan completed")

	return nil
}

// getMetaData build metadata
func (a *AWSEC2) getMetaData(instance types.Instance) map[string]string {
	mappings := map[string]func() string{
		fieldInstanceID:        func() string { return safeDereference(instance.InstanceId) },
		fieldImageID:           func() string { return safeDereference(instance.ImageId) },
		fieldPrivateDNSName:    func() string { return safeDereference(instance.PrivateDnsName) },
		fieldInstanceType:      func() string { return stringValueOrDefault(string(instance.InstanceType)) },
		fieldArchitecture:      func() string { return stringValueOrDefault(string(instance.Architecture)) },
		fieldInstanceLifecycle: func() string { return stringValueOrDefault(string(instance.InstanceLifecycle)) },
		fieldInstanceState:     func() string { return stringValueOrDefault(string(instance.State.Name)) },
		fieldVpcID:             func() string { return safeDereference(instance.VpcId) },
	}
	return NewFieldMapper(mappings, func() []ResourceTag {
		var tags []ResourceTag
		for _, tag := range instance.Tags {
			tags = append(tags, ResourceTag{
				Key:   *tag.Key,
				Value: *tag.Value,
			})
		}
		return tags
	}, a.Fields).getResourceMetaData()
}

func (a *AWSEC2) makeAPICallToAWS(nextToken string) (*ec2.DescribeInstancesOutput, error) {
	// Describe instances for current page
	params := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []string{
					"running",
					"pending",
					"shutting-down",
					"terminated",
					"stopping",
					"stopped",
				},
			},
		},
		MaxResults: aws.Int32(int32(perPage)),
	}

	if nextToken != "" {
		params.NextToken = aws.String(nextToken)
	}

	resp, err := a.Ec2Client.DescribeInstances(context.TODO(), params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
