package main

import (
	"strconv"

	metro "github.com/lokomotes/metro/api"
	"github.com/lokomotes/station-go/types"
)

// BlockMsg holds message to be sent
type BlockMsg struct {
	sender *Station
	msg    string
}

// From resolve BlockMsg
func (m BlockMsg) From(s types.StationDesc) error {
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
		req = &metro.BlockRequest{
			Token:   token,
			Src:     src,
			Dst:     dst,
			Message: m.msg,
		}
	)

	res, err := routerCli.Block(m.sender.context, req)
	if err != nil {
		return err
	}

	code := res.GetCode()
	switch code {
	default:
		m.sender.blocked[s.Serialize()] = true
		return &types.UnexpectedErr{What: strconv.Itoa(int(code))}
	case 200:
		return nil
	case 404:
		return &types.ImageNotFoundErr{What: s.Image}
	}
}
