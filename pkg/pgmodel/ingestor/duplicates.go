// This file and its contents are licensed under the Apache License 2.0.
// Please see the included NOTICE for copyright information and
// LICENSE for a copy of the license.

package ingestor

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/timescale/promscale/pkg/log"
	"github.com/timescale/promscale/pkg/pgmodel/metrics"
)

const reportDuplicatesInterval = time.Minute

var (
	launchReporterOnce    sync.Once
	duplicateMetricsTotal uint64
)

func init() {
	atomic.StoreUint64(&duplicateMetricsTotal, 0)
}

func registerDuplicates(duplicateSamples int64) {
	metrics.IngestorDuplicates.With(prometheus.Labels{"type": "metric", "kind": "sample"}).Add(float64(duplicateSamples))
	metrics.IngestorDuplicates.With(prometheus.Labels{"type": "metric", "kind": "writes_to_db"}).Inc()
}

func reportDuplicates(duplicateMetrics uint64) {
	atomic.AddUint64(&duplicateMetricsTotal, duplicateMetrics)
	metrics.IngestorDuplicates.With(prometheus.Labels{"type": "metric", "kind": "metric"}).Add(float64(duplicateMetrics))
	launchReporterOnce.Do(func() {
		go func() {
			report := time.NewTicker(reportDuplicatesInterval)
			for range report.C {
				if atomic.LoadUint64(&duplicateMetricsTotal) != 0 {
					log.Warn("msg", "duplicate data in sample", "total-duplicate-metrics", atomic.LoadUint64(&duplicateMetricsTotal))
					atomic.StoreUint64(&duplicateMetricsTotal, 0)
				}
			}
		}()
	})
}
