package definition

import (
	pb "abouroumine.com/server/grpc-v2/service"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"log"
	"strings"
)

const (
	exampleScheme      = "example"
	exampleServiceName = "local.lo.com"
)

var addrs = []string{"localhost:50051", "localhost:50052"}

// The Server is used to Implement the generated/product
type Server struct {
	pb.UnimplementedProductInfoServer // This is used here to Implement Server in Registration
	pb.UnimplementedOrderManagementServer
	productMap map[string]*pb.Product
	orderMap   map[string]*pb.Order
}

// AddProduct implements the generated.AddProduct
func (s *Server) AddProduct(ctx context.Context, in *pb.Product) (*pb.ProductID, error) {
	// Providing the Generated ID using UUID Library.
	out, err := uuid.NewUUID()

	// Verifying if everything went well during the ID Generation.
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error While Generation PRODUCT ID", err)
	}

	// Inject the new ID generated to the Product 'in'
	in.Id = out.String()

	// Checking if the Server ProductMap is Empty and initialize if so.
	if s.productMap == nil {
		s.productMap = make(map[string]*pb.Product)
	}

	// Inject the Product to the Server ProductMap
	// The Key of the Product in the Server ProductMap will be its ID
	s.productMap[in.Id] = in

	// Return Address of the ProductID and error.
	return &pb.ProductID{Value: in.Id}, status.New(codes.OK, "").Err()
}

// GetProduct implements the generated/GetProduct
func (s *Server) GetProduct(ctx context.Context, in *pb.ProductID) (*pb.Product, error) {
	// We Inject the (Product,error) to value, exists from Server productMap using the ProductID as key
	value, exists := s.productMap[in.Value]

	// If found Return the Product
	if exists {
		return value, status.New(codes.OK, "").Err()
	}

	// Else return nil and Error Message.
	return nil, status.Errorf(codes.NotFound, "Product does not exist.", in.Value)
}

func (s *Server) GetProducts(ctx context.Context, empty *emptypb.Empty) (*pb.Products, error) {
	// Check if Server productMap empty
	if len(s.productMap) == 0 {
		return nil, status.Errorf(codes.NotFound, "Empty List.")
	}
	// We Create a Product Variable to store our data
	var products pb.Products

	// We loop over the Product Map to Inject all the products
	for _, v := range s.productMap {
		products.Products = append(products.Products, v)
	}

	// We return the Products and ok message.
	return &products, status.New(codes.OK, "").Err()
}

func (s *Server) GetOrder(ctx context.Context, in *wrapperspb.StringValue) (*pb.Order, error) {
	order, exists := s.orderMap[in.Value]
	if exists {
		return order, status.New(codes.OK, "").Err()
	}
	return nil, status.Errorf(codes.NotFound, "No Order with such Name...", in.Value)
}

func (s *Server) SearchOrders(value *wrapperspb.StringValue, stream pb.OrderManagement_SearchOrdersServer) error {
	// We start a Loop over the Order Map.
	for k, v := range s.orderMap {
		// We print each Order that exists.
		log.Print(k, v)
		// We Loop over the items of each Order.
		for _, itemStr := range v.Items {
			log.Print(itemStr)
			// We check the value of each Item we have in our Order.
			// In order to send the Order if item found.
			if strings.Contains(itemStr, value.Value) {
				err := stream.Send(v)
				if err != nil {
					return fmt.Errorf("Error Sending Message to Stream: %v", err)
				}
				log.Print("Matching Order Found: ", k)
				break
			}
		}
	}
	return nil
}

func (s *Server) UpdateOrder(stream pb.OrderManagement_UpdateOrderServer) error {
	orderStr := "Updated Order IDs: "
	// Check if Order Map is Empty.
	if len(s.orderMap) == 0 {
		s.orderMap = make(map[string]*pb.Order)
	}
	// We Loop Until no more received Streams.
	for {
		// Get the Order from the Stream.
		order, err := stream.Recv()
		// If no more Order received we close and break from loop & Return from Function.
		if err == io.EOF {
			return stream.SendAndClose(&wrappers.StringValue{Value: "Order Processed " + orderStr})
		}
		// We either Update Order or Add if not exist.
		s.orderMap[order.Id] = order
		log.Println("Order ID: ", order.Id, " has been Updated!")
		orderStr += order.Id + ", "
	}
}
