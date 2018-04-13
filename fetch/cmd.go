package fetch

import (
	"sync"
	"time"

	"github.com/puigfp/observer/util"
	"github.com/urfave/cli"
)

var Command = cli.Command{
	Name:    "fetch",
	Aliases: []string{"f"},
	Usage:   "Polls websites and sends relevant metrics to influxDB",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Usage: "path to an observer configuration file",
			Value: "config.json",
		},
	},
	Action: func(c *cli.Context) error {
		// read config
		config, err := util.ReadConfigJSON(c.String("config"))
		if err != nil {
			return err
		}

		// init influxDB client
		influxdbClient, err := util.NewInfluxDBClient(config.InfluxDB)
		if err != nil {
			return err
		}

		// init channel use to get back the metrics
		metricsChan := make(chan metricPoint)

		// init thread-safe buffer used to temporalily store the metrics in memory, before sending them to influxDB in a batch
		metricsBuf := metricsBuffer{
			buffer: make([]metricPoint, 0),
			lock:   sync.Mutex{},
		}

		// launch polling goroutines: one for each website to monitor
		for _, website := range config.Websites {
			go poll(website, metricsChan)
		}

		// launch buffer emptying goroutine, which regularily sends the buffer's content to the database
		go storeMetrics(influxdbClient, &metricsBuf, time.Duration(5)*time.Second)

		// retrieve the metrics from the channel synchronously to keep the process running
		for metric := range metricsChan {
			metricsBuf.lock.Lock()
			metricsBuf.buffer = append(metricsBuf.buffer, metric)
			metricsBuf.lock.Unlock()
		}

		return nil
	},
}
