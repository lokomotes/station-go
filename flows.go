package main

import (
	"context"

	"github.com/lokomotes/station-go/types"
)

type stationBody struct {
	cancel  context.CancelFunc
	station *Station
}

// 			    flowID      name
type flowsT map[string]map[string]*stationBody

var flows flowsT

func createStatoinBody(flowID, name string) *stationBody {
	ctx, c := context.WithCancel(context.Background())

	s := &Station{
		context:   ctx,
		flowID:    flowID,
		name:      name,
		onLinked:  make([]types.Handler, 0, 1),
		onBlocked: make([]types.Handler, 0, 1),
		onSignal:  make([]types.Handler, 0, 1),
		onClosed:  make([]types.Handler, 0, 1),
		grabbed:   make(map[types.EventName]map[string][]types.Handler),
		blocked:   make(map[string]bool),
	}

	s.grabbed[types.Signal] = make(map[string][]types.Handler)
	s.grabbed[types.Linked] = make(map[string][]types.Handler)
	s.grabbed[types.Blocked] = make(map[string][]types.Handler)

	s.On(types.Closed, func(_ string, _ types.StationDesc) {
		flows.del(flowID, name)
	})

	return &stationBody{
		cancel:  c,
		station: s,
	}
}

func (f flowsT) create(flowID, name string) (*stationBody, bool) {
	stations, ok := f[flowID]
	if !ok {
		stations = make(map[string]*stationBody)
		body := createStatoinBody(flowID, name)
		stations[name] = body
		f[flowID] = stations
		return body, true
	}

	body, ok := stations[name]
	if ok {
		return nil, false
	}

	body = createStatoinBody(flowID, name)
	stations[name] = body
	return body, true
}

func (f flowsT) have(flowID, name string) bool {
	stations, ok := flows[flowID]
	if !ok {
		return false
	}

	_, ok = stations[name]
	return ok
}

func (f flowsT) get(flowID, name string) (*stationBody, bool) {
	stations, ok := f[flowID]
	if !ok {
		return nil, false
	}

	body, ok := stations[name]
	if !ok {
		return nil, false
	}

	return body, true
}

func (f flowsT) del(flowID, name string) bool {
	stations, ok := f[flowID]
	if !ok {
		return false
	}

	body, ok := stations[name]
	if !ok {
		return false
	}

	body.cancel()
	delete(stations, name)

	if len(stations) == 0 {
		delete(f, flowID)
	}

	return true
}

func init() {
	flows = make(flowsT)
}
