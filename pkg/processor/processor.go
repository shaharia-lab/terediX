// Package processor handles the processing of resources.
package processor

import (
	"fmt"
	"runtime"
	"time"

	"github.com/shaharia-lab/teredix/pkg/metrics"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/scanner"
	"github.com/shaharia-lab/teredix/pkg/scheduler"
	"github.com/shaharia-lab/teredix/pkg/storage"
	"github.com/sirupsen/logrus"
)

// Processor manages the processing of resources from various sources.
type Processor struct {
	Config   Config
	Storage  storage.Storage
	scanners []scanner.Scanner
	logger   *logrus.Logger
	metrics  *metrics.Collector
}

// Config holds configuration values for the Processor.
type Config struct {
	BatchSize int
}

// NewProcessor initializes a new Processor instance.
func NewProcessor(config Config, storage storage.Storage, scanners []scanner.Scanner, logger *logrus.Logger, metrics *metrics.Collector) Processor {
	return Processor{Config: config, Storage: storage, scanners: scanners, logger: logger, metrics: metrics}
}

// Process initiates resource processing.
// It starts scanners for all sources and processes resources in batches.
func (p *Processor) Process(resourceChan chan resource.Resource, sch scheduler.Scheduler) error {
	// Start a goroutine to process resources as they become available
	go func() {
		p.logger.Info("Starting resource processing goroutine")
		p.processResources(resourceChan)
	}()
	for _, sc := range p.scanners {
		lf := logrus.Fields{"scanner_name": sc.GetName(), "scanner_kind": sc.GetKind()}

		// necessary to store the current scanner in a new variable because of passing this inside closure
		scannerCopyForClosure := sc
		err := sch.AddFunc(sc.GetSchedule(), func() {
			start := time.Now()
			// Track initial memory
			var m1 runtime.MemStats
			runtime.ReadMemStats(&m1)

			err := scannerCopyForClosure.Scan(resourceChan)

			if err != nil {
				p.metrics.CollectTotalProcessErrorCount("scanner_scan")
				p.logger.WithFields(lf).WithError(err).Error("Failed to scan resources")
			}

			duration := time.Since(start)
			// Track end memory
			var m2 runtime.MemStats
			runtime.ReadMemStats(&m2)

			p.metrics.CollectTotalScanTimeDurationInSecs(sc.GetName(), sc.GetKind(), duration.Seconds())
			p.metrics.CollectTotalScanTimeDurationInMs(sc.GetName(), sc.GetKind(), float64(duration.Milliseconds()))

			p.metrics.CollectTotalMemoryUsageByScannerInMB(sc.GetName(), sc.GetKind(), float64(m2.TotalAlloc-m1.TotalAlloc)/(1024*1024))
			p.metrics.CollectTotalMemoryUsageByScannerInKB(sc.GetName(), sc.GetKind(), float64(m2.TotalAlloc-m1.TotalAlloc)/1024)
			p.metrics.CollectTotalMemoryUsageByScannerInBytes(sc.GetName(), sc.GetKind(), float64(m2.TotalAlloc-m1.TotalAlloc))
		})

		if err != nil {
			p.logger.WithFields(lf).WithError(err).Error("Failed to schedule scanner in job queue")
			p.metrics.CollectTotalScannerJobAddedToQueue(sc.GetName(), sc.GetKind(), "failed")
		}

		if err == nil {
			p.metrics.CollectTotalScannerJobAddedToQueue(sc.GetName(), sc.GetKind(), "success")
		}

		p.logger.WithFields(lf).WithFields(logrus.Fields{"schedule": sc.GetSchedule()}).Info("Scanner has been scheduled to run")
	}

	err := sch.Start()
	if err != nil {
		p.logger.WithError(err).Error("Failed to start scheduler")
		p.metrics.CollectTotalProcessErrorCount("scheduler_start")
		return err
	}

	p.logger.Info("Scheduler has been started")
	p.metrics.CollectTotalSchedulerStartCount()
	return nil
}

func (p *Processor) processResources(resourceChan <-chan resource.Resource) {
	var resources []resource.Resource

	const flushTimerInterval = 2 * time.Second

	flushTimer := time.NewTimer(flushTimerInterval)
	defer flushTimer.Stop()

	p.logger.WithFields(logrus.Fields{"resource_channel_flushed_interval_in_secs": flushTimerInterval.Seconds(), "processing_batch_size": p.Config.BatchSize}).Info("Resource channel config")

	for {
		select {
		case res, ok := <-resourceChan:
			if !ok {
				// Channel closed, break out of the loop
				p.logger.Info("Channel closed, break out of the loop")
				break
			}

			data := res.GetMetaData()
			rlf := logrus.Fields{"resource_kind": res.GetKind(), "resource_name": res.GetExternalID(), "resource_version": res.GetVersion(), "total_metadata": len(data.Get())}
			p.logger.WithFields(rlf).Info("Received resource from resource channel")

			resources = append(resources, res)

			if p.shouldProcessBatch(resources) {
				if err := p.processBatch(resources); err != nil {
					p.metrics.CollectTotalProcessErrorCount("batch_resource_processing")
					p.logger.WithFields(logrus.Fields{"total_resources_in_batch": len(resources)}).WithError(err).Error("Error processing batch")
				}
				resources = resetResourceBatch(p.Config.BatchSize)
				flushTimer.Reset(flushTimerInterval)
			}

		case <-flushTimer.C:
			if len(resources) > 0 {
				if err := p.processBatch(resources); err != nil {
					p.metrics.CollectTotalProcessErrorCount("resource_processing_channel_flush")
					p.logger.WithFields(logrus.Fields{"total_resources_in_batch": len(resources)}).WithError(err).Error("Error processing batch during resource channel flush")
				}
				resources = resetResourceBatch(p.Config.BatchSize)
			}
			flushTimer.Reset(flushTimerInterval)
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
	start := time.Now()
	p.logger.WithFields(logrus.Fields{"total_resources_in_batch": len(resources)}).Info("Processing batch of resources")

	if err := p.Storage.Persist(resources); err != nil {
		p.logger.WithField("total_resources_affected", len(resources)).WithError(err).Errorf("failed to persist resources")
		p.metrics.CollectTotalProcessErrorCount("resource_persistence")
		return fmt.Errorf("failed to persist resources: %w", err)
	}

	p.logger.WithField("total_resources", len(resources)).Info("Batch of resources has been processed successfully")

	duration := time.Since(start)

	p.metrics.CollectTotalStorageBatchPersistingLatencyInMs(float64(duration.Milliseconds()))
	return nil
}
