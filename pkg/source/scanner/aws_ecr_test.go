package scanner

import (
	"context"
	"teredix/pkg"
	"teredix/pkg/resource"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"

	"github.com/aws/aws-sdk-go/aws"

	ecrTypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// EcrClient is an autogenerated mock type for the EcrClient type
type EcrClientMock struct {
	mock.Mock
}

// DescribeImages provides a mock function with given fields: _a0, _a1, _a2
func (_m *EcrClientMock) DescribeImages(_a0 context.Context, _a1 *ecr.DescribeImagesInput, _a2 ...func(*ecr.Options)) (*ecr.DescribeImagesOutput, error) {
	_va := make([]interface{}, len(_a2))
	for _i := range _a2 {
		_va[_i] = _a2[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *ecr.DescribeImagesOutput
	if rf, ok := ret.Get(0).(func(context.Context, *ecr.DescribeImagesInput, ...func(*ecr.Options)) *ecr.DescribeImagesOutput); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ecr.DescribeImagesOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ecr.DescribeImagesInput, ...func(*ecr.Options)) error); ok {
		r1 = rf(_a0, _a1, _a2...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DescribeRepositories provides a mock function with given fields: _a0, _a1, _a2
func (_m *EcrClientMock) DescribeRepositories(_a0 context.Context, _a1 *ecr.DescribeRepositoriesInput, _a2 ...func(*ecr.Options)) (*ecr.DescribeRepositoriesOutput, error) {
	_va := make([]interface{}, len(_a2))
	for _i := range _a2 {
		_va[_i] = _a2[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *ecr.DescribeRepositoriesOutput
	if rf, ok := ret.Get(0).(func(context.Context, *ecr.DescribeRepositoriesInput, ...func(*ecr.Options)) *ecr.DescribeRepositoriesOutput); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ecr.DescribeRepositoriesOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ecr.DescribeRepositoriesInput, ...func(*ecr.Options)) error); ok {
		r1 = rf(_a0, _a1, _a2...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRepositoryPolicy provides a mock function with given fields: ctx, params, optFns
func (_m *EcrClientMock) GetRepositoryPolicy(ctx context.Context, params *ecr.GetRepositoryPolicyInput, optFns ...func(*ecr.Options)) (*ecr.GetRepositoryPolicyOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *ecr.GetRepositoryPolicyOutput
	if rf, ok := ret.Get(0).(func(context.Context, *ecr.GetRepositoryPolicyInput, ...func(*ecr.Options)) *ecr.GetRepositoryPolicyOutput); ok {
		r0 = rf(ctx, params, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ecr.GetRepositoryPolicyOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ecr.GetRepositoryPolicyInput, ...func(*ecr.Options)) error); ok {
		r1 = rf(ctx, params, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResourceTaggingServiceClient is an autogenerated mock type for the ResourceTaggingServiceClient type
type ResourceTaggingServiceClientMock struct {
	mock.Mock
}

// GetResources provides a mock function with given fields: _a0, _a1, _a2
func (_m *ResourceTaggingServiceClientMock) GetResources(_a0 context.Context, _a1 *resourcegroupstaggingapi.GetResourcesInput, _a2 ...func(*resourcegroupstaggingapi.Options)) (*resourcegroupstaggingapi.GetResourcesOutput, error) {
	_va := make([]interface{}, len(_a2))
	for _i := range _a2 {
		_va[_i] = _a2[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *resourcegroupstaggingapi.GetResourcesOutput
	if rf, ok := ret.Get(0).(func(context.Context, *resourcegroupstaggingapi.GetResourcesInput, ...func(*resourcegroupstaggingapi.Options)) *resourcegroupstaggingapi.GetResourcesOutput); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*resourcegroupstaggingapi.GetResourcesOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *resourcegroupstaggingapi.GetResourcesInput, ...func(*resourcegroupstaggingapi.Options)) error); ok {
		r1 = rf(_a0, _a1, _a2...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func TestAWSECR_Scan(t *testing.T) {
	// create mock ECR client
	mockEcrClient := new(EcrClientMock)

	// Create a mock client
	mockSvc := new(ResourceTaggingServiceClientMock)

	expectedOutput := &resourcegroupstaggingapi.GetResourcesOutput{
		ResourceTagMappingList: []types.ResourceTagMapping{
			{
				Tags: []types.Tag{
					{
						Key:   aws.String("Environment"),
						Value: aws.String("prod"),
					},
					{
						Key:   aws.String("Owner"),
						Value: aws.String("john@example.com"),
					},
				},
			},
		},
	}

	mockSvc.On("GetResources", mock.Anything, mock.Anything, mock.Anything).Return(expectedOutput, nil)

	// create an instance of AWSECR that uses the mock ECR client
	awsecr := NewAWSECR("test-source", "us-west-2", "xxx", mockEcrClient, mockSvc)

	// Define mock output
	mockOutput := &ecr.DescribeRepositoriesOutput{
		Repositories: []ecrTypes.Repository{
			{
				RepositoryUri:  aws.String("something"),
				RepositoryName: &[]string{"test-repo"}[0],
				RegistryId:     &[]string{"1234567890"}[0],
				RepositoryArn:  &[]string{"arn:aws:ecr:us-west-2:1234567890:repository/test-repo"}[0],
			},
		},
	}
	mockEcrClient.On("DescribeRepositories", mock.Anything, mock.Anything).Return(mockOutput, nil)

	// set expectations on the mock ECR client's DescribeImages method
	mockEcrClient.On("DescribeImages",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(
		&ecr.DescribeImagesOutput{
			ImageDetails: []ecrTypes.ImageDetail{
				{
					ImageDigest: aws.String("sha256:1234567890"),
					ImageTags:   []string{"tag1", "tag2", "tag3"},
				},
			},
		}, nil,
	)

	// create channel for resource output
	resourceChannel := make(chan resource.Resource, 1)

	// run Scan method
	err := awsecr.Scan(resourceChannel)

	// assert no errors occurred
	assert.Nil(t, err)

	// assert expected resource was sent to resource channel
	expectedResource := resource.Resource{
		Name:       "test-repo",
		Kind:       pkg.ResourceKindAWSECR,
		UUID:       "arn:aws:ecr:us-west-2:1234567890:repository/test-repo",
		ExternalID: "arn:aws:ecr:us-west-2:1234567890:repository/test-repo",
		MetaData: []resource.MetaData{
			{
				Key:   "AWS-ECR-Repository-Name",
				Value: "test-repo",
			},
			{
				Key:   "AWS-ECR-Image-Digest",
				Value: "sha256:1234567890",
			},
			{
				Key:   "AWS-ECR-Image-Tag",
				Value: "tag1",
			},
			{
				Key:   "AWS-ECR-Image-Tag",
				Value: "tag1",
			},
			{
				Key:   "AWS-ECR-Image-Tag",
				Value: "tag2",
			},
			{
				Key:   "AWS-ECR-Image-Tag",
				Value: "tag3",
			},
			{
				Key:   pkg.MetaKeyScannerLabel,
				Value: "tag3",
			},
		},
	}

	receivedResource := <-resourceChannel

	assert.Equal(t, len(expectedResource.MetaData), len(receivedResource.MetaData))
	assert.Equal(t, expectedResource.Kind, receivedResource.Kind)
	assert.Equal(t, expectedResource.Name, receivedResource.Name)
	assert.Equal(t, expectedResource.ExternalID, receivedResource.ExternalID)

	// assert that all expected calls to mock ECR client's methods were made
	//mockEcrClient.AssertExpectations(t)
}
