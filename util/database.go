package util

import (
	influxdb "github.com/influxdata/influxdb/client/v2"
)

// InfluxDBClient stores an InfluxDBConfig and a influxdb.Client (several funcs need to access both of these variables)
type InfluxDBClient struct {
	Config InfluxDBConfig
	Client influxdb.Client
}

// NewInfluxDBClient creates an InfluxDBClient from an InfluxDBConfig
//
// .Client.Close() should be called when the program is done with a client.
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
