package process

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	influxdb "github.com/influxdata/influxdb/client/v2"
	"github.com/puigfp/observer/util"
)

// computeStatusesCounts computes the number of times each status appears in the specified timeframe and sends this data to influxDB
func computeStatusesCounts(influxdbClient util.InfluxDBClient, dest string, begin, end time.Time) error {
	// get statuses
	statuses, err := retrieveStatuses(influxdbClient, begin, end)
	if err != nil {
		return err
	}

	// get statuses counts
	statusesCounts, err := retrieveStatusesCounts(influxdbClient, statuses, begin, end)
	if err != nil {
		return err
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
	for website, counts := range statusesCounts {
		// serialize to count object to JSON to store it in a field
		serialized, err := json.Marshal(counts)
		if err != nil {
			return err
		}

		tags := map[string]string{"website": website}
		fields := map[string]interface{}{"status_counts": string(serialized)}

		pt, err := influxdb.NewPoint(fmt.Sprintf("metrics_%v", dest), tags, fields, begin)
		if err != nil {
			return err
		}
		batchPoints.AddPoint(pt)
	}

	// write batch
	return influxdbClient.Client.Write(batchPoints)
}

// retrieveStatuses retrieves from the database a `website` -> `unique statuses slice` map
func retrieveStatuses(influxdbClient util.InfluxDBClient, begin, end time.Time) (map[string][]string, error) {
	status := make(map[string][]string)

	queryString := fmt.Sprintf(
		`
			SELECT DISTINCT(status) AS status
			FROM metrics
			WHERE
				time >= %v AND
				time < %v
			GROUP BY website
		`,
		begin.UnixNano(), end.UnixNano(),
	)

	query := influxdb.NewQuery(queryString, influxdbClient.Config.Database, influxdbClient.Config.Precision)

	response, err := influxdbClient.Client.Query(query)
	if err != nil {
		return status, err
	}
	if err := response.Error(); err != nil {
		return status, err
	}

	for _, series := range response.Results[0].Series {
		websiteStatuses := make([]string, 0)

		statusIndex, ok := index(series.Columns, "status")
		if !ok {
			return status, errors.New("influxDB did not return the expected field")
		}

		for _, value := range series.Values {
			if s, ok := value[statusIndex].(string); ok {
				websiteStatuses = append(websiteStatuses, s)
			} else {
				return status, errors.New("database 'status' field could not be type asserted to string")
			}
		}

		status[series.Tags["website"]] = websiteStatuses
	}

	return status, nil
}

// retrieveStatusCount retrieves from the database the number of times the status appears in the timeframe
func retrieveStatusCount(influxdbClient util.InfluxDBClient, status string, website string, begin, end time.Time) (int, error) {
	queryString := fmt.Sprintf(`
			SELECT COUNT(status) AS count
			FROM metrics
			WHERE
				website = '%v' AND
				status = '%v' AND
				time >= %v AND
				time < %v
		`,
		website, status,
		begin.UnixNano(), end.UnixNano(),
	)

	query := influxdb.NewQuery(queryString, influxdbClient.Config.Database, influxdbClient.Config.Precision)

	response, err := influxdbClient.Client.Query(query)
	if err != nil {
		return 0, err
	}
	if err := response.Error(); err != nil {
		return 0, err
	}

	if len(response.Results) == 0 || len(response.Results[0].Series) == 0 || len(response.Results[0].Series[0].Values) == 0 {
		return 0, errors.New("empty influxDB response")
	}

	columns := response.Results[0].Series[0].Columns
	value := response.Results[0].Series[0].Values[0]

	countIndex, ok := index(columns, "count")
	if !ok {
		return 0, errors.New("influxDB did not return the expected field")
	}

	count, ok := value[countIndex].(json.Number)
	if !ok {
		return 0, errors.New("database 'count' field could not be type asserted to int")
	}

	n, _ := strconv.ParseInt(string(count), 10, 32)

	return int(n), nil
}

// retrieveStatusesCounts retrieves from the database a `website` -> `status` -> `count` map
func retrieveStatusesCounts(influxdbClient util.InfluxDBClient, statuses map[string][]string, begin, end time.Time) (map[string]map[string]int, error) {
	tasks := 0
	counts := make(map[string]map[string]int)
	type result struct {
		website string
		status  string
		count   int
		err     error
	}
	results := make(chan result)

	for website, statusesList := range statuses {
		counts[website] = make(map[string]int)
		for _, status := range statusesList {
			tasks++
			go func(website, status string) {
				count, err := retrieveStatusCount(influxdbClient, status, website, begin, end)
				results <- result{
					website: website,
					status:  status,
					count:   count,
					err:     err,
				}
			}(website, status)
		}
	}

	for task := 0; task < tasks; task++ {
		res := <-results
		if res.err != nil {
			return counts, res.err
		}
		counts[res.website][res.status] = res.count
	}

	return counts, nil
}
