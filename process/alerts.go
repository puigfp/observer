package process

import (
	"errors"
	"fmt"
	"time"

	influxdb "github.com/influxdata/influxdb/client/v2"
	"github.com/puigfp/observer/util"
)

func computeAlertsLoop(influxdbClient util.InfluxDBClient, window string, rate time.Duration) {
	statuses := make(map[string]bool)
	for range time.Tick(rate) {
		if err := computeAlerts(influxdbClient, window, &statuses); err != nil {
			util.ErrorLogger.Println(err)
		} else {
			util.InfoLogger.Println("Computed alerts.")
		}
	}
}

func computeAlerts(influxdbClient util.InfluxDBClient, window string, statuses *map[string]bool) error {
	// get last success true/false counts from influxDB
	queryString := fmt.Sprintf(
		`
			SELECT success_true_count, success_false_count
			FROM metrics_%v
			GROUP BY website
			ORDER BY time DESC
			LIMIT 1
		`,
		window,
	)

	query := influxdb.NewQuery(queryString, influxdbClient.Config.Database, influxdbClient.Config.Precision)

	response, err := influxdbClient.Client.Query(query)
	if err != nil {
		return err
	}

	if len(response.Results) == 0 || len(response.Results[0].Series) == 0 {
		return errors.New("empty influxDB response")
	}

	alerts := make([]alert, 0)

	for _, series := range response.Results[0].Series {
		website := series.Tags["website"]

		// parse/type assert influxDB values
		trueCountIndex, ok1 := index(series.Columns, "success_true_count")
		falseCountIndex, ok2 := index(series.Columns, "success_false_count")
		timestampIndex, ok3 := index(series.Columns, "time")
		if !ok1 || !ok2 || !ok3 {
			return errors.New("influxDB dit not return the expected fields")
		}

		trueCount, ok1 := parseJSONNumber(series.Values[0][trueCountIndex])
		falseCount, ok2 := parseJSONNumber(series.Values[0][falseCountIndex])
		timestamp, ok3 := parseJSONNumber(series.Values[0][timestampIndex])
		if !ok1 || !ok2 || !ok3 {
			return errors.New("success counts returned by influxDB could not be interpreted as integers")
		}

		// comparing trueCount/totalCount with 0.8 threshold
		curStatus := falseCount == 0 || float64(trueCount)/float64(trueCount+falseCount) > 0.8
		prevStatus, ok := (*statuses)[website]

		// comparing previous status with current status
		if ok == true && curStatus != prevStatus {
			alerts = append(alerts, alert{
				timestamp: time.Unix(0, timestamp),
				website:   website,
				status:    curStatus,
			})
			util.InfoLogger.Printf("ALERT %v %v", website, curStatus)
		}

		// store current status in memory
		(*statuses)[website] = curStatus
	}

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
	for _, alert := range alerts {
		pt, err := influxdb.NewPoint("alerts", map[string]string{
			"website": alert.website,
		}, map[string]interface{}{
			"status": alert.status,
		}, alert.timestamp)
		if err != nil {
			return err
		}
		batchPoints.AddPoint(pt)
	}

	// write batch
	return influxdbClient.Client.Write(batchPoints)
}
