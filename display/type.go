package display

import (
	"sync"
	"time"

	ui "github.com/gizak/termui"
)

type state struct {
	websites        map[string]website
	websitesOrder   []string
	selectedWebsite int
	alerts          []alert
	lock            sync.Mutex
}

type website struct {
	name      string
	status    bool
	metrics2m metrics
	metrics1h metrics
}

type metrics struct {
	availability               *float64
	responseTimeAvg            *float64
	responseTimeMin            *int64
	responseTimeMax            *int64
	responseTime99thPercentile *int64
	statuses                   map[string]int
}

type alert struct {
	website   string
	timestamp time.Time
	status    bool
}

type widgets struct {
	summary    *ui.List
	statistics *ui.Par
	alerts     *ui.List
}
