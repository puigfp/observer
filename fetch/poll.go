package fetch

import (
	"net/http"
	"strings"
	"time"

	"github.com/puigfp/observer/util"
)

// pollOnce makes a GET HTTP request and returns a metric point describing how that request went
func pollOnce(site util.Website) metricPoint {
	begin := time.Now()

	resp, err := http.Get(site.URL)
	if err != nil {
		// Unfortunately, golang errors do not have an "error type" field.
		// The error message contains some informations specific to this instance of the error,
		// which is not helpful because we want to be able to get "error types" count.
		// The following lines are not pretty clean, but the idea is to extract the last part of
		// the error message which can look like `Get https://aol.com/: dial tcp: lookup aol.com: no such host`
		// (random example). In that case, the following code would extracts `no such host`.
		// Only storing a value that does not change between instances of the same error
		// simplifies the code that counts the error types over different windows of time.
		errString := err.Error()
		errStringSplits := strings.Split(errString, ":")
		errString = errStringSplits[len(errStringSplits)-1]
		errString = strings.TrimSpace(errString)

		return metricPoint{
			website:      site.Name,
			timestamp:    time.Now(),
			responseTime: -1,
			statusCode:   -1,
			status:       errString,
			success:      false,
		}
	}
	defer resp.Body.Close()

	return metricPoint{
		website:      site.Name,
		timestamp:    time.Now(),
		responseTime: time.Since(begin).Nanoseconds(),
		statusCode:   resp.StatusCode,
		status:       resp.Status,
		success:      200 <= resp.StatusCode && resp.StatusCode < 400,
	}
}

// poll runs indefinitely, calls pollOnce regularily, and sends back the metric points through the channel
func poll(site util.Website, metricsChan chan<- metricPoint) {
	for range time.Tick(site.PollRate) {
		// launch a new goroutine for each request
		// in this way, if the website takes longer than the poll rate to respond, a new request is made anyway
		go func() {
			metricsChan <- pollOnce(site)
		}()
	}
}
