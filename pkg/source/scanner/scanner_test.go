package scanner

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/shaharia-lab/teredix/pkg/resource"
)

// Data provider structure
type getResourceMetaDataTestCase struct {
	name           string
	inputMapper    *FieldMapper
	expectedOutput []resource.MetaData
}

func TestGetResourceMetaData(t *testing.T) {
	// Mocked functions for demonstration purposes
	mockMappingFunc := func() string {
		return "value"
	}
	mockTagsFunc := func() []types.Tag {
		return []types.Tag{{Key: aws.String("tagKey"), Value: aws.String("tagValue")}}
	}

	// Your data provider test cases
	testCases := []getResourceMetaDataTestCase{
		{
			name: "Basic Case",
			inputMapper: NewFieldMapper(
				map[string]func() string{"field1": mockMappingFunc},
				mockTagsFunc,
				[]string{"field1", "tags"},
			),
			expectedOutput: []resource.MetaData{
				{Key: "field1", Value: "value"},
				{Key: "tag_tagKey", Value: "tagValue"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualOutput := testCase.inputMapper.getResourceMetaData()

			if !reflect.DeepEqual(actualOutput, testCase.expectedOutput) {
				t.Errorf("Expected %v, but got %v", testCase.expectedOutput, actualOutput)
			}
		})
	}
}

// MockScanner is a mock implementation of the Scanner interface
type MockScanner struct {
	resources []resource.Resource
}

// Scan is a mock method that sends resources to the given channel
func (m *MockScanner) Scan(ch chan<- resource.Resource) {
	for _, r := range m.resources {
		ch <- r
	}
}
