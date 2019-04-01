package main

import (
	"github.com/lokomotes/station-go/types"
)

type grabber struct {
	grabbed map[types.EventName]map[string][]types.Handler
	serizSd string
}

func (g grabber) grab(e types.EventName, h types.Handler) {
	gs, _ := g.grabbed[e]
	gs[g.serizSd] = append(gs[g.serizSd], h)
}

func (g grabber) Signal(h types.Handler) {
	g.grab(types.Signal, h)
}

func (g grabber) Linked(h types.Handler) {
	g.grab(types.Linked, h)
}

func (g grabber) Blocked(h types.Handler) {
	g.grab(types.Blocked, h)
}
