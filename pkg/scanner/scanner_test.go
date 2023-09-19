package scanner

import (
	"reflect"
	"testing"

	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/stretchr/testify/assert"
)

// Data provider structure
type getResourceMetaDataTestCase struct {
	name           string
	inputMapper    *FieldMapper
	expectedOutput map[string]string
}

func TestGetResourceMetaData(t *testing.T) {
	// Mocked functions for demonstration purposes
	mockMappingFunc := func() string {
		return "value"
	}
	mockTagsFunc := func() []ResourceTag {
		return []ResourceTag{{Key: "tagKey", Value: "tagValue"}}
	}

	// Your data provider test cases
	testCases := []getResourceMetaDataTestCase{
		{
			name: "Basic Case",
			inputMapper: NewFieldMapper(
				map[string]func() string{"field1": mockMappingFunc},
				mockTagsFunc,
				[]string{"field1", fieldTags},
			),
			expectedOutput: map[string]string{
				"field1":     "value",
				"tag_tagKey": "tagValue",
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

func TestSources_BuildFromAppConfig(t *testing.T) {
	sources := NewSourceRegistry(GetScannerRegistries())
	tests := []struct {
		name            string
		sourceConfig    config.Source
		expectedScanner Scanner
	}{
		{
			name: "build aws ec2 scanner",
			sourceConfig: config.Source{
				Type: pkg.SourceTypeAWSEC2,
				Configuration: map[string]string{
					"region":        "ap-southeast-1",
					"account_id":    "xxx",
					"access_key":    "xxx",
					"secret_key":    "xxx",
					"session_token": "xxx",
				},
				Fields:   []string{fieldInstanceID, fieldImageID},
				Schedule: config.Schedule{},
			},
			expectedScanner: &AWSEC2{},
		},
		{
			name: "build aws ecr scanner",
			sourceConfig: config.Source{
				Type: pkg.SourceTypeAWSECR,
				Configuration: map[string]string{
					"region":        "ap-southeast-1",
					"account_id":    "xxx",
					"access_key":    "xxx",
					"secret_key":    "xxx",
					"session_token": "xxx",
				},
				Fields:   []string{ecrFieldRegistryID},
				Schedule: config.Schedule{},
			},
			expectedScanner: &AWSECR{},
		},
		{
			name: "build aws rds scanner",
			sourceConfig: config.Source{
				Type: pkg.SourceTypeAWSRDS,
				Configuration: map[string]string{
					"region":        "ap-southeast-1",
					"account_id":    "xxx",
					"access_key":    "xxx",
					"secret_key":    "xxx",
					"session_token": "xxx",
				},
				Fields:   []string{rdsFieldInstanceID},
				Schedule: config.Schedule{},
			},
			expectedScanner: &AWSRDS{},
		},
		{
			name: "build AWS S3 scanner",
			sourceConfig: config.Source{
				Type: pkg.SourceTypeAWSS3,
				Configuration: map[string]string{
					"region":        "ap-southeast-1",
					"account_id":    "xxx",
					"access_key":    "xxx",
					"secret_key":    "xxx",
					"session_token": "xxx",
				},
				Fields:   []string{s3fieldRegion},
				Schedule: config.Schedule{},
			},
			expectedScanner: &AWSS3{},
		},
		{
			name: "build file system scanner",
			sourceConfig: config.Source{
				Type: pkg.SourceTypeFileSystem,
				Configuration: map[string]string{
					"root_directory": "/tmp",
				},
				Fields:   []string{fileSystemFieldMachineHost},
				Schedule: config.Schedule{},
			},
			expectedScanner: &FsScanner{},
		},
		{
			name: "build GitHub repository scanner",
			sourceConfig: config.Source{
				Type: pkg.SourceTypeGitHubRepository,
				Configuration: map[string]string{
					"token":       "xxx",
					"user_or_org": "xxx",
				},
				Fields:   []string{fieldCompany},
				Schedule: config.Schedule{},
			},
			expectedScanner: &GitHubRepositoryScanner{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanners := sources.BuildFromAppConfig(map[string]config.Source{
				"hello": tt.sourceConfig,
			})

			assert.IsType(t, tt.expectedScanner, scanners[0])
		})
	}
}
