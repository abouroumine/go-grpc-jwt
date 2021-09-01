package main

import (
	d "abouroumine.com/server/grpc-v2/definition"
	pb "abouroumine.com/server/grpc-v2/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":50051"
)

func main() {
	// We initialize our TCP Server on port 50051 to start listening.
	lis, err := net.Listen("tcp", port)

	// if there is a problem we cancel our action and abandon the listening.
	if err != nil {
		log.Fatalf("Failed to Listen to Server: %v", err.Error())
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(d.OrderUnaryServerInterceptor), grpc.StreamInterceptor(d.OrderServerStreamInterceptor))

	// We used &d.Server{} since we did add the UnimplementedProductInfoServer in the Server Structure.
	pb.RegisterProductInfoServer(s, &d.Server{})
	pb.RegisterOrderManagementServer(s, &d.Server{})

	log.Printf("Starting gRPC listener on port " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to Serve: %v", err)
	}

}
