package util

import (
	"encoding/json"
	"io/ioutil"
)

type Website struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	PollRate int    `json:"poll_rate"`
}

type InfluxDBConfig struct {
	Addr            string `json:"address"`
	Username        string `json:"username"`
	Database        string `json:"database"`
	RetentionPolicy string `json:"retention_policy"`
	Precision       string `json:"precision"`
}

type ConfigJSON struct {
	InfluxDB InfluxDBConfig `json:"influxdb"`
	Websites []Website      `json:"websites"`
}

func ReadConfigJSON(path string) (ConfigJSON, error) {
	var config ConfigJSON

	configFileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(configFileContent, &config)

	return config, err
}
