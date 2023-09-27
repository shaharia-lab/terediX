package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var totalScannerBuild = promauto.NewCounter(prometheus.CounterOpts{
	Name: "teredix_scanner_build_total",
	Help: "The total number of scanner built",
})

var totalScannerBuildByName = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "teredix_scanner_build_by_name",
	Help: "The total number of scanner build by name",
}, []string{"scanner_name", "scanner_kind"})

var totalScannerJobAddedToQueue = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "teredix_scanner_job_added_to_queue",
	Help: "The total number of scanner job added to queue",
}, []string{"scanner_name", "scanner_kind", "result"})

var totalSchedulerStartCount = promauto.NewCounter(prometheus.CounterOpts{
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

var totalScanTimeDurationInSecs = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "teredix_scan_time_duration_in_secs",
	Help: "The total number of scan time duration in seconds",
}, []string{"scanner_name", "scanner_kind"})

var totalScanTimeDurationInMs = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "teredix_scan_time_duration_in_ms",
	Help: "The total number of scan time duration in milliseconds",
}, []string{"scanner_name", "scanner_kind"})

var totalMemoryUsageByScannerInMB = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "teredix_memory_usage_by_scanner_in_mb",
	Help: "The total memory usage by scanner in Megabytes",
}, []string{"scanner_name", "scanner_kind"})

var totalMemoryUsageByScannerInKB = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "teredix_memory_usage_by_scanner_in_kb",
	Help: "The total memory usage by scanner in KiloBytes",
}, []string{"scanner_name", "scanner_kind"})

type Collector struct {
}

func NewCollector() *Collector {
	return &Collector{}
}

// CollectTotalScannerBuild collect total scanner build
func (c *Collector) CollectTotalScannerBuild(totalScanners float64) {
	totalScannerBuild.Add(totalScanners)
}

// CollectTotalScannerBuildByName collect total scanner build by name
func (c *Collector) CollectTotalScannerBuildByName(scannerName, scannerKind string) {
	totalScannerBuildByName.WithLabelValues(scannerName, scannerKind).Inc()
}

// CollectTotalScannerJobAddedToQueue collect total scanner job added to queue
func (c *Collector) CollectTotalScannerJobAddedToQueue(scannerName, scannerKind, result string) {
	totalScannerJobAddedToQueue.WithLabelValues(scannerName, scannerKind, result).Inc()
}

// CollectTotalSchedulerStartCount collect total scheduler start count
func (c *Collector) CollectTotalSchedulerStartCount() {
	totalSchedulerStartCount.Inc()
}

// CollectTotalProcessErrorCount collect total process error count
func (c *Collector) CollectTotalProcessErrorCount(failureType string) {
	totalProcessErrorCount.WithLabelValues(failureType).Inc()
}

// CollectTotalScanTimeDurationInSecs collect total scan time duration
func (c *Collector) CollectTotalScanTimeDurationInSecs(scannerName, scannerKind string, duration float64) {
	totalScanTimeDurationInSecs.WithLabelValues(scannerName, scannerKind).Add(duration)
}

// CollectTotalScanTimeDurationInMs collect total scan time duration
func (c *Collector) CollectTotalScanTimeDurationInMs(scannerName, scannerKind string, duration float64) {
	totalScanTimeDurationInMs.WithLabelValues(scannerName, scannerKind).Add(duration)
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
