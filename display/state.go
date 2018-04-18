package display

import (
	"sort"

	"github.com/puigfp/observer/util"
)

func initState(config util.Config, st *state) error {
	st.websites = make(map[string]website)
	st.websitesOrder = make([]string, 0)

	for name := range config.Websites {
		st.websites[name] = website{}
		st.websitesOrder = append(st.websitesOrder, name)
	}

	sort.Slice(st.websitesOrder, func(i, j int) bool {
		return st.websitesOrder[i] < st.websitesOrder[j]
	})

	return nil
}

func updateState(influxdbClient util.InfluxDBClient, w widgets, st *state) {
	fetchState(influxdbClient, st)
	st.lock.Lock()
	refreshSummaryWidget(w.summary, st)
	refreshStatisticsWidget(w.statistics, st)
	refreshAlertsWidget(w.alerts, st)
	st.lock.Unlock()
	render()
}
