package process

import (
	"errors"
	"fmt"
	"time"

	influxdb "github.com/influxdata/influxdb/client/v2"
	"github.com/puigfp/observer/util"
)

func computeAlerts(influxdbClient util.InfluxDBClient, window string, statuses *map[string]bool) error {
	// get the possible alerts
	alertsMap, err := retrieveAlerts(influxdbClient, window)
	if err != nil {
		return err
	}

	// perform edge detection to only keep the relevant status alerts
	alerts := filterAlerts(alertsMap, statuses)

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

func retrieveAlerts(influxdbClient util.InfluxDBClient, window string) (map[string]alert, error) {
	alerts := make(map[string]alert)

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
		return alerts, err
	}

	// process influxDB response
	if len(response.Results) == 0 || len(response.Results[0].Series) == 0 {
		return alerts, errors.New("empty influxDB response")
	}

	for _, series := range response.Results[0].Series {
		website := series.Tags["website"]

		// get columns IDs
		trueCountIndex, ok1 := index(series.Columns, "success_true_count")
		falseCountIndex, ok2 := index(series.Columns, "success_false_count")
		timestampIndex, ok3 := index(series.Columns, "time")
		if !ok1 || !ok2 || !ok3 {
			return alerts, errors.New("influxDB dit not return the expected fields")
		}

		// parse/type assert influxDB values
		trueCount, ok1 := util.ParseJSONNumber(series.Values[0][trueCountIndex])
		falseCount, ok2 := util.ParseJSONNumber(series.Values[0][falseCountIndex])
		timestamp, ok3 := util.ParseJSONNumber(series.Values[0][timestampIndex])
		if !ok1 || !ok2 || !ok3 {
			return alerts, errors.New("success counts returned by influxDB could not be interpreted as integers")
		}

		// compare trueCount/totalCount with 0.8 threshold
		status := falseCount == 0 || float64(trueCount)/float64(trueCount+falseCount) > 0.8

		// add alert object to the map
		alerts[website] = alert{
			timestamp: time.Unix(0, timestamp),
			website:   website,
			status:    status,
		}
	}

	return alerts, nil
}

func filterAlerts(alerts map[string]alert, statuses *map[string]bool) []alert {
	filteredAlerts := make([]alert, 0)

	for website, alert := range alerts {
		curStatus := alert.status
		prevStatus, ok := (*statuses)[website]

		// compare previous status with current status
		if ok == true && curStatus != prevStatus {
			filteredAlerts = append(filteredAlerts, alert)
		}

		// store current status in memory
		(*statuses)[website] = curStatus
	}

	return filteredAlerts
}
