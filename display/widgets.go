package display

import (
	"fmt"
	"sort"

	ui "github.com/gizak/termui"
)

// refreshSummaryWidget updates an already instanciated summary widget with the current state data
func refreshSummaryWidget(summary *ui.List, st *state) *ui.List {
	summaryList := make([]string, 0, len(st.websites))
	for i, website := range st.websitesOrder {
		line := ""

		if i == st.selectedWebsite {
			line += "[->](fg-yellow) "
		} else {
			line += "   "
		}

		status := st.websites[website].status
		if status != nil {
			if *(st.websites[website].status) {
				line += "[[UP]](fg-green)      "
			} else {
				line += "[[DOWN]](fg-red)    "
			}
		} else {
			line += "[[NO DATA]](fg-yellow) "
		}

		line += st.websites[website].name
		summaryList = append(summaryList, line)
	}

	summary.Items = summaryList
	summary.ItemFgColor = ui.ColorWhite
	summary.BorderLabel = "Summary"
	summary.Height = ui.TermHeight()

	return summary
}

// refreshAlertsWidget updates an already instanciated alerts widget with the current state data
func refreshAlertsWidget(alerts *ui.List, st *state) *ui.List {
	alertsList := make([]string, 0)
	for _, alert := range st.alerts {
		line := ""
		line += fmt.Sprintf("[%02d-%02d %02d:%02d](fg-yellow) ",
			alert.timestamp.Month(), alert.timestamp.Day(),
			alert.timestamp.Hour(), alert.timestamp.Minute())

		if alert.status {
			line += "[[UP]](fg-green)   "
		} else {
			line += "[[DOWN]](fg-red) "
		}

		line += alert.website
		alertsList = append(alertsList, line)
	}

	alerts.Items = alertsList
	alerts.ItemFgColor = ui.ColorWhite
	alerts.BorderLabel = "Alerts"
	alerts.Height = ui.TermHeight()

	return alerts
}

// refreshStatisticsWidget updates an already instanciated statistics widget with the current state data
func refreshStatisticsWidget(statistics *ui.Par, st *state) *ui.Par {

	if st.selectedWebsite < len(st.websites) {
		website := st.websites[st.websitesOrder[st.selectedWebsite]]
		if website.metrics10m != nil && website.metrics1h != nil {
			(*statistics) = *ui.NewPar(fmt.Sprintf(`[Last 10 minutes](fg-yellow)
			[---------------](fg-yellow)
			
			%v
			
			[Last hour](fg-yellow)
			[---------](fg-yellow)
			
			%v`,
				getStatisticsString(*(website.metrics10m)),
				getStatisticsString(*(website.metrics1h)),
			))
		} else {
			(*statistics) = *ui.NewPar("NO DATA")
		}
	} else {
		(*statistics) = *ui.NewPar("NO DATA")
	}

	statistics.BorderLabel = "Statistics"
	statistics.TextFgColor = ui.ColorWhite
	statistics.Height = ui.TermHeight()

	return statistics
}

// getStatisticsString computes a string that is used by refreshStatisticsWidget
func getStatisticsString(m metrics) string {
	s := fmt.Sprintf(`[Availability](fg-bold) %.1f%%

[Response time](fg-bold)
- avg:   %.0fms
- min:   %vms
- max:   %vms
- 99th:  %vms

[Status](fg-bold)
`,
		m.availability*100,
		m.responseTimeAvg/1000000,
		m.responseTimeMin/1000000,
		m.responseTimeMax/1000000,
		m.responseTime99thPercentile/1000000,
	)

	statuses := make([]string, 0)

	for status := range m.statuses {
		statuses = append(statuses, status)
	}

	sort.Slice(statuses, func(i, j int) bool {
		return m.statuses[statuses[i]] > m.statuses[statuses[j]] || statuses[i] < statuses[j]
	})

	for _, status := range statuses {
		s += fmt.Sprintf("- '%v': %v\n", status, m.statuses[status])
	}

	return s
}
