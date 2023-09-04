// Package scanner scans targets
package scanner

import (
	"context"

	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/resource"

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
}

// NewAWSEC2 construct AWS EC2 source
func NewAWSEC2(sourceName string, region string, accountID string, ec2Client Ec2Client, fields []string) *AWSEC2 {
	return &AWSEC2{
		SourceName: sourceName,
		Ec2Client:  ec2Client,
		Region:     region,
		AccountID:  accountID,
		Fields:     fields,
	}
}

// Scan discover resource and send to resource channel
func (a *AWSEC2) Scan(resourceChannel chan resource.Resource) error {
	pageNum := 0
	nextToken := ""

	for {
		resp, err := a.makeAPICallToAWS(nextToken)
		if err != nil {
			return err
		}

		// Loop through instances and their tags
		for _, reservation := range resp.Reservations {
			for _, instance := range reservation.Instances {
				res := resource.NewResource(pkg.ResourceKindAWSEC2, *instance.InstanceId, *instance.InstanceId, a.SourceName, "")
				res.MetaData = a.getMetaData(instance)
				resourceChannel <- res
			}
		}

		if resp.NextToken == nil {
			break
		}
		nextToken = *resp.NextToken
		pageNum++
	}

	return nil
}

func (a *AWSEC2) getMetaData(instance types.Instance) []resource.MetaData {
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
