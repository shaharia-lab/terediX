// Package metrics provide metrics for teredix
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var totalScannerBuild = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "teredix_scanner_build_total",
	Help: "The total number of scanner built",
})

var totalScannerBuildByKind = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "teredix_scanner_build_by_kind",
	Help: "The total number of scanner build by kind",
}, []string{"scanner_kind"})

var totalScannerJobAddedToQueue = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "teredix_scanner_job_added_to_queue",
	Help: "The total number of scanner job added to queue",
}, []string{"scanner_name", "scanner_kind", "result"})

var totalSchedulerStartCount = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "teredix_scheduler_start_count",
	Help: "The total number of scheduler start count",
})

var totalScannerJobStatusCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "teredix_scanner_job_status_count",
	Help: "The total number of scanner job status count",
}, []string{"scanner_name", "scanner_kind", "status"})

var totalProcessErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "teredix_process_error_count",
	Help: "The total number of process error count",
}, []string{"failure_type"})

var scanDurationInSec = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "teredix_scan_duration_seconds",
	Help: "Duration taken by scanner to finish a job",
}, []string{"scanner_name", "scanner_kind"})

var totalMemoryUsageByScannerInMB = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "teredix_memory_usage_by_scanner_in_mb",
	Help: "The total memory usage by scanner in Megabytes",
}, []string{"scanner_name", "scanner_kind"})

var totalMemoryUsageByScannerInKB = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "teredix_memory_usage_by_scanner_in_kb",
	Help: "The total memory usage by scanner in KiloBytes",
}, []string{"scanner_name", "scanner_kind"})

var totalMemoryUsageByScannerInBytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "teredix_memory_usage_by_scanner_in_bytes",
	Help: "The total memory usage by scanner in bytes",
}, []string{"scanner_name", "scanner_kind"})

var totalStorageBatchPersistingLatencyInMs = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "teredix_storage_batch_persisting_latency_in_ms",
	Help: "The total storage batch persisting latency in milliseconds",
})

var totalResourceDiscoveredByScanner = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "teredix_resource_discovered_by_scanner",
	Help: "The total resource discovered by scanner",
}, []string{"scanner_name", "scanner_kind"})

var resourceCounter = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "teredix_total_resource",
		Help: "Total count of resources by source and kind for the latest version.",
	},
	[]string{"scanner_name", "scanner_kind"},
)

var resourceCounterByMetadata = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "teredix_total_resource_by_metadata_key_value",
		Help: "Total count of resources by metadata for the latest version.",
	},
	[]string{"scanner_name", "scanner_kind", "metadata_key", "metadata_value"},
)

// Collector collect metrics
type Collector struct {
}

// NewCollector create new collector
func NewCollector() *Collector {
	return &Collector{}
}

// CollectTotalScannerBuild collect total scanner build
func (c *Collector) CollectTotalScannerBuild(totalScanners float64) {
	totalScannerBuild.Set(totalScanners)
}

// CollectTotalScannerBuildByKind collect total scanner build by name
func (c *Collector) CollectTotalScannerBuildByKind(scannerKind string) {
	totalScannerBuildByKind.WithLabelValues(scannerKind).Inc()
}

// CollectTotalScannerJobAddedToQueue collect total scanner job added to queue
func (c *Collector) CollectTotalScannerJobAddedToQueue(scannerName, scannerKind, result string) {
	totalScannerJobAddedToQueue.WithLabelValues(scannerName, scannerKind, result).Inc()
}

// CollectTotalSchedulerStartCount collect total scheduler start count
func (c *Collector) CollectTotalSchedulerStartCount() {
	totalSchedulerStartCount.Set(1)
}

// CollectTotalProcessErrorCount collect total process error count
func (c *Collector) CollectTotalProcessErrorCount(failureType string) {
	totalProcessErrorCount.WithLabelValues(failureType).Inc()
}

// RecordScanTimeInSecs collect total scan time duration in seconds
func (c *Collector) RecordScanTimeInSecs(scannerName, scannerKind string, duration float64) {
	scanDurationInSec.WithLabelValues(scannerName, scannerKind).Set(duration)
}

// CollectTotalScannerJobStatusCount collect total scanner job status count
func (c *Collector) CollectTotalScannerJobStatusCount(scannerName, scannerKind, status string) {
	totalScannerJobStatusCount.WithLabelValues(scannerName, scannerKind, status).Inc()
}

// CollectTotalMemoryUsageByScannerInMB collect total memory usage by scanner
func (c *Collector) CollectTotalMemoryUsageByScannerInMB(scannerName, scannerKind string, memoryUsage float64) {
	totalMemoryUsageByScannerInMB.WithLabelValues(scannerName, scannerKind).Set(memoryUsage)
}

// CollectTotalMemoryUsageByScannerInKB collect total memory usage by scanner
func (c *Collector) CollectTotalMemoryUsageByScannerInKB(scannerName, scannerKind string, memoryUsage float64) {
	totalMemoryUsageByScannerInKB.WithLabelValues(scannerName, scannerKind).Set(memoryUsage)
}

// CollectTotalMemoryUsageByScannerInBytes collect total memory usage by scanner
func (c *Collector) CollectTotalMemoryUsageByScannerInBytes(scannerName, scannerKind string, memoryUsage float64) {
	totalMemoryUsageByScannerInBytes.WithLabelValues(scannerName, scannerKind).Set(memoryUsage)
}

// CollectTotalStorageBatchPersistingLatencyInMs collect total storage batch persisting latency
func (c *Collector) CollectTotalStorageBatchPersistingLatencyInMs(totalLatency float64) {
	totalStorageBatchPersistingLatencyInMs.Set(totalLatency)
}

// CollectTotalResourceDiscoveredByScanner collect total resource discovered by scanner
func (c *Collector) CollectTotalResourceDiscoveredByScanner(scannerName, scannerKind string, totalResourceDiscovered float64) {
	totalResourceDiscoveredByScanner.WithLabelValues(scannerName, scannerKind).Set(totalResourceDiscovered)
}

// CollectTotalResourceCount collect total resource count
func (c *Collector) CollectTotalResourceCount(source, kind string, totalResource int) {
	resourceCounter.WithLabelValues(source, kind).Set(float64(totalResource))
}

// CollectTotalResourceCountByMetaData collect total resource count by metadata key
func (c *Collector) CollectTotalResourceCountByMetaData(source, kind, metaDataKey string, metaDataValue string, totalResource int) {
	resourceCounterByMetadata.WithLabelValues(source, kind, metaDataKey, metaDataValue).Set(float64(totalResource))
}
