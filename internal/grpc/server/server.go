package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"

	pb "github.com/daneofmanythings/calcuroller/internal/grpc/proto"
	"github.com/daneofmanythings/calcuroller/pkg/interpreter/object"
	"github.com/daneofmanythings/calcuroller/pkg/interpreter/repl"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
func (s *rollerServer) Roll(ctx context.Context, req *pb.RollRequest) (*pb.RollResponse, error) {
	result, metadata := repl.RunFromGRPC(req.GetDiceString())

	// this is pure chaos
	if result.Type() == object.ERROR_OBJ {
		return &pb.RollResponse{
			Message: &pb.RollResponse_Status{
				Status: &pb.MyStatus{
					Code:    int32(codes.InvalidArgument),
					Message: result.Inspect(),
				},
			},
		}, nil
	}

	metadataJSON, err := json.Marshal(metadata.Store)
	if err != nil {
		metadataJSON = []byte{}
	}

	// and this is pure chaos
	return &pb.RollResponse{
		Message: &pb.RollResponse_Data{
			Data: &pb.RollResponseData{
				Data: &pb.RollData{
					Literal:  result.Inspect(),
					Metadata: metadataJSON,
				},
				CallerId: req.GetCallerId(),
			},
		},
	}, nil
}

func main() {
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
