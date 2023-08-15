package processor

import (
	"errors"
	"fmt"
	"testing"

	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/source"
	"github.com/shaharia-lab/teredix/pkg/source/scanner"
	"github.com/shaharia-lab/teredix/pkg/storage"
	"github.com/stretchr/testify/mock"
)

type mockScanner struct {
	resources []resource.Resource
	err       error
}

func (ms *mockScanner) Scan(ch chan<- resource.Resource) error {
	for _, res := range ms.resources {
		ch <- res
	}
	return ms.err
}

type mockStorage struct {
	err error
}

func (ms *mockStorage) Persist(resources []resource.Resource) error {
	return ms.err
}

func TestProcess(t *testing.T) {
	// Mock resources to be scanned.
	resources := []resource.Resource{
		{Kind: "GitHubRepository", Name: "repo1"},
		{Kind: "GitHubRepository", Name: "repo2"},
		{Kind: "GitHubRepository", Name: "repo3"},
	}

	// Scenario 1: Test successful scanning and processing of resources.
	t.Run("Successful processing", func(t *testing.T) {
		resourceChan := make(chan resource.Resource)

		// Setup mock scanner and storage
		firstScanner := new(scanner.Mock)
		firstScanner.On("Scan", resourceChan).Run(func(args mock.Arguments) {
			for _, res := range resources {
				resourceChan <- res
			}
		}).Return(nil)

		mockStorage := new(storage.Mock)
		mockStorage.On("Persist", resources).Return(nil)

		p := NewProcessor(Config{BatchSize: 3}, mockStorage, []source.Source{{Scanner: firstScanner}})
		p.Process(resourceChan)

		// Assert that the expected calls were made
		firstScanner.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})

	// Scenario 2: Test scanning error.
	t.Run("Scanner error", func(t *testing.T) {
		resourceChan := make(chan resource.Resource)

		firstScanner := new(scanner.Mock)
		firstScanner.On("Scan", resourceChan).Return(errors.New("failed scanner"))

		mockStorage := new(storage.Mock)
		p := NewProcessor(Config{BatchSize: 2}, mockStorage, []source.Source{{Scanner: firstScanner}})

		p.Process(resourceChan)

		// Assert that the expected call was made
		firstScanner.AssertExpectations(t)
	})
}

func TestProcessWithDifferentBatchSizes(t *testing.T) {
	// Define a helper function to generate n resources
	generateResources := func(n int) []resource.Resource {
		var resources []resource.Resource
		for i := 0; i < n; i++ {
			resources = append(resources, resource.Resource{
				Kind: "GitHubRepository",
				Name: fmt.Sprintf("repo%d", i+1),
			})
		}
		return resources
	}

	testCases := []struct {
		totalResources  int
		batchSize       int
		expectedBatches int
	}{
		{5, 2, 3},
		{5, 3, 2},
		{5, 5, 1},
		// ... add more scenarios if needed
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("TotalResources:%d,BatchSize:%d", tc.totalResources, tc.batchSize), func(t *testing.T) {
			resources := generateResources(tc.totalResources)
			resourceChan := make(chan resource.Resource)

			// Setup mock scanner and storage
			firstScanner := new(scanner.Mock)
			firstScanner.On("Scan", resourceChan).Run(func(args mock.Arguments) {
				for _, res := range resources {
					resourceChan <- res
				}
			}).Return(nil)

			mockStorage := new(storage.Mock)
			// This will ensure the Persist method is called expectedBatches times
			mockStorage.On("Persist", mock.Anything).Times(tc.expectedBatches).Return(nil)

			p := NewProcessor(Config{BatchSize: tc.batchSize}, mockStorage, []source.Source{{Scanner: firstScanner}})
			p.Process(resourceChan)

			// Assert that the expected calls were made
			firstScanner.AssertExpectations(t)
			mockStorage.AssertExpectations(t)
		})
	}
}
