package main

import (
	d "abouroumine.com/server/grpc-v2/definition"
	pb "abouroumine.com/server/grpc-v2/service"
	"context"
	"crypto/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"strings"
)

var (
	port    = ":50051"
	crtFile = "./cert/server.crt"
	keyFile = "./cert/server.key"
	// caFile             = "./cert/ca.crt"
	errMissingMetaData = status.Errorf(codes.InvalidArgument, "Missing MetaData")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "Invalid Credentials")
)

func main() {
	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("Failed To Load Key Pairs: %s\n", err.Error())
	}
	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		grpc.UnaryInterceptor(ensureValidToken),
	}

	s := grpc.NewServer(opts...)

	// We used &d.Server{} since we did add the UnimplementedProductInfoServer in the Server Structure.
	pb.RegisterProductInfoServer(s, &d.Server{})

	// We initialize our TCP Server on port 50051 to start listening.
	lis, err := net.Listen("tcp", port)

	// if there is a problem we cancel our action and abandon the listening.
	if err != nil {
		log.Fatalf("Failed to Listen to Server: %v", err.Error())
	}

	log.Printf("Starting gRPC listener on port " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to Serve: %v", err)
	}

}

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	return token == "some-key"
}

func ensureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetaData
	}
	if !valid(md["authorization"]) {
		return nil, errInvalidToken
	}
	return handler(ctx, req)
}
