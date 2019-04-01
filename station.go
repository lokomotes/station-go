package main

import (
	"context"
	"fmt"
	"time"

	_ "github.com/lokomotes/station-go/template"
	"github.com/lokomotes/station-go/types"
)

type stationDesc types.StationDesc

// Station is adfs
type Station struct {
	context context.Context
	flowID  string
	name    string
	image   string

	onSignal  []types.Handler
	onLinked  []types.Handler
	onBlocked []types.Handler
	onClosed  []types.Handler
	//          evt                 desc     callbacks
	grabbed map[types.EventName]map[string][]types.Handler
	//			St.ID
	blocked map[string]bool
}

// Name is name of this `Station`
func (s *Station) Name() string { return s.name }

// Image is image of this `Station`
func (s *Station) Image() string { return s.image }

// Link writes a message to be sent to future reachable `Station`
//
// It writes a message but does not send it yet.
// You can send it by specifying the destination `Station` using returned `To` method.
func (s *Station) Link(msg string) types.MsgTo {
	return LinkMsg{
		sender: s,
		msg:    msg,
	}
}

// Block writes a message to be sent telling `Station` not to send any further messages.
//
// It writes a message but does not send it yet.
// You can send it by specifying the destination `Station` using returned `From` method.
func (s *Station) Block(msg string) types.MsgFrom {
	return BlockMsg{
		sender: s,
		msg:    msg,
	}
}

// Signal writes a message to be sent to other `Station`.
//
// It writes a message but does not send it yet.
// You can send it by specifying the destination `Station` using returned `To` method.
func (s *Station) Signal(msg string) types.MsgTo {
	return SignalMsg{
		sender: s,
		msg:    msg,
	}
}

// On register Handler for receiving messages
func (s *Station) On(event types.EventName, h types.Handler) {
	switch event {
	case types.Linked:
		s.onLinked = append(s.onLinked, h)
	case types.Blocked:
		s.onBlocked = append(s.onBlocked, h)
	case types.Signal:
		s.onSignal = append(s.onSignal, h)
	case types.Closed:
		s.onClosed = append(s.onClosed, h)
	}
}

// Grab intercept events from other `Station`.
//
// It describes which `Station` to intercept the messages from.
// You can resolve it using returned methods.
// Note that intercepted mesages can NOT be caught via `on` method.
// If you want to catch messages using `on` method again, release it using `lose` method.
func (s *Station) Grab(sd types.StationDesc) types.Grabber {
	return grabber{
		grabbed: s.grabbed,
		serizSd: stationDesc(sd).serialize(),
	}
}

// Lose releases grabbed messages from other `Station`
//
// It describes which `Station` to release.
// You can resolve it using returned methods.
func (s *Station) Lose(sd types.StationDesc) types.Loser {
	return loser{
		grabbed: s.grabbed,
		serizSd: stationDesc(sd).serialize(),
	}
}

// Close the `Station` and emits `closed` event.
//
// If you want to listen, obviously, should subscribe listener before close.
func (s *Station) Close() {
	s.Log("station closed")
}

// Log prints operand string
func (s *Station) Log(msg string) {
	t := time.Now().UnixNano() / int64(time.Millisecond)
	fmt.Println(fmt.Sprintf("%d %s", t, msg))
}

func (s *Station) emit(e types.EventName, msg string, from types.StationDesc) {
	if ok, _ := s.blocked[from.Serialize()]; ok {
		return
	}

	gs, _ := s.grabbed[e]
	sd := stationDesc(from)
	hs, ok := gs[sd.serialize()]

	if !ok {
		switch e {
		case types.Linked:
			hs = s.onLinked
		case types.Blocked:
			hs = s.onBlocked
		case types.Signal:
			hs = s.onSignal
		case types.Closed:
			hs = s.onClosed
		}
	}

	for _, h := range hs {
		h(msg, from)
	}
}

func (sd stationDesc) serialize() string {
	n := sd.Name
	if len(n) > 0 {
		return sd.Image + "~" + n
	}

	return sd.Image
}
