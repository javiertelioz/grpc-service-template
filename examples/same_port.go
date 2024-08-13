package examples

import (
	"context"
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/javiertelioz/grpc-templates/examples/server"
	"github.com/javiertelioz/grpc-templates/proto/helloworld/v1"
	"github.com/javiertelioz/grpc-templates/proto/payments/v1"
)

// RunServerWithSamePort godoc
// Combine gRPC and REST in a single server on the same port
func RunServerWithSamePort() error {
	grpcServer := grpc.NewServer()

	helloworld.RegisterGreeterServiceServer(grpcServer, &server.GreeterServer{})
	payments.RegisterPaymentServiceServer(grpcServer, &server.PaymentServer{})

	reflection.Register(grpcServer)

	// Configure the HTTP multiplexer to handle both gRPC and REST
	mux := runtime.NewServeMux()
	ctx := context.Background()

	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	// Register REST gateway handlers in the multiplexer
	err = helloworld.RegisterGreeterServiceHandler(ctx, mux, conn)
	if err != nil {
		return err
	}

	err = payments.RegisterPaymentServiceHandler(ctx, mux, conn)
	if err != nil {
		log.Fatalf("failed to start PaymentServer HTTP gateway: %v", err)
	}

	httpServer := &http.Server{
		Addr: ":8080",
		Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && r.Header.Get("Content-Type") == "application/grpc" {
				grpcServer.ServeHTTP(w, r)
			} else {
				mux.ServeHTTP(w, r)
			}
		}), &http2.Server{}),
	}

	fmt.Println("Starting gRPC and gRPC-Gateway services on :8080...")
	if err := httpServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
