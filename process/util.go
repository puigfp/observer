package process

import "time"

func roundSub(t time.Time, d time.Duration) time.Time {
	res := t.Round(d)
	if res.After(t) {
		res = res.Add(-d)
	}
	return res
}
