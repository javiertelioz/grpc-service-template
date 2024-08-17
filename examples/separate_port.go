package examples

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/javiertelioz/grpc-service-template/examples/server"
	"github.com/javiertelioz/grpc-service-template/proto/helloworld/v1"
	"github.com/javiertelioz/grpc-service-template/proto/payments/v1"
)

// RunGRPCServer
// Set up the gRPC server on port 8081 and serve requests indefinitely
func RunGRPCServer() error {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	helloworld.RegisterGreeterServiceServer(s, &server.GreeterServer{})
	payments.RegisterPaymentServiceServer(s, &server.PaymentServer{})

	reflection.Register(s)

	fmt.Println("Starting gRPC services on :8081...")
	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}

// RunRESTServer godoc
// Set up the REST server on port 8080 and handle requests by proxying them to the gRPC server
func RunRESTServer() error {
	ctx := context.Background()
	mux := runtime.NewServeMux()

	conn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	//conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	if err := helloworld.RegisterGreeterServiceHandler(ctx, mux, conn); err != nil {
		return err
	}

	err = payments.RegisterPaymentServiceHandler(ctx, mux, conn)
	if err != nil {
		log.Fatalf("failed to start PaymentService HTTP gateway: %v", err)
	}

	fmt.Println("Starting gRPC-Gateway services on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		return err
	}
	return nil
}

func SetupGRPCServer() {
	go func() {
		if err := RunRESTServer(); err != nil {
			log.Fatalf("Failed to run REST services: %v", err)
		}
	}()

	if err := RunGRPCServer(); err != nil {
		log.Fatalf("Failed to run gRPC services: %v", err)
	}
}
