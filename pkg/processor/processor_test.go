package processor

import (
	"math"
	"sync"
	"teredix/pkg/resource"
	"teredix/pkg/source"
	"teredix/pkg/source/scanner"
	"teredix/pkg/storage"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestProcessor_Process(t *testing.T) {
	firstScannerResources := []resource.Resource{
		{
			Kind:        "Test",
			UUID:        "UUID",
			Name:        "Label",
			ExternalID:  "ExternalID",
			RelatedWith: nil,
			MetaData:    nil,
		},
		{
			Kind:        "Test2",
			UUID:        "UUID2",
			Name:        "Name2",
			ExternalID:  "ExternalID2",
			RelatedWith: nil,
			MetaData:    nil,
		},
	}

	secondScannerResources := []resource.Resource{
		{
			Kind:        "Test",
			UUID:        "UUID",
			Name:        "Label",
			ExternalID:  "ExternalID",
			RelatedWith: nil,
			MetaData:    nil,
		},
		{
			Kind:        "Test2",
			UUID:        "UUID2",
			Name:        "Name2",
			ExternalID:  "ExternalID2",
			RelatedWith: nil,
			MetaData:    nil,
		},
	}

	firstScanner := new(scanner.ScannerMock)
	firstScanner.On("Scan").Return(firstScannerResources)

	secondScanner := new(scanner.ScannerMock)
	secondScanner.On("Scan").Return(secondScannerResources)

	mockStorage := new(storage.StorageMock)
	mockStorage.On("Persist", mock.Anything).Return(nil)

	sources := []source.Source{
		{
			Name:    "test",
			Scanner: firstScanner,
		},
		{
			Name:    "test2",
			Scanner: secondScanner,
		},
	}

	processBatchSize := 2
	pr := NewProcessor(Config{BatchSize: processBatchSize}, mockStorage, sources)
	pr.Process()

	// Ensure all goroutines have completed before proceeding with assertion
	var wg sync.WaitGroup
	wg.Add(len(firstScannerResources) + len(secondScannerResources))
	for i := 0; i < len(firstScannerResources); i++ {
		wg.Done()
	}
	for i := 0; i < len(secondScannerResources); i++ {
		wg.Done()
	}
	wg.Wait()

	firstScanner.AssertNumberOfCalls(t, "Scan", 1)
	secondScanner.AssertNumberOfCalls(t, "Scan", 1)

	// Determine how many times Persist should be called
	expectedPersistCalls := int(math.Ceil(float64(len(firstScannerResources)+len(secondScannerResources)) / float64(processBatchSize)))

	mockStorage.AssertNumberOfCalls(t, "Persist", expectedPersistCalls)
}
