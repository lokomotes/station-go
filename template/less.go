package template

import (
	"github.com/lokomotes/station-go/types"
)

// LessCallback is invoked when `Station` is linked.
type LessCallback func(string, types.StationDesc) string

// Less is message-in and message-out style stateless activity.
func Less(h LessCallback) func(types.Station, []string) {
	return AcceptN(1, func(s types.Station, msg string, from types.StationDesc) {
		rst := h(msg, from)

		s.Block(rst).From(from)
		s.Close()
	})
}
