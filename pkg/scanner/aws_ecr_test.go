package scanner

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/metrics"
	"github.com/shaharia-lab/teredix/pkg/scheduler"
	"github.com/shaharia-lab/teredix/pkg/storage"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"

	ecrTypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
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
	testCases := []struct {
		name                       string
		sourceFields               []string
		awsECRRepositories         []ecrTypes.Repository
		awsECRTags                 []types.Tag
		expectedTotalResource      int
		expectedTotalMetaDataCount int
		expectedMetaDataKeys       []string
	}{
		{
			name: "returns resources",
			sourceFields: []string{
				ecrFieldRepositoryName,
				ecrFieldRepositoryURI,
				ecrFieldArn,
				ecrFieldRegistryID,
				ecrFieldTags,
			},
			awsECRRepositories: []ecrTypes.Repository{
				{
					RepositoryUri:  aws.String("something"),
					RepositoryName: &[]string{"test-repo"}[0],
					RegistryId:     &[]string{"1234567890"}[0],
					RepositoryArn:  &[]string{"arn:aws:ecr:us-west-2:1234567890:repository/test-repo"}[0],
				},
			},
			awsECRTags: []types.Tag{
				{
					Key:   aws.String("Environment"),
					Value: aws.String("prod"),
				},
			},
			expectedTotalResource:      1,
			expectedTotalMetaDataCount: 5,
			expectedMetaDataKeys: []string{
				ecrFieldRepositoryName,
				ecrFieldRepositoryURI,
				ecrFieldArn,
				ecrFieldRegistryID,
				"tag_Environment",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// create mock ECR client
			mockEcrClient := new(EcrClientMock)
			// Define mock output
			mockOutput := &ecr.DescribeRepositoriesOutput{
				Repositories: tc.awsECRRepositories,
			}
			mockEcrClient.On("DescribeRepositories", mock.Anything, mock.Anything).Return(mockOutput, nil)

			// Create a mock client
			mockSvc := new(ResourceTaggingServiceClientMock)

			expectedOutput := &resourcegroupstaggingapi.GetResourcesOutput{
				ResourceTagMappingList: []types.ResourceTagMapping{
					{
						Tags: tc.awsECRTags,
					},
				},
			}

			mockSvc.On("GetResources", mock.Anything, mock.Anything, mock.Anything).Return(expectedOutput, nil)

			mockStorage := new(storage.Mock)
			mockStorage.On("GetNextVersionForResource", "test-source", pkg.ResourceKindAWSECR).Return(1, nil)

			sc := config.Source{
				Type: pkg.SourceTypeAWSECR,
				Configuration: map[string]string{
					"region":        "us-west-2",
					"account_id":    "1234567890",
					"access_key":    "xxx",
					"secret_key":    "xxx",
					"session_token": "xxx",
				},
				Fields:   tc.sourceFields,
				Schedule: "@every 1s",
			}
			er := AWSECR{}
			er.Setup("test-source", sc, NewScannerDependencies(scheduler.NewGoCron(), mockStorage, &logrus.Logger{}, metrics.NewCollector()))
			er.ECRClient = mockEcrClient
			er.ResourceTaggingService = mockSvc

			RunCommonScannerAssertionTest(t, &er, tc.expectedTotalResource, tc.expectedTotalMetaDataCount, tc.expectedMetaDataKeys)
		})

	}
}
