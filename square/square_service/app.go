package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	sq "gprc"
)

// Declaring the service that will be exposed over RPC
type server struct {
	sq.UnimplementedSquareServiceServer
}

// Our exposed method
func (s server) GetSquare(ctx context.Context, in *sq.GetSquareRequest) (*sq.GetSquareResponse, error) {

	return &sq.GetSquareResponse{
		Result: in.Number * in.Number,
	}, nil

}

func run() {

	lis, err := net.Listen("tcp", ":1234")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	sq.RegisterSquareServiceServer(s, &server{})

	log.Println("Listening on port 1234")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
