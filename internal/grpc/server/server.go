package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/daneofmanythings/diceroni/internal/grpc/proto"
	"google.golang.org/grpc"
)

var port int = 8080

type rollerServer struct {
	pb.UnimplementedRollerServer
}

func newServer() *rollerServer {
	return &rollerServer{}
}

func (s *rollerServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Ping: "pong"}, nil
}

func (s *rollerServer) Roll(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	return &pb.CreateResponse{
		Data:     nil,
		CallerID: "not yet implemented",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("...could not listen: %v", err)
	}

	var opts []grpc.ServerOption

	serviceRegistrar := grpc.NewServer(opts...)
	pb.RegisterRollerServer(serviceRegistrar, newServer())

	err = serviceRegistrar.Serve(lis)
	if err != nil {
		log.Fatalf("...could not serve: %v", err)
	}
}
