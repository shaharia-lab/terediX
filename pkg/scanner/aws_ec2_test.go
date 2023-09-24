package scanner

import (
	"context"
	"testing"

	"github.com/go-co-op/gocron"
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/storage"
	"github.com/shaharia-lab/teredix/pkg/util"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/stretchr/testify/mock"
)

// Ec2ClientMock is an autogenerated mock type for the Ec2ClientMock type
type Ec2ClientMock struct {
	mock.Mock
}

// DescribeInstances provides a mock function with given fields: ctx, params, optFns
func (_m *Ec2ClientMock) DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *ec2.DescribeInstancesOutput
	if rf, ok := ret.Get(0).(func(context.Context, *ec2.DescribeInstancesInput, ...func(*ec2.Options)) *ec2.DescribeInstancesOutput); ok {
		r0 = rf(ctx, params, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ec2.DescribeInstancesOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ec2.DescribeInstancesInput, ...func(*ec2.Options)) error); ok {
		r1 = rf(ctx, params, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func TestAWSEC2_Scan(t *testing.T) {
	testCases := []struct {
		name                  string
		sourceFields          []string
		awsEc2Instances       []types.Instance
		expectedTotalResource int
		expectedMetaDataKeys  []string
	}{
		{
			name: "returns resources",
			sourceFields: []string{
				fieldInstanceID,
				fieldImageID,
				fieldPrivateDNSName,
				fieldInstanceType,
				fieldArchitecture,
				fieldInstanceLifecycle,
				fieldInstanceState,
				fieldVpcID,
				fieldTags,
			},
			awsEc2Instances: []types.Instance{
				{
					InstanceId:        aws.String("i-1234567890"),
					ImageId:           aws.String("ami-1234567890"),
					PrivateDnsName:    aws.String("ip-10-0-0-1.us-west-2.compute.internal"),
					InstanceType:      types.InstanceTypeT2Micro,
					Architecture:      types.ArchitectureValuesI386,
					InstanceLifecycle: "scheduled",
					State: &types.InstanceState{
						Name: types.InstanceStateNameRunning,
					},
					VpcId: aws.String("vpc-1234567890"),
					Tags: []types.Tag{
						{
							Key:   aws.String("Name"),
							Value: aws.String("test-instance"),
						},
					},
					SecurityGroups: []types.GroupIdentifier{
						{
							GroupId:   aws.String("sg-1234567890"),
							GroupName: aws.String("test-security-group"),
						},
					},
				},
			},
			expectedTotalResource: 1,
			expectedMetaDataKeys: []string{
				fieldInstanceID,
				fieldImageID,
				fieldPrivateDNSName,
				fieldInstanceType,
				fieldArchitecture,
				fieldInstanceLifecycle,
				fieldInstanceState,
				fieldVpcID,
				"tag_Name",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mc := new(Ec2ClientMock)
			instanceOutput := &ec2.DescribeInstancesOutput{
				Reservations: []types.Reservation{
					{
						Instances: tc.awsEc2Instances,
					},
				},
			}
			mc.On("DescribeInstances", mock.Anything, mock.Anything, mock.Anything).Return(instanceOutput, nil)

			sm := new(storage.Mock)
			sm.On("GetNextVersionForResource", mock.Anything, mock.Anything).Return(1, nil)

			sc := config.Source{
				Type:          "",
				ConfigFrom:    "",
				Configuration: map[string]string{"region": "us-west-2", "account_id": "1234567890"},
				Fields:        tc.sourceFields,
				DependsOn:     nil,
				Schedule:      config.Schedule{},
			}
			e := AWSEC2{}
			e.Build("test-source", sc, sm, &gocron.Scheduler{}, &logrus.Logger{})
			e.setEC2Client(mc)

			res := RunScannerForTests(&e)
			assert.Equal(t, tc.expectedTotalResource, len(res))
			util.CheckIfMetaKeysExistsInResources(t, res, tc.expectedMetaDataKeys)
		})
	}
}
