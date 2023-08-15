// Package processor handles the processing of resources.
package processor

import (
	"errors"
	"fmt"
	"sync"

	"github.com/kyokomi/emoji"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/source"
	"github.com/shaharia-lab/teredix/pkg/storage"
)

// Processor manages the processing of resources from various sources.
type Processor struct {
	Sources []source.Source
	Config  Config
	Storage storage.Storage
}

// Config holds configuration values for the Processor.
type Config struct {
	BatchSize int
}

// NewProcessor initializes a new Processor instance.
func NewProcessor(config Config, storage storage.Storage, sources []source.Source) Processor {
	return Processor{Sources: sources, Config: config, Storage: storage}
}

// Process initiates resource processing.
// It starts scanners for all sources and processes resources in batches.
func (p *Processor) Process(resourceChan chan resource.Resource) {
	var wg sync.WaitGroup

	// This WaitGroup will be for the processResources goroutine.
	var processWg sync.WaitGroup
	processWg.Add(1) // Add 1 for the processResources goroutine

	// Start a goroutine to process resources as they become available
	go func() {
		defer processWg.Done() // Decrement the counter when processResources completes
		p.processResources(resourceChan)
	}()

	// Start goroutines to scan in parallel
	for _, s := range p.Sources {
		wg.Add(1)
		go func(s source.Source) {
			defer wg.Done()
			if err := s.Scanner.Scan(resourceChan); err != nil {
				fmt.Printf("Failed to start the scanner. Scanner: %s. Error: %s\n", s.Name, err)
			}
		}(s)
	}

	// Wait for all the scanners to finish
	wg.Wait()

	// Close the channel to signal that all resources have been sent
	close(resourceChan)
	fmt.Println("Resource channel has been closed.")

	// Wait for processResources goroutine to complete
	processWg.Wait()
	fmt.Println("processResources has completed.")
}

func (p *Processor) processResources(resourceChan <-chan resource.Resource) {
	var resources []resource.Resource

	for res := range resourceChan {
		fmt.Println("Received resource:", res.Kind, res.Name)
		resources = append(resources, res)

		if p.shouldProcessBatch(resources) {
			if err := p.processBatch(resources); err != nil {
				fmt.Println("Error processing batch:", err)
			}
			resources = resetResourceBatch(p.Config.BatchSize)
		}
	}
	fmt.Println("Exited the resource channel loop.")

	fmt.Println("Checking if there are remaining resources to be processed...")
	if len(resources) > 0 {
		fmt.Println("Processing remaining resources...")
		if err := p.processBatch(resources); err != nil {
			fmt.Println("Error processing batch:", err)
		}
	} else {
		fmt.Println("No remaining resources to be processed.")
	}
}

func (p *Processor) shouldProcessBatch(resources []resource.Resource) bool {
	return len(resources) == p.Config.BatchSize
}

func resetResourceBatch(capacity int) []resource.Resource {
	return make([]resource.Resource, 0, capacity)
}

func (p *Processor) processBatch(resources []resource.Resource) error {
	fmt.Println("\nProcessing batch of", len(resources), "resources...")

	for _, res := range resources {
		fmt.Println(emoji.Sprintf(":check_mark: Processed resource: [ %s ] - %s", res.Kind, res.Name))
	}

	if err := p.Storage.Persist(resources); err != nil {
		return errors.New("Failed to persist resources: " + err.Error())
	}
	return nil
}
