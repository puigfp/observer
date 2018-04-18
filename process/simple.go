package process

import (
	"fmt"
	"time"

	influxdb "github.com/influxdata/influxdb/client/v2"
	"github.com/puigfp/observer/util"
)

var computeResponseTimeAggregateMetrics = computeMetricsHOF(`
	SELECT
		MEAN(responseTime) AS responseTime_avg,
		MIN(responseTime) AS responseTime_min,
		MAX(responseTime) AS responseTime_max
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
		PERCENTILE(responseTime, 99) AS responseTime_99thPercentile
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
		COUNT(success) AS success_true_count
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
		COUNT(success) AS success_false_count
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

		query := influxdb.NewQuery(queryString, influxdbClient.Config.Database, influxdbClient.Config.Precision)

		response, err := influxdbClient.Client.Query(query)
		if err != nil {
			return err
		}

		return response.Error()
	}
}
