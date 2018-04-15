package util

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type InfluxDBConfig struct {
	Addr            string `json:"address"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	Database        string `json:"database"`
	RetentionPolicy string `json:"retention_policy"`
	Precision       string `json:"precision"`
}

type ConfigJSON struct {
	InfluxDB InfluxDBConfig `json:"influxdb"`
	Websites []struct {
		Name     string `json:"name"`
		URL      string `json:"url"`
		PollRate int    `json:"poll_rate"`
	} `json:"websites"`
}

type Config struct {
	InfluxDB InfluxDBConfig
	Websites map[string]Website
}

type Website struct {
	Name     string
	URL      string
	PollRate time.Duration
}

func ReadConfigJSON(path string) (Config, error) {
	// read file
	configFileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	// parse JSON
	var config ConfigJSON
	err = json.Unmarshal(configFileContent, &config)
	if err != nil {
		return Config{}, err
	}

	// extract websites
	websites := make(map[string]Website)
	for _, website := range config.Websites {
		websites[website.Name] = Website{
			Name:     website.Name,
			URL:      website.URL,
			PollRate: time.Duration(website.PollRate) * time.Millisecond, // convert int to time.Duration
		}
	}

	return Config{
		InfluxDB: config.InfluxDB,
		Websites: websites,
	}, nil
}
