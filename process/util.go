package process

import "time"

func roundSub(t time.Time, d time.Duration) time.Time {
	res := t.Round(d)
	if res.After(t) {
		res = res.Add(-d)
	}
	return res
}

func index(slice []string, str string) (int, bool) {
	for i, elem := range slice {
		if elem == str {
			return i, true
		}
	}
	return 0, false
}
