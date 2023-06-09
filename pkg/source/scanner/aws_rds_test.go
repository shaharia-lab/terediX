package scanner

import (
	"testing"

	"github.com/shahariaazam/teredix/pkg/resource"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/stretchr/testify/mock"
)

// RdsClient is an autogenerated mock type for the RdsClient type
type RdsClientMock struct {
	mock.Mock
}

// DescribeDBInstancesPages provides a mock function with given fields: _a0, _a1
func (_m *RdsClientMock) DescribeDBInstancesPages(_a0 *rds.DescribeDBInstancesInput, _a1 func(*rds.DescribeDBInstancesOutput, bool) bool) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*rds.DescribeDBInstancesInput, func(*rds.DescribeDBInstancesOutput, bool) bool) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListTagsForResource provides a mock function with given fields: _a0
func (_m *RdsClientMock) ListTagsForResource(_a0 *rds.ListTagsForResourceInput) (*rds.ListTagsForResourceOutput, error) {
	ret := _m.Called(_a0)

	var r0 *rds.ListTagsForResourceOutput
	if rf, ok := ret.Get(0).(func(*rds.ListTagsForResourceInput) *rds.ListTagsForResourceOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rds.ListTagsForResourceOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*rds.ListTagsForResourceInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func TestAWSRDS_Scan(t *testing.T) {
	testCases := []struct {
		name          string
		rdsInstances  []*rds.DBInstance
		tags          []*rds.Tag
		expectedError error
		expectedCount int
	}{
		{
			name: "successfully list RDS instances and map resources",
			rdsInstances: []*rds.DBInstance{
				{DBInstanceIdentifier: aws.String("instance1")},
				{DBInstanceIdentifier: aws.String("instance2")},
			},
			tags: []*rds.Tag{
				{Key: aws.String("Environment"), Value: aws.String("Production")},
			},
			expectedError: nil,
			expectedCount: 2,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			rdsClientMock := new(RdsClientMock)
			rdsClientMock.On("DescribeDBInstancesPages", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
				fn := args.Get(1).(func(*rds.DescribeDBInstancesOutput, bool) bool)
				for _, rdsInstance := range tt.rdsInstances {
					fn(&rds.DescribeDBInstancesOutput{DBInstances: []*rds.DBInstance{rdsInstance}}, false)
				}
			})

			rdsClientMock.On("ListTagsForResource", mock.Anything).Return(&rds.ListTagsForResourceOutput{TagList: tt.tags}, nil)

			resourceChannel := make(chan resource.Resource, len(tt.rdsInstances))
			var res []resource.Resource

			go func() {
				// Create an RDS scanner and scan
				a := NewAWSRDS("source-name", "us-east-1", "123456789012", rdsClientMock)
				a.Scan(resourceChannel)

				close(resourceChannel)
			}()

			for r := range resourceChannel {
				res = append(res, r)
			}

			if len(res) != tt.expectedCount {
				t.Errorf("unexpected number of resources: got %d, want %d", len(res), tt.expectedCount)
			}
		})
	}
}
