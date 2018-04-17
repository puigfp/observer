package display

import ui "github.com/gizak/termui"

// registerListeners registers functions when some events happen
//
// for now, it only setups functions that react to a key press
// those functions may update the state and trigger a re-render of the screen
func registerListeners(w *widgets, st *state) {
	// press q to quit
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	// press up to change website
	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		st.lock.Lock()
		st.selectedWebsite = (st.selectedWebsite + len(st.websites) - 1) % len(st.websites)
		refreshSummaryWidget(w.summary, st)
		refreshStatisticsWidget(w.statistics, st)
		st.lock.Unlock()
		render()
	})

	// press down to change website
	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		st.lock.Lock()
		st.selectedWebsite = (st.selectedWebsite + 1) % len(st.websites)
		refreshSummaryWidget(w.summary, st)
		refreshStatisticsWidget(w.statistics, st)
		st.lock.Unlock()
		render()
	})
}
