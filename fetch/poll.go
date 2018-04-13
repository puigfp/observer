package fetch

import (
	"net/http"
	"time"

	"github.com/puigfp/observer/util"
)

func pollOnce(site util.Website) metricPoint {
	begin := time.Now()

	resp, err := http.Get(site.URL)
	if err != nil {
		return metricPoint{
			website:      site.Name,
			timestamp:    time.Now(),
			responseTime: 0,
			status:       err.Error(),
			success:      false,
		}
	}
	defer resp.Body.Close()

	return metricPoint{
		website:      site.Name,
		timestamp:    time.Now(),
		responseTime: time.Since(begin).Nanoseconds(),
		status:       resp.Status,
		success:      200 <= resp.StatusCode && resp.StatusCode < 400,
	}
}

func poll(site util.Website, metricsChan chan<- metricPoint) {
	for range time.Tick(time.Duration(site.PollRate) * time.Millisecond) {
		go func() {
			metricsChan <- pollOnce(site)
		}()
	}
}
