package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

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

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to the calcuroller client REPL!")
	fmt.Print("enter your name >> ")
	callerId, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\nan error occurred reading input. err=%s", err)
	}

	fmt.Print("(enter dice strings, ex: d4 + 1)\n\n")

	for {
		fmt.Print(">> ")
		diceString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("\nan error occurred reading input. err=%s", err)
			continue
		}
		response, err := client.Roll(context.Background(), &pb.RollRequest{
			DiceString: diceString,
			CallerId:   callerId,
		})
		if err != nil {
			log.Fatalf("Roll failed: err=%s", err)
		}
		switch response.Message.(type) {
		case *pb.RollResponse_Data:
			log.Println("Request string: " + response.GetData().GetRequestLiteral())
			log.Printf("Value: %d\n", response.GetData().GetValue())
			log.Println("Metadata:" + prettyStringifyMetadata(response.GetData().GetMetadata()))
		case *pb.RollResponse_Status:
			log.Println("(error) " + response.GetStatus().Message + "\n")
		}
	}
}

func prettyStringifyMetadata(md []*pb.DiceRollMetadata) string {
	var out bytes.Buffer
	for _, data := range md {
		out.WriteString("\n{\n" + data.ResponseLiteral + ": ")
		out.WriteString(fmt.Sprintf("%d", data.Value) + "\n")
		// TODO: make the metadata print prettier. low priority.
		out.WriteString(data.String() + "\n}")
	}
	out.WriteString("\n")

	return out.String()
}
