package main

import (
	"context"
	"log"

	pb "github.com/daneofmanythings/calcuroller/internal/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var serverAddr string = "localhost:8080"

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	defer conn.Close()

	client := pb.NewRollerClient(conn)

	result, err := client.Roll(context.Background(), &pb.CreateRequest{
		DiceString: "d20 + 5",
		CallerId:   "Mary",
	})
	if err != nil {
		log.Fatalf("Roll failed: err=%s", err)
	}

	log.Println(result)
}
