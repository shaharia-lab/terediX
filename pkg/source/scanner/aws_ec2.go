// Package scanner scans targets
package scanner

import (
	"context"
	"fmt"
	"strings"

	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/util"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

const (
	perPage = 100

	fieldInstanceID        = "instanceId"
	fieldImageID           = "imageId"
	fieldPrivateDnsName    = "privateDnsName"
	fieldInstanceType      = "instanceType"
	fieldArchitecture      = "architecture"
	fieldInstanceLifecycle = "instanceLifecycle"
	fieldInstanceState     = "instanceState"
	fieldVpcID             = "vpcId"
	fieldTags              = "tags"
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
				resourceChannel <- a.mapToResource(instance)
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

func (a *AWSEC2) mapToResource(instance types.Instance) resource.Resource {
	var repoMeta []resource.MetaData
	for _, mapper := range a.fieldMapper(instance) {
		if util.IsFieldExistsInConfig(mapper.field, a.Fields) || strings.Contains(mapper.field, "tag_") {
			val := mapper.value()
			if val != "" {
				repoMeta = append(repoMeta, resource.MetaData{Key: mapper.field, Value: val})
			}
		}
	}

	res := resource.Resource{
		Name:       *instance.InstanceId,
		Kind:       pkg.ResourceKindAWSEC2,
		UUID:       *instance.InstanceId,
		ExternalID: *instance.InstanceId,
		MetaData: []resource.MetaData{
			{
				Key:   "AWS-EC2-Instance-ID",
				Value: *instance.InstanceId,
			},
			{
				Key:   "AWS-EC2-Image-ID",
				Value: *instance.ImageId,
			},
			{
				Key:   "AWS-EC2-PrivateDnsName",
				Value: *instance.PrivateDnsName,
			},
			{
				Key:   "AWS-EC2-InstanceType",
				Value: string(instance.InstanceType),
			},
			{
				Key:   "AWS-EC2-Architecture",
				Value: string(instance.Architecture),
			},
			{
				Key:   "AWS-EC2-InstanceLifecycle",
				Value: string(instance.InstanceLifecycle),
			},
			{
				Key:   "AWS-EC2-InstanceState",
				Value: string(instance.State.Name),
			},
			{
				Key:   "AWS-EC2-VpcId",
				Value: *instance.VpcId,
			},
		},
	}

	for _, tag := range instance.Tags {
		metaData := resource.MetaData{
			Key:   fmt.Sprintf("AWS-EC2-%s", *tag.Key),
			Value: *tag.Value,
		}
		res.MetaData = append(res.MetaData, metaData)
	}

	for _, sg := range instance.SecurityGroups {
		metaData := resource.MetaData{
			Key:   fmt.Sprintf("AWS-EC2-Security-Group-%s", *sg.GroupId),
			Value: *sg.GroupName,
		}
		res.MetaData = append(res.MetaData, metaData)
	}

	res.MetaData = append(res.MetaData, repoMeta...)
	return res
}

func (a *AWSEC2) fieldMapper(instance types.Instance) []MetaDataMapper {
	fieldMapper := []MetaDataMapper{
		{fieldInstanceID, func() string {
			if instance.InstanceId != nil {
				return *instance.InstanceId
			}

			return ""
		}},
		{fieldImageID, func() string {
			if instance.ImageId != nil {
				return *instance.ImageId
			}

			return ""
		}},
		{fieldPrivateDnsName, func() string {
			if instance.PrivateDnsName != nil {
				return *instance.PrivateDnsName
			}

			return ""
		}},
		{fieldInstanceType, func() string {
			it := string(instance.InstanceType)
			if it != "" {
				return it
			}

			return ""
		}},
		{fieldArchitecture, func() string {
			ia := string(instance.Architecture)
			if ia != "" {
				return ia
			}

			return ""
		}},
		{fieldInstanceLifecycle, func() string {
			il := string(instance.InstanceLifecycle)
			if il != "" {
				return il
			}

			return ""
		}},
		{fieldInstanceState, func() string {
			is := string(instance.State.Name)
			if is != "" {
				return is
			}

			return ""
		}},
		{fieldVpcID, func() string {
			if instance.VpcId != nil {
				return *instance.VpcId
			}

			return ""
		}},
	}

	if util.IsFieldExistsInConfig(fieldTags, a.Fields) {
		for _, tag := range instance.Tags {
			fieldMapper = append(fieldMapper, MetaDataMapper{
				field: fmt.Sprintf("tag_%s", *tag.Key),
				value: func() string {
					return *tag.Value
				},
			})
		}
	}

	return fieldMapper
}
