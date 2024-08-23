package main

import (
	"log"
	"net"

	example "grpc-server/proto" // Importe o pacote gerado pelo compilador protobuf

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	example.RegisterExampleServiceServer(s, &example.MyServer{})

	log.Println("Server started at :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
