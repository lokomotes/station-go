package main

import (
	"strconv"

	metro "github.com/lokomotes/metro/api"
	"github.com/lokomotes/station-go/types"
)

// SignalMsg holds message to be sent
type SignalMsg struct {
	sender *Station
	msg    string
}

// To resolve SignalMsg
func (m SignalMsg) To(s types.StationDesc) error {
	var (
		src = &metro.Station{
			Id:    m.sender.flowID,
			Name:  m.sender.name,
			Image: m.sender.image,
		}
		dst = &metro.Station{
			Name:  s.Name,
			Image: s.Image,
		}
		req = &metro.TransmitRequest{
			Token:   token,
			Src:     src,
			Dst:     dst,
			Message: m.msg,
		}
	)

	res, err := routerCli.Transmit(m.sender.context, req)
	if err != nil {
		return err
	}

	code := res.GetCode()
	switch code {
	default:
		return &types.UnexpectedErr{What: strconv.Itoa(int(code))}
	case 200:
		return nil
	case 404:
		return &types.ImageNotFoundErr{What: s.Image}
	}
}
