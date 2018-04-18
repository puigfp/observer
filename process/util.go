package process

import "time"

// roundSub rounds t to the nearest "multiple" of d that is lower that t
func roundSub(t time.Time, d time.Duration) time.Time {
	res := t.Round(d)
	if res.After(t) {
		res = res.Add(-d)
	}
	return res
}

// index finds the index of the first occurence of a string in a string slice
func index(slice []string, str string) (int, bool) {
	for i, elem := range slice {
		if elem == str {
			return i, true
		}
	}
	return 0, false
}
