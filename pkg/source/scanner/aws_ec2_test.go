package scanner

import (
	"context"
	"errors"
	"teredix/pkg"
	"teredix/pkg/resource"
	"testing"

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

func TestAWSEC2_Scan_Return_Data_Successfully(t *testing.T) {
	// Setup
	mc := new(Ec2ClientMock)
	instanceOutput := &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			{
				Instances: []types.Instance{
					{
						InstanceId:     aws.String("i-1234567890"),
						ImageId:        aws.String("ami-1234567890"),
						PrivateDnsName: aws.String("ip-10-0-0-1.us-west-2.compute.internal"),
						InstanceType:   types.InstanceTypeT2Micro,
						Architecture:   types.ArchitectureValuesI386,
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
			},
		},
	}

	mc.On("DescribeInstances", mock.Anything, mock.Anything, mock.Anything).Return(instanceOutput, nil)
	a := NewAWSEC2("test-source", "us-west-2", "1234567890", mc)

	resCh := make(chan resource.Resource, 1)
	err := a.Scan(resCh)
	assert.NoError(t, err)

	res := <-resCh
	assert.Equal(t, "i-1234567890", res.UUID)
	assert.Equal(t, pkg.ResourceKindAWSEC2, res.Kind)
	assert.Equal(t, "i-1234567890", res.ExternalID)
	assert.Len(t, res.MetaData, 10)

	/*// Case 2: Error
	mc.DescribeInstancesFn = func(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
		return nil, errors.New("test-error")
	}
	err = a.Scan(resCh)
	assert.Error(t, err)*/
}

func TestAWSEC2_Scan_Return_Error(t *testing.T) {
	// Setup
	mc := new(Ec2ClientMock)

	mc.On("DescribeInstances", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("failed to fetch"))
	a := NewAWSEC2("test-source", "us-west-2", "1234567890", mc)

	resCh := make(chan resource.Resource, 1)
	err := a.Scan(resCh)
	assert.Error(t, err)
}
