package process

import (
	"reflect"
	"sort"
	"testing"
	"time"
)

func checkAlertTestResult(alerts map[string]alert, statuses map[string]bool, expectedAlerts []alert, expectedStatuses map[string]bool, t *testing.T) {
	// call edge detection logic
	filteredAlerts := filterAlerts(alerts, &statuses)

	// sort alerts slice
	sort.Slice(filteredAlerts, func(i, j int) bool {
		return filteredAlerts[i].website < filteredAlerts[j].website
	})
	sort.Slice(expectedAlerts, func(i, j int) bool {
		return expectedAlerts[i].website < expectedAlerts[j].website
	})

	// check statuses
	if !reflect.DeepEqual(statuses, expectedStatuses) {
		t.Fail()
	}

	if !reflect.DeepEqual(filteredAlerts, expectedAlerts) {
		t.Fail()
	}
}

func TestEdgeDetectionBecomeUp(t *testing.T) {
	alerts := map[string]alert{
		"becomeUp": alert{
			timestamp: time.Date(2018, 02, 03, 17, 03, 0, 0, time.Local),
			status:    true,
			website:   "becomeUp",
		},
	}

	statuses := map[string]bool{
		"becomeUp": false,
	}

	expectedStatuses := map[string]bool{
		"becomeUp": true,
	}

	expectedAlerts := []alert{
		alert{
			timestamp: time.Date(2018, 02, 03, 17, 03, 0, 0, time.Local),
			status:    true,
			website:   "becomeUp",
		},
	}

	checkAlertTestResult(alerts, statuses, expectedAlerts, expectedStatuses, t)
}

func TestEdgeDetectionBecomeDown(t *testing.T) {
	alerts := map[string]alert{
		"becomeDown": alert{
			timestamp: time.Date(2018, 02, 03, 17, 03, 0, 0, time.Local),
			status:    false,
			website:   "becomeDown",
		},
	}

	statuses := map[string]bool{
		"becomeDown": true,
	}

	expectedStatuses := map[string]bool{
		"becomeDown": false,
	}

	expectedAlerts := []alert{
		alert{
			timestamp: time.Date(2018, 02, 03, 17, 03, 0, 0, time.Local),
			status:    false,
			website:   "becomeDown",
		},
	}

	checkAlertTestResult(alerts, statuses, expectedAlerts, expectedStatuses, t)
}

func TestEdgeDetectionStillUp(t *testing.T) {
	alerts := map[string]alert{
		"stillUp": alert{
			timestamp: time.Date(2018, 02, 03, 17, 03, 0, 0, time.Local),
			status:    true,
			website:   "stillUp",
		},
	}

	statuses := map[string]bool{
		"stillUp": true,
	}

	expectedStatuses := map[string]bool{
		"stillUp": true,
	}

	expectedAlerts := []alert{}

	checkAlertTestResult(alerts, statuses, expectedAlerts, expectedStatuses, t)
}

func TestEdgeDetectionStillDown(t *testing.T) {
	alerts := map[string]alert{
		"stillDown": alert{
			timestamp: time.Date(2018, 02, 03, 17, 03, 0, 0, time.Local),
			status:    false,
			website:   "stillDown",
		},
	}

	statuses := map[string]bool{
		"stillDown": false,
	}

	expectedStatuses := map[string]bool{
		"stillDown": false,
	}

	expectedAlerts := []alert{}

	checkAlertTestResult(alerts, statuses, expectedAlerts, expectedStatuses, t)
}

func TestEdgeDetectionAbsent(t *testing.T) {
	alerts := map[string]alert{
		"absent": alert{
			timestamp: time.Date(2018, 02, 03, 17, 03, 0, 0, time.Local),
			status:    true,
			website:   "absent",
		},
	}

	statuses := map[string]bool{}

	expectedStatuses := map[string]bool{
		"absent": true,
	}

	expectedAlerts := []alert{}

	checkAlertTestResult(alerts, statuses, expectedAlerts, expectedStatuses, t)
}

func TestEdgeDetectionGlobal(t *testing.T) {
	alerts := map[string]alert{
		"becomeUp": alert{
			timestamp: time.Date(2018, 02, 03, 17, 03, 0, 0, time.Local),
			status:    true,
			website:   "becomeUp",
		},
		"becomeDown": alert{
			timestamp: time.Date(2018, 02, 03, 17, 03, 0, 0, time.Local),
			status:    false,
			website:   "becomeDown",
		},
		"stillUp": alert{
			timestamp: time.Date(2018, 02, 03, 17, 03, 0, 0, time.Local),
			status:    true,
			website:   "stillUp",
		},
		"stillDown": alert{
			timestamp: time.Date(2018, 02, 03, 17, 03, 0, 0, time.Local),
			status:    false,
			website:   "stillDown",
		},
		"absent": alert{
			timestamp: time.Date(2018, 02, 03, 17, 03, 0, 0, time.Local),
			status:    true,
			website:   "absent",
		},
	}

	statuses := map[string]bool{
		"becomeUp":   false,
		"becomeDown": true,
		"stillUp":    true,
		"stillDown":  false,
	}

	expectedStatuses := map[string]bool{
		"becomeUp":   true,
		"becomeDown": false,
		"stillUp":    true,
		"stillDown":  false,
		"absent":     true,
	}

	expectedAlerts := []alert{
		alert{
			timestamp: time.Date(2018, 02, 03, 17, 03, 0, 0, time.Local),
			status:    true,
			website:   "becomeUp",
		},
		alert{
			timestamp: time.Date(2018, 02, 03, 17, 03, 0, 0, time.Local),
			status:    false,
			website:   "becomeDown",
		},
	}

	checkAlertTestResult(alerts, statuses, expectedAlerts, expectedStatuses, t)
}
