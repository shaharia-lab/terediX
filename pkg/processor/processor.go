// Package processor process resources
package processor

import (
	"fmt"
	"log"
	"sync"
	"teredix/pkg/source"
	"teredix/pkg/storage"

	"github.com/kyokomi/emoji"

	"teredix/pkg/resource"
)

// Processor represent resource processor
type Processor struct {
	Sources []source.Source
	Config  Config
	Storage storage.Storage
}

// Config represent configuration for processor
type Config struct {
	BatchSize int
}

// NewProcessor construct new processor
func NewProcessor(config Config, storage storage.Storage, sources []source.Source) Processor {
	return Processor{Sources: sources, Config: config, Storage: storage}
}

// Process start processing resources
func (p *Processor) Process(resourceChan chan resource.Resource) {

	// Start a goroutine to process resources as they become available
	go p.processResources(resourceChan)

	// Start goroutines to scan in parallel
	var wg sync.WaitGroup
	for _, s := range p.Sources {
		wg.Add(1)
		go func(s source.Source) {
			defer wg.Done()
			err := s.Scanner.ScanSource(resourceChan)
			if err != nil {
				log.Printf("failed to start the scanner. scanner: %s. Error: %s", s.Name, err)
			}
		}(s)
	}
	// Wait for all the scanners to finish
	wg.Wait()

	// Close the channel to signal that all resources have been sent
	close(resourceChan)
}

func (p *Processor) processResources(resourceChan <-chan resource.Resource) {
	batchSize := p.Config.BatchSize
	resources := make([]resource.Resource, 0, batchSize)

	for res := range resourceChan {
		resources = append(resources, res)

		if len(resources) == batchSize {
			err := p.processBatch(resources)
			if err != nil {
				fmt.Println(err.Error())
			}
			resources = make([]resource.Resource, 0, batchSize)
		}
	}

	// Process any remaining resources in the channel
	if len(resources) > 0 {
		err := p.processBatch(resources)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (p *Processor) processBatch(resources []resource.Resource) error {
	fmt.Println("\nProcessing batch of", len(resources), "resources...")
	for _, res := range resources {
		fmt.Println(emoji.Sprintf(":check_mark: Processed resource: [ %s ] - %s", res.Kind, res.Name))
	}

	return p.Storage.Persist(resources)
}
