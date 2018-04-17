package display

import (
	ui "github.com/gizak/termui"
)

// initBody initializes the ui.Body global variable with empty widgets and then does the first render
func initBody(st *state) widgets {
	w := widgets{
		summary:    ui.NewList(),
		statistics: ui.NewPar(""),
		alerts:     ui.NewList(),
	}

	st.lock.Lock()
	refreshSummaryWidget(w.summary, st)
	refreshStatisticsWidget(w.statistics, st)
	refreshAlertsWidget(w.alerts, st)
	st.lock.Unlock()

	ui.Body.Rows = []*ui.Row{ui.NewRow(
		ui.NewCol(4, 0, w.summary),
		ui.NewCol(4, 0, w.statistics),
		ui.NewCol(4, 0, w.alerts),
	)}

	return w
}

// render displays the widgets stored in ram on the screen
func render() {
	ui.Body.Align()
	ui.Render(ui.Body)
}

// func renderResizeLoop(st *state) {
// 	go func() {
// 		for range time.Tick(time.Second) {
// 			render(st)
// 		}
// 	}()
// }
