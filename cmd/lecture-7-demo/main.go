package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"gitlab.com/siriusfreak/lecture-7-demo/internal/api"
	"gitlab.com/siriusfreak/lecture-7-demo/internal/ml_service"
	"google.golang.org/grpc"

	desc "gitlab.com/siriusfreak/lecture-7-demo/pkg/lecture-7-demo"
)

const (
	grpcPort = ":82"
	grpcServerEndpoint = "localhost:82"
)


func run() error {
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	desc.RegisterLecture7DemoServer(s, api.NewLecture7DemoAPI())

	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return nil
}

func runJSON() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := desc.RegisterLecture7DemoHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts)
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(":8081", mux)
	if err != nil {
		panic(err)
	}
}

func runMLService() {
	serv := ml_service.InitMLService()
	err := serv.StartConsuming()
	if err != nil {
		panic(err)
	}
}


func main() {
	go runJSON()
	go runMLService()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
