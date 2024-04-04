package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	pb "github.com/daneofmanythings/calcuroller/internal/grpc/proto"
	"github.com/daneofmanythings/calcuroller/pkg/interpreter/object"
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
		if err == nil {
			response, err := client.Roll(context.Background(), &pb.CreateRequest{
				DiceString: diceString,
				CallerId:   callerId,
			})
			if err != nil {
				log.Fatalf("Roll failed: err=%s", err)
			}
			metadata := deserializeJSON(response.Data.Metadata)
			value := response.Data.Literal
			log.Println("CallerID: " + response.CallerId[:len(response.CallerId)-1])
			log.Println("Value: " + value)
			log.Println("Metadata:" + prettyStringifyMetadata(metadata))
		} else {
			fmt.Printf("\nan error occurred reading input. err=%s", err)
		}
	}
}

func deserializeJSON(input []byte) map[string]object.DiceData {
	reader := bytes.NewReader(input)
	container := make(map[string]object.DiceData)

	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&container)
	if err != nil {
		log.Fatalf("unable to deserialize json. err=%v", err)
	}

	return container
}

func prettyStringifyMetadata(md map[string]object.DiceData) string {
	var out bytes.Buffer
	for key, val := range md {
		out.WriteString("\n{\n" + key + ":\n")
		out.WriteString(val.Inspect() + "}")
	}
	out.WriteString("\n")

	return out.String()
}
