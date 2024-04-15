package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	pb "github.com/daneofmanythings/calcuroller/internal/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	serverAddr string = "localhost:8080"
	caCertPath string = "./internal/grpc/certs/ca-cert.pem"
)

func main() {
	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("...could not load TLS credentials: ", err)
	}

	conn, err := grpc.Dial(
		serverAddr,
		grpc.WithTransportCredentials(tlsCredentials),
	)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	defer conn.Close()

	client := pb.NewRollerClient(conn)

	runREPL(client)
}

func runREPL(client pb.RollerClient) {
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

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}

func prettyStringifyMetadata(md []*pb.DiceRollMetadata) string {
	var out bytes.Buffer
	for _, data := range md {
		out.WriteString("\n{\n" + data.ResponseLiteral + ": ")
		out.WriteString(fmt.Sprintf("%d", data.Value) + "\n")
		// TODO: make the metadata print prettier. low priority.
		out.WriteString(data.String() + "\n}") // <- This is very ugly
	}
	out.WriteString("\n")

	return out.String()
}
