package types

import "strings"

// EventName represents the name of the event.
type EventName int

const (
	_ EventName = iota
	// Closed is the event that a current `Station` closed
	Closed
	// Linked is the event that a `Station` made a `LINK` request to this `Station`.
	Linked
	// Blocked is the event that a `Station` made a `BLOCK` request to this `Station`.
	Blocked
	// Signal is an event that a `Station` sent a `SIGNAL` to this `Station`.
	Signal
)

// StationDesc holds description of Station
type StationDesc struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

// Serialize stringifies `StationDesc`
func (sd StationDesc) Serialize() string {
	rst := sd.Image
	if len(sd.Name) > 0 {
		rst = rst + "~" + sd.Name
	}
	return rst
}

// Deserialize set `StationDesc` from stringified `StationDesc`
func (sd StationDesc) Deserialize(v string) StationDesc {
	r := strings.Split(v, "~")
	sd.Image = r[0]
	if len(r) > 1 {
		sd.Name = r[1]
	}

	return sd
}

// Handler is callback for receiving messages
type Handler = func(string, StationDesc)

// MsgTo resolves messages
type MsgTo interface{ To(StationDesc) error }

// MsgFrom resolves messages
type MsgFrom interface{ From(StationDesc) error }

// Grabber grabs messages from specified `Station`
type Grabber interface {
	Signal(Handler)
	Linked(Handler)
	Blocked(Handler)
}

// Loser looses grabbed messages by `Grabber`
type Loser interface {
	Signal(Handler)
	Linked(Handler)
	Blocked(Handler)
}

// Station is adfs
type Station interface {
	Name() string
	Image() string

	On(EventName, Handler)
	Signal(string) MsgTo
	Link(string) MsgTo
	Block(string) MsgFrom
	Grab(StationDesc) Grabber
	Lose(StationDesc) Loser
	Close()
	Log(string)
}
