package util

import (
	"encoding/json"
	"strconv"
)

// parseJSONNumber casts an interface{} wrapping a json.Number variable to int64
func ParseJSONNumber(number interface{}) (int64, bool) {
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

// ParseJSONFloat casts an interface{} wrapping a json.Number variable to int64
func ParseJSONFloat(number interface{}) (float64, bool) {
	if number == nil {
		return 0, true
	}

	numberStr, ok := number.(json.Number)
	if !ok {
		return 0, false
	}

	numberFloat, err := strconv.ParseFloat(string(numberStr), 64)
	if err != nil {
		return 0, false
	}

	return numberFloat, true
}
