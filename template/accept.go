package template

import (
	"github.com/lokomotes/station-go/types"
)

// AcceptCallback is invoked when `Station` is linked.
type AcceptCallback func(types.Station, string, types.StationDesc)

// AcceptN accepts only first N `Station`s in this activity instance lifetime
func AcceptN(n int, h AcceptCallback) func(types.Station, []string) {
	return func(s types.Station, _ []string) {
		cnt := n
		var l types.Handler

		acceptor := func(msg string, from types.StationDesc) {
			cnt--
			if cnt == 0 {
				l = nil
			}
			h(s, msg, from)
		}

		l = acceptor

		s.On(types.Linked, func(msg string, from types.StationDesc) {
			if l == nil {
				return
			}
			l(msg, from)
		})
	}
}
