package process

import (
	"encoding/json"
	"strconv"
	"time"
)

// parseJSONNumber casts an interface{} wrapping a json.Number variable to int64
func parseJSONNumber(number interface{}) (int64, bool) {
	if number == nil {
		return 0, true
	}

	numberStr, ok := number.(json.Number)
	if !ok {
		return 0, false
	}

	numberInt, err := strconv.ParseInt(string(numberStr), 10, 64)
	if err != nil {
		return 0, false
	}

	return numberInt, true
}

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
