package main

import (
	"context"
	"io"
	"os"
	"strconv"
	"strings"

	// the script `build.sh` will create it
	// TODO: replace bianary-only-package to plug-in
	"app"

	metro "github.com/lokomotes/metro/api"
	"github.com/lokomotes/station-go/types"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	routerHost string
	routerPort uint16
	routerConn *grpc.ClientConn
	routerCli  metro.RouterClient
)

func signalHandler(sig *metro.Signal) {
	var (
		srcSt = sig.GetSrc()
		dstSt = sig.GetDst()

		flowID  = dstSt.GetId()
		dstName = dstSt.GetName()

		isExists = flows.have(flowID, dstName)
		fetch    = func() *stationBody { b, _ := flows.get(flowID, dstName); return b }
		create   = func() *stationBody { b, _ := flows.create(flowID, dstName); return b }
		start    = func() *stationBody {
			body := create()
			body.station.Log("new station is open")

			var args []string
			if msg := sig.GetMessage(); len(msg) > 1 {
				msg = msg[1 : len(msg)-1]
				args = strings.Split(msg, ",")
			}

			app.Open(body.station, args)

			return body
		}
		toDsc = func(s *metro.Station) types.StationDesc {
			return types.StationDesc{
				Name:  s.GetName(),
				Image: s.GetImage(),
			}
		}

		ctrl = sig.GetControl()
	)

	switch ctrl {
	default:
	case metro.Signal_START:
		if isExists {
			fetch().station.Log("already opened station")
			return
		}

		start()

	case metro.Signal_TERMINATE:
		if !isExists {
			return
		}
		flows.del(flowID, dstName)

	case metro.Signal_LINKED:
		var s *Station
		if isExists {
			s = fetch().station
		} else {
			s = start().station
		}

		s.emit(types.Linked, sig.GetMessage(), toDsc(srcSt))

	case metro.Signal_MESSAGE:
		if !isExists {
			return
		}
		fetch().station.emit(types.Signal, sig.GetMessage(), toDsc(srcSt))

	case metro.Signal_BLOCKED:
		if !isExists {
			return
		}
		fetch().station.emit(types.Blocked, sig.GetMessage(), toDsc(srcSt))
	}
}

func main() {
	var err error

	log.Info(getRouterAddress())
	log.Info(token.GetId())

	routerConn, err = grpc.Dial(getRouterAddress(), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial to Router: %v", err)
	}
	defer routerConn.Close()

	routerCli = metro.NewRouterClient(routerConn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := routerCli.Listen(ctx, &metro.ListenRequest{Token: token})
	if err != nil {
		log.Fatalf("Failed to listen from Router: %v", err)
	}

	for {
		sig, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Unexpected end of stream: %v", err)
		}

		signalHandler(sig)
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})

	if host := os.Getenv("LOKO_METRO_HOST"); len(host) > 0 {
		routerHost = host
	} else {
		routerHost = "0.0.0.0"
	}

	if port, err := strconv.ParseUint(os.Getenv("LOKO_METRO_PORT"), 10, 16); err == nil {
		routerPort = uint16(port)
	} else {
		routerPort = 50051
	}
}

func getRouterAddress() string {
	return routerHost + ":" + strconv.Itoa(int(routerPort))
}
