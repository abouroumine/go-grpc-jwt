package main

import (
	pb "abouroumine.com/client/grpc-v2/service"
	"context"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"log"
	"time"
)

const (
	address  = "localhost:50051"
	hostname = "localhost"
	//crtFile       = "./cert/client.crt"
	//keyFile       = "./cert/client.key"
	serverCrtFile = "./cert/server.crt"
	//caFile        = "./cert/ca.crt"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile(serverCrtFile, hostname)
	if err != nil {
		log.Fatalf("Failed To Load Credentials: %v\n", err.Error())
	}

	auth := oauth.NewOauthAccess(fetchToken())

	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(auth),
		grpc.WithTransportCredentials(creds),
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("Did not connect: %v\n", err.Error())
	}

	defer conn.Close()
	c2 := pb.NewProductInfoClient(conn)

	clientDeadline := time.Now().Add(5 * time.Second)

	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
	defer cancel()

	name := "Apple iPhone 11"
	description := "Meet Apple iPhone 11. All-new dual-camera system with Ultra Wide and Night Mode."

	price := float32(1000.0)
	r, err := c2.AddProduct(ctx, &pb.Product{Name: name, Description: description, Price: price})

	if err != nil {
		log.Fatalf("Could not add Product: %v\n", err)
	}

	log.Printf("Product ID: %s added successfully\n", r.Value)

	product, err := c2.GetProduct(ctx, &pb.ProductID{Value: r.Value})
	if err != nil {
		log.Fatalf("Could not get Product: %v\n", err)
	}

	log.Printf("Product: %s\n\n\n\n", product.String())

}

func fetchToken() *oauth2.Token {
	return &oauth2.Token{AccessToken: "some-key"}
}
