package display

import (
	"fmt"
	"sort"

	ui "github.com/gizak/termui"
)

var (
	summaryWidget    *ui.List
	statisticsWidget *ui.Par
	alertsWidget     *ui.List
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

		if st.websites[website].status {
			line += "[[UP]](fg-green)   "
		} else {
			line += "[[DOWN]](fg-red) "
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
	(*statistics) = *ui.NewPar(fmt.Sprintf(`[Last 10 minutes](fg-yellow)
[---------------](fg-yellow)

%v

[Last hour](fg-yellow)
[---------](fg-yellow)

%v`,
		getStatisticsString(st.websites[st.websitesOrder[st.selectedWebsite]].metrics10m),
		getStatisticsString(st.websites[st.websitesOrder[st.selectedWebsite]].metrics1h),
	))

	statistics.BorderLabel = "Statistics"
	statistics.TextFgColor = ui.ColorWhite
	statistics.Height = ui.TermHeight()

	return statistics
}

// getStatisticsString computes a string that is used by refreshStatisticsWidget
func getStatisticsString(m metrics) string {
	s := fmt.Sprintf(`[Availability](fg-bold) %v

[Response time](fg-bold)
- avg:   %v
- min:   %v
- max:   %v
- 99th:  %v

[Status](fg-bold)
`,
		m.availability,
		m.responseTimeAvg,
		m.responseTimeMin,
		m.responseTimeMax,
		m.responseTime99thPercentile,
	)

	statuses := make([]string, 0)

	for status := range m.statuses {
		statuses = append(statuses, status)
	}

	sort.Slice(statuses, func(i, j int) bool {
		fmt.Println(m.statuses[statuses[i]], m.statuses[statuses[i]])
		return m.statuses[statuses[i]] > m.statuses[statuses[j]] || statuses[i] < statuses[j]
	})

	for _, status := range statuses {
		s += fmt.Sprintf("- '%v': %v\n", status, m.statuses[status])
	}

	return s
}
