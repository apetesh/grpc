package main

import (
	"log"
	"net"

	srv "github.com/apetesh/grpc/server"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	phoneBook
	grpcServer := grpc.NewServer()
	srv.RegisterGrpcServer(grpcServer, &srv.Server{})
	grpcServer.Serve(lis)
}
