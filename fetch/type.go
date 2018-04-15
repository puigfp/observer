package fetch

import (
	"sync"
	"time"
)

// metricPoint stores the information about how a request went
type metricPoint struct {
	timestamp    time.Time //
	website      string    // website.Name
	responseTime int64     // response time, in nanoseconds (-1 if something went wrong before getting an HTTP response)
	statusCode   int       // HTTP status code (-1 if something wrong happened before getting an HTTP response)
	status       string    // either response.Status (ex: "200 OK", "404 Not Found") or a string describing what went wrong
	success      bool      // true only when the response HTTP code is in [200, 400[
}

// metricsBuffer is a thread-safe metricPoint slice
type metricsBuffer struct {
	buffer []metricPoint
	lock   sync.Mutex
}
