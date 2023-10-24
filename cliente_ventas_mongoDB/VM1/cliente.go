package main

import (
	"context"
	"encoding/json"
	"fmt"
	"grpc-golang/pb"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("Hello buyer ...")
	// Check if a JSON file name was provided as a command-line argument
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run client.go <json_file_name>")
		return
	}

	jsonFileName := os.Args[1]

	// Read order data from the specified JSON file
	data, err := os.ReadFile(jsonFileName)
	if err != nil {
		log.Fatal(err)
	}

	var request pb.OrderServiceRequest
	json.Unmarshal(data, &request)
	conn, err := grpc.Dial("10.10.10.221:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)
	resp, err := client.Order(context.Background(), &request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Receive response => %s ", resp.OrderResponse)
}
