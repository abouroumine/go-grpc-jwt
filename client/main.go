package main

import (
	i "abouroumine.com/client/grpc-v2/interceptors"
	pb "abouroumine.com/client/grpc-v2/service"
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(i.OrderUnaryClientInterceptor), grpc.WithStreamInterceptor(i.ClientStreamInterceptor))
	if err != nil {
		log.Fatalf("Did not connect: %v\n", err.Error())
	}

	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := pb.NewOrderManagementClient(conn)

	updateStream, err := c.UpdateOrder(ctx)
	if err != nil {
		log.Fatalf("%v.UpdateOrders(_) = _, %v\n", c, err)
	}

	updateOrder := pb.Order{
		Id:          "111",
		Description: "description 1",
	}
	updateOrder2 := pb.Order{
		Id:          "222",
		Description: "description 1",
	}

	if err := updateStream.Send(&updateOrder); err != nil {
		log.Fatalf("%v.Send(%v) = %v\n", updateStream, updateOrder, err)
	}

	if err := updateStream.Send(&updateOrder2); err != nil {
		log.Fatalf("%v.Send(%v) = %v\n", updateStream, updateOrder2, err)
	}

	updateRes, err := updateStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got Error %v, want %v\n", updateStream, err, nil)
	}
	log.Printf("Update Orders Res: %s\n", updateRes)

	/*
		c := pb.NewOrderManagementClient(conn)

		searchStream, _ := c.SearchOrders(ctx, &wrappers.StringValue{Value: "Google"})

		for {
			searchOrder, err := searchStream.Recv()
			if err == io.EOF {
				break
			}
			log.Println("Search Result: ", searchOrder)
		}*/

	/*c := pb.NewProductInfoClient(conn)


	name := "Apple iPhone 11"
	description := "Meet Apple iPhone 11. All-new dual-camera system with Ultra Wide and Night Mode."

	price := float32(1000.0)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, err := c.AddProduct(ctx, &pb.Product{Name: name, Description: description, Price: price})

	if err != nil {
		log.Fatalf("Could not add Product: %v\n", err)
	}

	log.Printf("Product ID: %s added successfully\n", r.Value)

	product, err := c.GetProduct(ctx, &pb.ProductID{Value: r.Value})

	if err != nil {
		log.Fatalf("Could not get Product: %v\n", err)
	}

	log.Printf("Product: %s\n\n\n\n", product.String())

	time.Sleep(2 * time.Second)

	products, err := c.GetProducts(ctx, &emptypb.Empty{})

	if err != nil {
		log.Fatalf("Could not get Product List: %v\n", err.Error())
	}

	log.Printf("Products List size is: %v \nThe Content is: %s\n", len(products.Products), products.String())
	*/
}
