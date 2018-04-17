package display

import (
	"sync"
	"time"

	ui "github.com/gizak/termui"
)

// state holds all the data needed to render the ui
type state struct {
	websites        map[string]website
	websitesOrder   []string
	selectedWebsite int
	alerts          []alert
	lock            sync.Mutex
}

// website holds a website's metrics
type website struct {
	name       string
	status     bool
	metrics10m metrics
	metrics1h  metrics
}

// metrics holds a website's metrics
type metrics struct {
	availability               string
	responseTimeAvg            string
	responseTimeMin            string
	responseTimeMax            string
	responseTime99thPercentile string
	statuses                   map[string]int
}

// alert holds all the data about an alert
type alert struct {
	website   string
	timestamp time.Time
	status    bool
}

// widgets stores pointers to the 3 widgets that make the ui
type widgets struct {
	summary    *ui.List
	statistics *ui.Par
	alerts     *ui.List
}
