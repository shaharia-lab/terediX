package scanner

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/metrics"
	"github.com/shaharia-lab/teredix/pkg/scheduler"
	"github.com/shaharia-lab/teredix/pkg/storage"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

// AWSS3Client is an autogenerated mock type for the AWSS3Client type
type AWSS3ClientMock struct {
	mock.Mock
}

// GetBucketTagging provides a mock function with given fields: ctx, params, optFns
func (_m *AWSS3ClientMock) GetBucketTagging(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *s3.GetBucketTaggingOutput
	if rf, ok := ret.Get(0).(func(context.Context, *s3.GetBucketTaggingInput, ...func(*s3.Options)) *s3.GetBucketTaggingOutput); ok {
		r0 = rf(ctx, params, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*s3.GetBucketTaggingOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *s3.GetBucketTaggingInput, ...func(*s3.Options)) error); ok {
		r1 = rf(ctx, params, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListBuckets provides a mock function with given fields: ctx, params, optFns
func (_m *AWSS3ClientMock) ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *s3.ListBucketsOutput
	if rf, ok := ret.Get(0).(func(context.Context, *s3.ListBucketsInput, ...func(*s3.Options)) *s3.ListBucketsOutput); ok {
		r0 = rf(ctx, params, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*s3.ListBucketsOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *s3.ListBucketsInput, ...func(*s3.Options)) error); ok {
		r1 = rf(ctx, params, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func TestAWSS3_Scan(t *testing.T) {
	testCases := []struct {
		name                  string
		sourceFields          []string
		buckets               []types.Bucket
		tags                  []types.Tag
		expectedTotalResource int
		expectedTotalMetaData int
		expectedMetaDataKeys  []string
	}{
		{
			name: "successfully list buckets and map resources",
			sourceFields: []string{
				s3fieldRegion,
				s3fieldARN,
				s3fieldBucketName,
				s3fieldTags,
			},
			buckets: []types.Bucket{
				{Name: aws.String("bucket1")},
				{Name: aws.String("bucket2")},
			},
			tags: []types.Tag{
				{Key: aws.String("tag1"), Value: aws.String("value1")},
			},
			expectedTotalResource: 2,
			expectedTotalMetaData: 4,
			expectedMetaDataKeys:  []string{s3fieldRegion, s3fieldARN, s3fieldBucketName, "tag_tag1"},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s3ClientMock := new(AWSS3ClientMock)
			s3ClientMock.On("ListBuckets", mock.Anything, mock.Anything, mock.Anything).Return(&s3.ListBucketsOutput{Buckets: tt.buckets}, nil)
			s3ClientMock.On("GetBucketTagging", mock.Anything, mock.Anything, mock.Anything).Return(&s3.GetBucketTaggingOutput{TagSet: tt.tags}, nil)

			storageMock := new(storage.Mock)
			storageMock.On("GetNextVersionForResource", mock.Anything, mock.Anything).Return(1, nil)

			sc := config.Source{
				Type: pkg.SourceTypeAWSS3,
				Configuration: map[string]string{
					"region":        "us-east-1",
					"account_id":    "123456789012",
					"access_key":    "xxx",
					"secret_key":    "xxx",
					"session_token": "xxx",
				},
				Fields:    tt.sourceFields,
				DependsOn: nil,
				Schedule:  "",
			}

			s := AWSS3{}
			s.Setup("source-name", sc, NewScannerDependencies(scheduler.NewGoCron(), storageMock, &logrus.Logger{}, metrics.NewCollector()))
			s.S3Client = s3ClientMock

			RunCommonScannerAssertionTest(t, &s, tt.expectedTotalResource, tt.expectedTotalMetaData, tt.expectedMetaDataKeys)
		})
	}
}
