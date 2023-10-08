package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCollectTotalScannerBuild(t *testing.T) {
	c := NewCollector()

	expected := 5.0
	c.CollectTotalScannerBuild(expected)

	assert.Equal(t, expected, testutil.ToFloat64(totalScannerBuild))
}

func TestCollectTotalScannerBuildByKind(t *testing.T) {
	c := NewCollector()
	kind := "someKind"

	c.CollectTotalScannerBuildByKind(kind)

	assert.Equal(t, 1.0, testutil.ToFloat64(totalScannerBuildByKind.WithLabelValues(kind)))
}

//... (you can follow similar patterns for the other methods)

func TestCollectTotalMemoryUsageByScannerInMB(t *testing.T) {
	c := NewCollector()
	name := "someScanner"
	kind := "someKind"
	expectedMemoryUsage := 500.0

	c.CollectTotalMemoryUsageByScannerInMB(name, kind, expectedMemoryUsage)

	assert.Equal(t, expectedMemoryUsage, testutil.ToFloat64(totalMemoryUsageByScannerInMB.WithLabelValues(name, kind)))
}

func TestCollectTotalScannerJobAddedToQueue(t *testing.T) {
	c := NewCollector()
	name := "someScanner"
	kind := "someKind"
	result := "someResult"

	c.CollectTotalScannerJobAddedToQueue(name, kind, result)

	assert.Equal(t, 1.0, testutil.ToFloat64(totalScannerJobAddedToQueue.WithLabelValues(name, kind, result)))
}

func TestCollectTotalSchedulerStartCount(t *testing.T) {
	c := NewCollector()

	c.CollectTotalSchedulerStartCount()

	assert.Equal(t, 1.0, testutil.ToFloat64(totalSchedulerStartCount))
}

func TestCollectTotalProcessErrorCount(t *testing.T) {
	c := NewCollector()
	failureType := "someFailure"

	c.CollectTotalProcessErrorCount(failureType)

	assert.Equal(t, 1.0, testutil.ToFloat64(totalProcessErrorCount.WithLabelValues(failureType)))
}

func TestRecordScanDurationInSecs(t *testing.T) {
	c := NewCollector()
	name := "someScanner"
	kind := "someKind"
	duration := 2.5

	c.RecordScanTimeInSecs(name, kind, duration)

	assert.Equal(t, duration, testutil.ToFloat64(totalScanTimeDurationInSecs.WithLabelValues(name, kind)))
}

func TestCollectTotalScannerJobStatusCount(t *testing.T) {
	c := NewCollector()
	name := "someScanner"
	kind := "someKind"
	status := "completed"

	c.CollectTotalScannerJobStatusCount(name, kind, status)

	assert.Equal(t, 1.0, testutil.ToFloat64(totalScannerJobStatusCount.WithLabelValues(name, kind, status)))
}

func TestCollectTotalMemoryUsageByScannerInKB(t *testing.T) {
	c := NewCollector()
	name := "someScanner"
	kind := "someKind"
	memoryUsage := 256.0

	c.CollectTotalMemoryUsageByScannerInKB(name, kind, memoryUsage)

	assert.Equal(t, memoryUsage, testutil.ToFloat64(totalMemoryUsageByScannerInKB.WithLabelValues(name, kind)))
}

func TestCollectTotalMemoryUsageByScannerInBytes(t *testing.T) {
	c := NewCollector()
	name := "someScanner"
	kind := "someKind"
	memoryUsage := 1024.0

	c.CollectTotalMemoryUsageByScannerInBytes(name, kind, memoryUsage)

	assert.Equal(t, memoryUsage, testutil.ToFloat64(totalMemoryUsageByScannerInBytes.WithLabelValues(name, kind)))
}

func TestCollectTotalStorageBatchPersistingLatencyInMs(t *testing.T) {
	c := NewCollector()
	latency := 5.5

	c.CollectTotalStorageBatchPersistingLatencyInMs(latency)

	assert.Equal(t, latency, testutil.ToFloat64(totalStorageBatchPersistingLatencyInMs))
}

func TestCollectTotalResourceDiscoveredByScanner(t *testing.T) {
	c := NewCollector()
	name := "someScanner"
	kind := "someKind"
	resourceCount := 50.0

	c.CollectTotalResourceDiscoveredByScanner(name, kind, resourceCount)

	assert.Equal(t, resourceCount, testutil.ToFloat64(totalResourceDiscoveredByScanner.WithLabelValues(name, kind)))
}
