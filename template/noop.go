package template

import (
	"github.com/lokomotes/station-go/types"
)

// Noop prints "noop" and close `Station`
func Noop(s types.Station, msgs []string) {
	s.Log("noop")
	s.Close()
}
