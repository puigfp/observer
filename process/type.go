package process

import "time"

// alert stores the information about an alert
type alert struct {
	timestamp time.Time
	website   string
	status    bool
}
