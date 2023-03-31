package processor

import (
	"teredix/pkg/resource"
	"teredix/pkg/source"
	"teredix/pkg/source/scanner"
	"teredix/pkg/storage"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

func TestProcessor_Process(t *testing.T) {

	resourceChan := make(chan resource.Resource, 10)

	firstScanner := new(scanner.Mock)
	firstScanner.On("Scan", resourceChan).Return(nil)

	secondScanner := new(scanner.Mock)
	secondScanner.On("Scan", resourceChan).Return(nil)

	mockStorage := new(storage.Mock)
	mockStorage.On("Persist").Return(nil)

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

	// Set up test data
	config := Config{BatchSize: 2}

	processor := NewProcessor(config, mockStorage, sources)

	// Set up mock storage to return nil error when Persist is called
	mockStorage.On("Persist", mock.Anything).Return(nil)

	// Call the Process method
	go func() {
		processor.Process(resourceChan)
	}()

	// Send some test resources on the channel
	res1 := resource.Resource{Kind: "test1", Name: "resource1"}
	res2 := resource.Resource{Kind: "test2", Name: "resource2"}
	res3 := resource.Resource{Kind: "test3", Name: "resource3"}

	resourceChan <- res1
	resourceChan <- res2
	resourceChan <- res3

	// Wait for the processing to finish
	time.Sleep(500 * time.Millisecond)

	// Check that the expected resources were processed
	mockStorage.AssertNumberOfCalls(t, "Persist", 2)
	mockStorage.AssertCalled(t, "Persist", []resource.Resource{res1, res2})
	mockStorage.AssertCalled(t, "Persist", []resource.Resource{res3})
}
