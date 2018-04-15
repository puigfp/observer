package process

import (
	"time"

	"github.com/puigfp/observer/util"
)

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

		if err := computeStatusesCounts(influxdbClient, dest, begin, end); err != nil {
			util.ErrorLogger.Println(err)
		}

		util.InfoLogger.Printf("Computed metrics for ['%v', '%v'[ window.", begin, end)
	}
}
