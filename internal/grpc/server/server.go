package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"

	pb "github.com/daneofmanythings/diceroni/internal/grpc/proto"
	"github.com/daneofmanythings/diceroni/pkg/interpreter/evaluator"
	"github.com/daneofmanythings/diceroni/pkg/interpreter/lexer"
	"github.com/daneofmanythings/diceroni/pkg/interpreter/parser"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

// TODO: DO REAL ERROR CHECKING ON BOTH RETURN VALUES AND JSON ENCODING!
func (s *rollerServer) Roll(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	l := lexer.New(req.DiceString)
	p := parser.New(l)
	program := p.ParseProgram()

	result, metadata := evaluator.EvalFromRequest(program)

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		metadataJSON = []byte{}
	}

	return &pb.CreateResponse{
		Data: &pb.RollData{
			Literal:  fmt.Sprintf("%d", result),
			Metadata: metadataJSON,
		},
		CallerId: req.CallerId,
	}, nil
}

func Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("...could not listen: %v", err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterRollerServer(grpcServer, newServer())
	reflection.Register(grpcServer)

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("...could not serve: %v", err)
	}
}
