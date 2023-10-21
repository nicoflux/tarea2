package main

import (
	"context"
	"encoding/json"
	"fmt"
	"grpc-golang/pb"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

type server struct {
	pb.OrderServiceServer
}

func (s *server) Order(ctx context.Context, req *pb.OrderServiceRequest) (*pb.OrderServiceReply, error) {
	receivedJSON, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	log.Printf("Reçu : %s", string(receivedJSON))
	orderNumber := 1234567890
	return &pb.OrderServiceReply{
		OrderResponse: fmt.Sprintf("You bought : " + req.Products[0].Title + " and your order number is " + strconv.Itoa(orderNumber)),
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	fmt.Println("Server is running on port 8080")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, &server{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
