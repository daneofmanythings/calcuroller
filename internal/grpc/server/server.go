package main

import (
	"context"
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
	requestLiteral := req.GetDiceString()
	result, metadata := repl.RunFromGRPC(requestLiteral)

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

	value := result.(*object.Integer).Value
	diceRollMetadata := []*pb.DiceRollMetadata{}

	for _, rollData := range metadata.Store {
		rollMetadata := &pb.DiceRollMetadata{
			ResponseLiteral: rollData.Literal,
			Tags:            rollData.Tags,
			RawRolls:        rollData.RawRolls,
			FinalRolls:      rollData.FinalRolls,
			Value:           rollData.Value,
		}
		diceRollMetadata = append(diceRollMetadata, rollMetadata)
	}

	// and this is pure chaos
	return &pb.RollResponse{
		Message: &pb.RollResponse_Data{
			Data: &pb.RollData{
				RequestLiteral: requestLiteral,
				Value:          value,
				Metadata:       diceRollMetadata,
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
