//pkg/processor/processor.go
package processor

import (
	"fmt"
	"github.com/kyokomi/emoji"
	"infrastructure-discovery/pkg/source"
	"infrastructure-discovery/pkg/storage"
	"sync"

	"infrastructure-discovery/pkg/resource"
)

type Processor struct {
	Sources []source.Source
	Config  Config
	Storage storage.Storage
}

type Config struct {
	BatchSize int
}

func NewProcessor(config Config, storage storage.Storage, sources []source.Source) Processor {
	return Processor{Sources: sources, Config: config, Storage: storage}
}

func (p *Processor) Process() {
	// Create a channel to receive resources from scanners
	resourceChan := make(chan resource.Resource)

	// Start a goroutine to process resources as they become available
	go p.processResources(resourceChan)

	// Start goroutines to scan in parallel
	var wg sync.WaitGroup
	for _, s := range p.Sources {
		wg.Add(1)
		go func(s source.Source) {
			defer wg.Done()
			resources := s.Scanner.Scan()
			// Send resources to the shared channel
			for _, res := range resources {
				resourceChan <- res
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
