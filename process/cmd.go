package process

import (
	"time"

	"github.com/puigfp/observer/util"
	"github.com/urfave/cli"
)

// Command is the main function of this package
var Command = cli.Command{
	Name:    "process",
	Aliases: []string{"p"},
	Usage:   "Computes pseudo-metrics from the raw metrics stored in influxDB",
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
		defer influxdbClient.Client.Close()

		// We do not process any window overlapping with the last 5 seconds because the fetcher
		// sends back the metrics every 5 seconds. Computing metrics on a window that overlaps
		// with this timeframe could make this program miss some points that have not yet been
		// sent by the fetcher.
		security := time.Duration(5) * time.Second

		// compute metrics for 2m window every 5 seconds
		go computeMetricsLoop(
			influxdbClient, "2m",
			time.Duration(2)*time.Minute, time.Duration(5)*time.Second,
			security,
		)

		// compute metrics for 10m window every 10 seconds
		go computeMetricsLoop(
			influxdbClient, "10m",
			time.Duration(10)*time.Minute, time.Duration(10)*time.Second,
			security,
		)

		// compute metrics for 1h window every minute
		go computeMetricsLoop(
			influxdbClient, "1h",
			time.Hour, time.Minute,
			security,
		)

		// compute alerts every 10 seconds
		// (this last one is also used to keep the process running)
		computeAlertsLoop(influxdbClient, "2m", time.Duration(10)*time.Second)

		return nil
	},
}
