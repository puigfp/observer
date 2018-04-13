package fetch

import (
	"log"
	"time"

	influxdb "github.com/influxdata/influxdb/client/v2"
)

func storeMetricsOnce(influxdbClient influxdb.Client, metrics []metricPoint) error {
	// create batch
	batchPoints, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Database:        "observer", //INFLUX_DB_DATABASE,
		Precision:       "ms",       //INFLUX_DB_PRECISION,
		RetentionPolicy: "autogen",  // INFLUX_DB_RETENTION_POLICY,
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
	return influxdbClient.Write(batchPoints)
}

func storeMetrics(influxdbClient influxdb.Client, metricsBuf *metricsBuffer, rate time.Duration) {
	for range time.Tick(rate) {
		metricsBuf.lock.Lock()
		if err := storeMetricsOnce(influxdbClient, metricsBuf.buffer); err != nil {
			log.Println("Could not send metrics to InfluxDB:", err)

			// when the database is not available, only the last 1000 metric points are kept in memory
			if len(metricsBuf.buffer) > 1000 {
				metricsBuf.buffer = metricsBuf.buffer[len(metricsBuf.buffer)-1000:]
			}
		} else {
			log.Printf("Emptied buffer, sent %v metrics to InfluxDB.\n", len(metricsBuf.buffer))
			metricsBuf.buffer = make([]metricPoint, 0)
		}
		metricsBuf.lock.Unlock()
	}
}

func convertMetricPointToTagsAndFields(point metricPoint) (map[string]string, map[string]interface{}) {
	tags := map[string]string{
		"website": point.website,
	}

	fields := map[string]interface{}{
		"responseTime": point.responseTime,
		"status":       point.status,
		"success":      point.success,
	}

	return tags, fields
}
