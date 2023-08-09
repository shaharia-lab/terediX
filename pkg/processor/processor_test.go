package processor

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/shahariaazam/teredix/pkg/resource"
	"github.com/shahariaazam/teredix/pkg/source"
	"github.com/shahariaazam/teredix/pkg/source/scanner"
	"github.com/shahariaazam/teredix/pkg/storage"

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
	config := Config{BatchSize: 2, WorkerPoolSize: 1}

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

func TestProcessor_Process_Handle_Error_From_Scanner(t *testing.T) {

	resourceChan := make(chan resource.Resource, 10)

	sc := new(scanner.Mock)
	sc.On("Scan", resourceChan).Return(errors.New("failed scanner"))

	mockStorage := new(storage.Mock)
	mockStorage.On("Persist").Return(nil)

	sources := []source.Source{
		{
			Name:    "test",
			Scanner: sc,
		},
	}

	// Set up test data
	config := Config{BatchSize: 2, WorkerPoolSize: 1}

	processor := NewProcessor(config, mockStorage, sources)

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the Process method
	go func() {
		processor.Process(resourceChan)
	}()

	// Wait for the processing to finish
	time.Sleep(500 * time.Millisecond)

	// Restore stdout
	w.Close()
	os.Stdout = old

	// Read the output from Process method
	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	// Check that the output contains the expected text
	expected := "failed to start the scanner. scanner: test. Error: failed scanner"
	actual := buf.String()
	if actual != expected {
		t.Errorf("Process output = %q, expected %q", actual, expected)
	}
}

func TestProcessor_Process_Handle_Error_From_Storage_During_Persist(t *testing.T) {

	resourceChan := make(chan resource.Resource, 10)

	// Send some test resources on the channel
	res1 := resource.Resource{Kind: "test1", Name: "resource1"}
	res2 := resource.Resource{Kind: "test2", Name: "resource2"}

	sc := new(scanner.Mock)
	sc.On("Scan", resourceChan).Return(nil)

	mockStorage := new(storage.Mock)
	mockStorage.On("Persist", []resource.Resource{res1, res2}).Return(errors.New("error from persist call"))

	sources := []source.Source{
		{
			Name:    "test",
			Scanner: sc,
		},
	}

	// Set up test data
	config := Config{BatchSize: 2, WorkerPoolSize: 1}

	processor := NewProcessor(config, mockStorage, sources)

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the Process method
	go func() {
		processor.Process(resourceChan)
	}()

	resourceChan <- res1
	resourceChan <- res2

	// Wait for the processing to finish
	time.Sleep(500 * time.Millisecond)

	// Restore stdout
	w.Close()
	os.Stdout = old

	// Read the output from Process method
	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	// Check that the output contains the expected text
	expected := "error from persist call"
	actual := buf.String()

	if !strings.Contains(actual, expected) {
		t.Errorf("Process output = %q, expected %q", actual, expected)
	}
}
