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
		Addr:     config.Addr,
		Username: config.Username,
		Password: config.Password,
	})

	return InfluxDBClient{
		Config: config,
		Client: client,
	}, err
}
