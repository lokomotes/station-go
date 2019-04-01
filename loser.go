package main

import (
	"reflect"

	"github.com/lokomotes/station-go/types"
)

type loser struct {
	grabbed map[types.EventName]map[string][]types.Handler
	serizSd string
}

func (g loser) lose(e types.EventName, h types.Handler) {
	gs, _ := g.grabbed[e]
	hs, ok := gs[g.serizSd]
	if !ok {
		return
	}

	hpl := reflect.ValueOf(h)
	for i, h := range hs {
		hpr := reflect.ValueOf(h)
		if hpl != hpr {
			continue
		}
		if i == len(hs)-1 {
			gs[g.serizSd] = hs[:i]
		} else {
			gs[g.serizSd] = append(hs[:i], hs[i+1:]...)
		}
	}
}

func (g loser) Signal(h types.Handler) {
	g.lose(types.Signal, h)
}

func (g loser) Linked(h types.Handler) {
	g.lose(types.Linked, h)
}

func (g loser) Blocked(h types.Handler) {
	g.lose(types.Blocked, h)
}
