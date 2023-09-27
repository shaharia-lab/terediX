// Package processor handles the processing of resources.
package processor

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kyokomi/emoji"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/scanner"
	"github.com/shaharia-lab/teredix/pkg/scheduler"
	"github.com/shaharia-lab/teredix/pkg/storage"
)

// Processor manages the processing of resources from various sources.
type Processor struct {
	Config   Config
	Storage  storage.Storage
	scanners []scanner.Scanner
}

// Config holds configuration values for the Processor.
type Config struct {
	BatchSize int
}

// NewProcessor initializes a new Processor instance.
func NewProcessor(config Config, storage storage.Storage, scanners []scanner.Scanner) Processor {
	return Processor{Config: config, Storage: storage, scanners: scanners}
}

// Process initiates resource processing.
// It starts scanners for all sources and processes resources in batches.
func (p *Processor) Process(resourceChan chan resource.Resource, sch scheduler.Scheduler) {
	// Start a goroutine to process resources as they become available
	go func() {
		p.processResources(resourceChan)
	}()

	for _, sc := range p.scanners {
		err := sch.AddFunc(sc.GetSchedule(), func() {
			sc.Scan(resourceChan)
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	err := sch.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Processor) processResources(resourceChan <-chan resource.Resource) {
	var resources []resource.Resource

	flushTimer := time.NewTimer(2 * time.Second)
	defer flushTimer.Stop()

	for {
		select {
		case res, ok := <-resourceChan:
			if !ok {
				// Channel closed, break out of the loop
				break
			}
			fmt.Println("Received resource:", res.GetKind(), res.GetName())
			resources = append(resources, res)

			if p.shouldProcessBatch(resources) {
				if err := p.processBatch(resources); err != nil {
					fmt.Println("Error processing batch:", err)
				}
				resources = resetResourceBatch(p.Config.BatchSize)
				flushTimer.Reset(2 * time.Second)
			}

		case <-flushTimer.C:
			if len(resources) > 0 {
				if err := p.processBatch(resources); err != nil {
					fmt.Println("Error processing batch:", err)
				}
				resources = resetResourceBatch(p.Config.BatchSize)
			}
			flushTimer.Reset(2 * time.Second)
		}
	}
}

func (p *Processor) shouldProcessBatch(resources []resource.Resource) bool {
	return len(resources) >= p.Config.BatchSize
}

func resetResourceBatch(capacity int) []resource.Resource {
	return make([]resource.Resource, 0, capacity)
}

func (p *Processor) processBatch(resources []resource.Resource) error {
	fmt.Println("\nProcessing batch of", len(resources), "resources...")

	for _, res := range resources {
		fmt.Println(emoji.Sprintf(":check_mark: Processed resource: [ %s ] - %s", res.GetKind(), res.GetName()))
	}

	if err := p.Storage.Persist(resources); err != nil {
		return errors.New("Failed to persist resources: " + err.Error())
	}
	return nil
}
