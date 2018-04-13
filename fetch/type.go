package fetch

import (
	"sync"
	"time"
)

type metricPoint struct {
	timestamp    time.Time
	website      string
	responseTime int64
	status       string
	success      bool
}

type metricsBuffer struct {
	buffer []metricPoint
	lock   sync.Mutex
}
