package process

import (
	"time"

	"github.com/puigfp/observer/util"
	"github.com/urfave/cli"
)

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

		go computeResponseTimeMetricsLoop(
			influxdbClient, "metrics_2m",
			time.Duration(2)*time.Minute, time.Duration(10)*time.Second, time.Duration(5)*time.Second,
		)

		go computeResponseTimeMetricsLoop(
			influxdbClient, "metrics_10m",
			time.Duration(10)*time.Minute, time.Duration(30)*time.Second, time.Duration(5)*time.Second,
		)

		computeResponseTimeMetricsLoop(
			influxdbClient, "metrics_1h",
			time.Hour, time.Minute, time.Duration(5)*time.Second,
		)

		return nil
	},
}
