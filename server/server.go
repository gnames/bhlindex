package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/gnames/bhlindex/protob"
	"google.golang.org/grpc"
)

var version string

type bhlServer struct{}

func (bhlServer) Ver(ctx context.Context, void *protob.Void) (*protob.Version, error) {
	ver := protob.Version{Value: version}
	return &ver, nil
}

func Serve(port int, ver string) {
	version = ver
	srv := grpc.NewServer()
	var bhl bhlServer
	protob.RegisterBHLIndexServer(srv, bhl)
	portVal := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", portVal)
	if err != nil {
		log.Fatalf("could not listen to %s: %v", portVal, err)
	}
	log.Fatal(srv.Serve(l))

	time.Sleep(20 * time.Second)
}
