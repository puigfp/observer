package process

import (
	"fmt"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
	"github.com/puigfp/observer/util"
)

var computeResponseTimeAggregateMetrics = computeMetricsHOF(`
	SELECT
		MEAN(responseTime) as responseTime_avg,
		MIN(responseTime) as responseTime_min,
		MAX(responseTime) as responseTime_max 
	INTO metrics_%v
	FROM metrics
	WHERE
		responseTime > 0
		AND time >= %v
		AND time < %v
	GROUP BY time(%v, %v), website
`)

var computeResponseTimePercentileMetrics = computeMetricsHOF(`
	SELECT
		PERCENTILE(responseTime, 99) as responseTime_99thPercentile
	INTO metrics_%v
	FROM metrics
	WHERE
		responseTime > 0
		AND time >= %v
		AND time < %v
	GROUP BY time(%v, %v), website
`)

var computeSuccessCountMetrics = computeMetricsHOF(`
	SELECT
		COUNT(success) as success_true_count
	INTO metrics_%v
	FROM metrics
	WHERE
		success = true
		AND time >= %v
		AND time < %v
	GROUP BY time(%v, %v), website
`)

var computeFailCountMetrics = computeMetricsHOF(`
	SELECT
		COUNT(success) as success_false_count
	INTO metrics_%v
	FROM metrics
	WHERE
		success = false
		AND time >= %v
		AND time < %v
	GROUP BY time(%v, %v), website
`)

func computeMetricsHOF(template string) func(influxdbClient util.InfluxDBClient, dest string, begin, end time.Time, window, offset time.Duration) error {
	return func(influxdbClient util.InfluxDBClient, dest string, begin, end time.Time, window, offset time.Duration) error {
		queryString := fmt.Sprintf(
			template,
			dest,
			begin.UnixNano(), end.UnixNano(),
			window.String(), offset.String(),
		)

		query := client.NewQuery(queryString, influxdbClient.Config.Database, influxdbClient.Config.Precision)

		response, err := influxdbClient.Client.Query(query)
		if err != nil {
			return err
		}

		return response.Error()
	}
}

func computeMetricsLoop(influxdbClient util.InfluxDBClient, dest string, window, tick, security time.Duration) {
	for t := range time.Tick(tick) {
		end := roundSub(t.Add(-security), tick)
		begin := end.Add(-window)
		offset := begin.Sub(roundSub(begin, window))

		if err := computeResponseTimeAggregateMetrics(influxdbClient, dest, begin, end, window, offset); err != nil {
			util.ErrorLogger.Println(err)
		}

		if err := computeResponseTimePercentileMetrics(influxdbClient, dest, begin, end, window, offset); err != nil {
			util.ErrorLogger.Println(err)
		}

		if err := computeSuccessCountMetrics(influxdbClient, dest, begin, end, window, offset); err != nil {
			util.ErrorLogger.Println(err)
		}

		if err := computeFailCountMetrics(influxdbClient, dest, begin, end, window, offset); err != nil {
			util.ErrorLogger.Println(err)
		}

		if err := computeStatusCounts(influxdbClient, dest, begin, end); err != nil {
			util.ErrorLogger.Println(err)
		}

		util.InfoLogger.Printf("Computed metrics for ['%v', '%v'[ window.", begin, end)
	}
}
