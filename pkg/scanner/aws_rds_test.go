package scanner

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/scheduler"
	"github.com/shaharia-lab/teredix/pkg/storage"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

// RdsClient is an autogenerated mock type for the RdsClient type
type RdsClientMock struct {
	mock.Mock
}

// DescribeDBInstances provides a mock function with given fields: _a0, _a1, _a2
func (_m *RdsClientMock) DescribeDBInstances(_a0 context.Context, _a1 *rds.DescribeDBInstancesInput, _a2 ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error) {
	_va := make([]interface{}, len(_a2))
	for _i := range _a2 {
		_va[_i] = _a2[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *rds.DescribeDBInstancesOutput
	if rf, ok := ret.Get(0).(func(context.Context, *rds.DescribeDBInstancesInput, ...func(*rds.Options)) *rds.DescribeDBInstancesOutput); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rds.DescribeDBInstancesOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *rds.DescribeDBInstancesInput, ...func(*rds.Options)) error); ok {
		r1 = rf(_a0, _a1, _a2...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func TestAWSRDS_Scan(t *testing.T) {
	testCases := []struct {
		name                  string
		sourceFields          []string
		rdsInstances          []types.DBInstance
		expectedTotalResource int
		expectedTotalMetaData int
		expectedMetaDataKeys  []string
	}{
		{
			name: "successfully list RDS instances and map resources",
			sourceFields: []string{
				rdsFieldInstanceID,
				rdsFieldTags,
				rdsFieldARN,
				rdsFieldRegion,
			},
			rdsInstances: []types.DBInstance{
				{DBInstanceIdentifier: aws.String("instance1"), TagList: []types.Tag{{Key: aws.String("Environment"), Value: aws.String("Production")}}},
				{DBInstanceIdentifier: aws.String("instance2"), TagList: []types.Tag{{Key: aws.String("Environment"), Value: aws.String("Production")}}},
			},
			expectedTotalResource: 2,
			expectedTotalMetaData: 4,
			expectedMetaDataKeys: []string{
				rdsFieldInstanceID,
				"tag_Environment",
				rdsFieldARN,
				rdsFieldRegion,
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			rdsClientMock := new(RdsClientMock)
			rdsInput := &rds.DescribeDBInstancesInput{}

			mockOutput := &rds.DescribeDBInstancesOutput{
				DBInstances: tt.rdsInstances,
			}

			// Return mockOutput as a result of the DescribeDBInstances method
			rdsClientMock.On("DescribeDBInstances", mock.Anything, rdsInput, mock.Anything).Return(mockOutput, nil)

			mockStorage := new(storage.Mock)
			mockStorage.On("GetNextVersionForResource", mock.Anything, mock.Anything).Return(1, nil)

			sc := config.Source{
				Type: pkg.ResourceKindAWSRDS,
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

			rd := AWSRDS{}
			rd.Setup("source-name", sc, NewScannerDependencies(scheduler.NewCron(), mockStorage, &logrus.Logger{}))
			rd.RdsClient = rdsClientMock

			RunCommonScannerAssertionTest(t, &rd, tt.expectedTotalResource, tt.expectedTotalMetaData, tt.expectedMetaDataKeys)
		})
	}
}
