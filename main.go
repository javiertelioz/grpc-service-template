package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/javiertelioz/grpc-templates/services"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pbGreeterService "github.com/javiertelioz/grpc-templates/proto/helloworld/v1"
	pbPaymentsService "github.com/javiertelioz/grpc-templates/proto/payments/v1"
)

// Set up the gRPC server on port 8080 and serve requests indefinitely
func runGRPCServer() error {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pbGreeterService.RegisterGreeterServiceServer(s, &services.GreeterService{})
	pbPaymentsService.RegisterPaymentServiceServer(s, &services.PaymentService{})

	reflection.Register(s)

	fmt.Println("Starting gRPC services on :8080...")
	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}

// Set up the REST server on port 8081 and handle requests by proxying them to the gRPC server
func runRESTServer() error {
	ctx := context.Background()
	mux := runtime.NewServeMux()

	conn, err := grpc.DialContext(ctx, "localhost:8080", grpc.WithInsecure())
	if err != nil {
		return err
	}

	if err := pbGreeterService.RegisterGreeterServiceHandler(ctx, mux, conn); err != nil {
		return err
	}

	err = pbPaymentsService.RegisterPaymentServiceHandler(ctx, mux, conn)
	if err != nil {
		log.Fatalf("failed to start PaymentService HTTP gateway: %v", err)
	}

	fmt.Println("Starting gRPC-Gateway services on :8081...")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		return err
	}
	return nil
}

func main() {
	go func() {
		if err := runRESTServer(); err != nil {
			log.Fatalf("Failed to run REST services: %v", err)
		}
	}()

	if err := runGRPCServer(); err != nil {
		log.Fatalf("Failed to run gRPC services: %v", err)
	}
}
