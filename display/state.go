package display

import (
	"sort"

	"github.com/puigfp/observer/util"
)

// initState initialize the state object with the data contained in the config file
func initState(config util.Config, st *state) error {
	st.websites = make(map[string]website)
	st.websitesOrder = make([]string, 0)

	for name := range config.Websites {
		st.websites[name] = website{name: name}
		st.websitesOrder = append(st.websitesOrder, name)
	}

	sort.Slice(st.websitesOrder, func(i, j int) bool {
		return st.websitesOrder[i] < st.websitesOrder[j]
	})

	return nil
}
