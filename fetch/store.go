package fetch

import (
	"time"

	influxdb "github.com/influxdata/influxdb/client/v2"
	"github.com/puigfp/observer/util"
)

// storeMetricsOnce sends metric points to influxDB in a single batch
func storeMetricsOnce(influxdbClient util.InfluxDBClient, metrics []metricPoint) error {
	// create batch
	batchPoints, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Database:        influxdbClient.Config.Database,
		Precision:       influxdbClient.Config.Precision,
		RetentionPolicy: influxdbClient.Config.RetentionPolicy,
	})
	if err != nil {
		return err
	}

	// add points to batch
	for _, metric := range metrics {
		tags, fields := convertMetricPointToTagsAndFields(metric)
		pt, err := influxdb.NewPoint("metrics", tags, fields, metric.timestamp)
		if err != nil {
			return err
		}
		batchPoints.AddPoint(pt)
	}

	// write batch
	return influxdbClient.Client.Write(batchPoints)
}

// storeMetrics runs indefinitely and regularily sends the metrics stored in the provited thread-safe buffer to influxDB
func storeMetrics(influxdbClient util.InfluxDBClient, metricsBuf *metricsBuffer, rate time.Duration) {
	for range time.Tick(rate) {
		metricsBuf.lock.Lock()
		if err := storeMetricsOnce(influxdbClient, metricsBuf.buffer); err != nil {
			util.ErrorLogger.Println("Could not send metrics to InfluxDB:", err)

			// when the database is not available, only the last 1000 metric points are kept in memory
			if len(metricsBuf.buffer) > 1000 {
				metricsBuf.buffer = metricsBuf.buffer[len(metricsBuf.buffer)-1000:]
			}
		} else {
			util.InfoLogger.Printf("Emptied buffer, sent %v metrics to InfluxDB.\n", len(metricsBuf.buffer))
			metricsBuf.buffer = make([]metricPoint, 0)
		}
		metricsBuf.lock.Unlock()
	}
}

// convertMetricPointToTagsAndFields transforms a metric point into values that can be sent directly to influxDB
func convertMetricPointToTagsAndFields(point metricPoint) (map[string]string, map[string]interface{}) {
	tags := map[string]string{
		"website": point.website,
	}

	fields := map[string]interface{}{
		"responseTime": point.responseTime,
		"status":       point.status,
		"statusCode":   point.statusCode,
		"success":      point.success,
	}

	return tags, fields
}
