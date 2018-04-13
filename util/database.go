package util

import (
	influxdb "github.com/influxdata/influxdb/client/v2"
)

type InfluxDBClient struct {
	Config InfluxDBConfig
	Client influxdb.Client
}

func NewInfluxDBClient(config InfluxDBConfig) (InfluxDBClient, error) {
	client, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr:     "http://localhost:8086", //"INFLUX_DB_ADDR",
		Username: "admin",                 //INFLUX_DB_USERNAME,
	})

	return InfluxDBClient{
		Config: config,
		Client: client,
	}, err
}
