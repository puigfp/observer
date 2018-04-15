package fetch

import (
	"sync"
	"time"
)

// metricPoint stores the information about how a request went
type metricPoint struct {
	timestamp    time.Time
	website      string
	responseTime int64
	statusCode   int
	status       string
	success      bool
}

// metricsBuffer is a thread-safe metricPoint slice
type metricsBuffer struct {
	buffer []metricPoint
	lock   sync.Mutex
}
