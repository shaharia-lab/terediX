package scanner

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/metrics"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/scheduler"
	"github.com/shaharia-lab/teredix/pkg/storage"
	"github.com/sirupsen/logrus"
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

// runScannerForTest initiates a scan using the provided scanner and collects
// the resources it discovers into a slice. This function is specifically
// designed to help with testing, allowing you to run a scanner and easily
// gather its results for verification.
func runScannerForTest(scanner Scanner) []resource.Resource {
	resourceChannel := make(chan resource.Resource)

	var res []resource.Resource

	go func() {
		scanner.Scan(resourceChannel)
		close(resourceChannel)
	}()

	for r := range resourceChannel {
		res = append(res, r)
	}

	return res
}

func RunCommonScannerAssertionTest(t *testing.T, scanner Scanner, expectedResourceCount int, expectedMetaDataCount int, expectedMetaDataKeys []string) {
	res := runScannerForTest(scanner)

	assert.Equal(t, expectedResourceCount, len(res), fmt.Sprintf("expected %d resource, but got %d resource", expectedResourceCount, len(res)))
	data := res[0].GetMetaData()
	assert.Equal(t, expectedMetaDataCount, len(data.Get()))

	checkIfMetaKeysExistsInResources(t, res, expectedMetaDataKeys)
}

// TestBuildScanners tests the BuildScanners function
func TestBuildScanners(t *testing.T) {
	var testCases = []struct {
		name                 string
		sources              map[string]config.Source
		expectedTotalScanner int
	}{
		{
			name: "build file system scanner successfully",
			sources: map[string]config.Source{
				"fs_one": {
					Type: pkg.SourceTypeFileSystem,
					Configuration: map[string]string{
						"root_directory": "/tmp",
					},
				},
			},
			expectedTotalScanner: 1,
		},
		{
			name: "should not build scanner if scanner type is not supported",
			sources: map[string]config.Source{
				"fs_one": {
					Type: "unsupported type",
					Configuration: map[string]string{
						"root_directory": "/tmp",
					},
				},
			},
			expectedTotalScanner: 0,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			storageMock := new(storage.Mock)

			scanners := BuildScanners(testCase.sources, NewScannerDependencies(scheduler.NewStaticScheduler(), storageMock, &logrus.Logger{}, metrics.NewCollector()))
			assert.Equal(t, testCase.expectedTotalScanner, len(scanners))
		})
	}
}

// checkIfMetaKeysExistsInResources Checks if all the keys in the given list exist in the metaData of all the Resources
func checkIfMetaKeysExistsInResources(t *testing.T, res []resource.Resource, expectedMetaDataKeys []string) {
	for k, v := range res {
		data := v.GetMetaData()
		missingKeys := data.FindMissingKeys(expectedMetaDataKeys)
		if len(missingKeys) > 0 {
			t.Errorf("Metadata missing. Missing keys [%d]: %v", k, missingKeys)
		}
	}
}
