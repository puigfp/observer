package display

import (
	"encoding/json"
	"errors"
	"fmt"

	influxdb "github.com/influxdata/influxdb/client/v2"
	"github.com/puigfp/observer/util"
)

func fetchState(influxdbClient util.InfluxDBClient, st *state) error {
	metrics2m, err := fetchMetrics(influxdbClient, "2m")
	if err != nil {
		return err
	}

	metrics10m, err := fetchMetrics(influxdbClient, "10m")
	if err != nil {
		return err
	}

	metrics1h, err := fetchMetrics(influxdbClient, "1h")
	if err != nil {
		return err
	}

	// TODO: fetch alerts

	st.lock.Lock()
	for name := range st.websites {
		site := website{
			name: name,
		}

		// status
		if m, ok := metrics2m[name]; ok {
			status := m.availability > 0.8
			site.status = &status
		}

		m10m, ok := metrics10m[name]
		if ok {
			site.metrics10m = &m10m
		}
		m1h, ok := metrics1h[name]
		if ok {
			site.metrics1h = &m1h
		}

		st.websites[name] = site
	}

	st.lock.Unlock()
	return nil
}

func fetchMetrics(influxdbClient util.InfluxDBClient, window string) (map[string]metrics, error) {
	queryString := fmt.Sprintf(`
		SELECT
			time, website, success_true_count, success_false_count, status_counts,
			responseTime_avg, responseTime_min, responseTime_max, responseTime_99thPercentile 
		FROM
			metrics_%v
		WHERE
			time > now() - %v - 5m
			AND time < now() - %v - 10s
			AND (success_false_count > 0 OR success_true_count > 0)
		GROUP BY
			website
		ORDER BY time DESC
		LIMIT 1
	`, window, window, window)

	query := influxdb.NewQuery(queryString, influxdbClient.Config.Database, influxdbClient.Config.Precision)

	response, err := influxdbClient.Client.Query(query)
	if err != nil {
		return nil, err
	}
	if err := response.Error(); err != nil {
		return nil, err
	}

	if len(response.Results) == 0 || len(response.Results[0].Series) == 0 {
		return nil, errors.New("empty influxDB response")
	}

	m := make(map[string]metrics)

	for _, series := range response.Results[0].Series {
		if len(series.Values) == 0 {
			continue
		}

		value := series.Values[0]

		websiteMetrics := metrics{}

		name, ok := value[1].(string)
		if !ok {
			fmt.Println("bad name")
			continue
		}

		var availability float64
		trueCount, _ := util.ParseJSONNumber(value[2])
		falseCount, _ := util.ParseJSONNumber(value[3])
		if falseCount == 0 {
			availability = 1
		} else {
			availability = float64(trueCount) / (float64(falseCount) + float64(trueCount))
		}
		websiteMetrics.availability = availability

		responseTimeAvg, _ := util.ParseJSONFloat(value[5])
		websiteMetrics.responseTimeAvg = responseTimeAvg
		responseTimeMin, _ := util.ParseJSONNumber(value[6])
		websiteMetrics.responseTimeMin = responseTimeMin
		responseTimeMax, _ := util.ParseJSONNumber(value[7])
		websiteMetrics.responseTimeMax = responseTimeMax
		responseTime99thPercentile, _ := util.ParseJSONNumber(value[8])
		websiteMetrics.responseTime99thPercentile = responseTime99thPercentile

		statusesJSON, ok := value[4].(string)
		if !ok {
			continue
		}

		var statuses map[string]int
		if err := json.Unmarshal([]byte(statusesJSON), &statuses); err != nil {
			continue
		}
		websiteMetrics.statuses = statuses

		m[name] = websiteMetrics
	}

	return m, nil
}
