package display

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	ui "github.com/gizak/termui"
	influxdb "github.com/influxdata/influxdb/client/v2"
	"github.com/puigfp/observer/util"
)

// updateState fetchs the last metrics from the database and updates the UI using the new data
func updateState(influxdbClient util.InfluxDBClient, w widgets, st *state) {
	fetchState(influxdbClient, st)
	st.lock.Lock()
	refreshSummaryWidget(w.summary, st)
	refreshStatisticsWidget(w.statistics, st)
	refreshAlertsWidget(w.alerts, st)
	st.lock.Unlock()
	render()
}

// updateState fetches the last metrics from the database and update the state
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

	alerts, err := fetchAlerts(influxdbClient)
	if err != nil {
		return err
	}

	st.lock.Lock()
	st.alerts = alerts

	for name := range st.websites {
		site := website{
			name: name,
		}

		// status
		if m, ok := metrics2m[name]; ok {
			status := m.availability > 0.8
			site.status = &status
		}

		// metrics over 10m windows
		m10m, ok := metrics10m[name]
		if ok {
			site.metrics10m = &m10m
		}

		// metrics over 1h windows
		m1h, ok := metrics1h[name]
		if ok {
			site.metrics1h = &m1h
		}

		st.websites[name] = site
	}

	st.lock.Unlock()
	return nil
}

// fetchMetrics fetches the last metrics over a certain window from the database
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

// fetchAlerts fetches the last alerts from the database
func fetchAlerts(influxdbClient util.InfluxDBClient) ([]alert, error) {
	queryString := fmt.Sprintf(`
		SELECT
			time, website, status
		FROM
			alerts
		ORDER BY time DESC
		LIMIT %v
	`, ui.TermHeight())

	query := influxdb.NewQuery(queryString, influxdbClient.Config.Database, influxdbClient.Config.Precision)

	response, err := influxdbClient.Client.Query(query)
	if err != nil {
		return nil, err
	}
	if err := response.Error(); err != nil {
		return nil, err
	}

	if !(len(response.Results) == 0 || len(response.Results[0].Series) == 0) {
		alerts := make([]alert, 0)
		for _, a := range response.Results[0].Series[0].Values {
			timestamp, ok := util.ParseJSONNumber(a[0])
			if !ok {
				continue
			}
			website, ok := a[1].(string)
			if !ok {
				fmt.Println(website)
				continue
			}
			status, ok := a[2].(bool)
			if !ok {
				fmt.Println(status)
				continue
			}

			alerts = append(alerts, alert{
				timestamp: time.Unix(0, timestamp),
				website:   website,
				status:    status,
			})
		}

		return alerts, nil
	}

	return nil, nil
}
