package app

import (
	"github.com/lokomotes/station-go/types"
)

// Open is entry point of user `Station`
func Open(s types.Station, msgs []string) {
	s.Close()
}
