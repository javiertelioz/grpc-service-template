package examples

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/javiertelioz/grpc-service-template/examples/server"
	"github.com/javiertelioz/grpc-service-template/proto/helloworld/v1"
	"github.com/javiertelioz/grpc-service-template/proto/payments/v1"
)

// RunServerWithChi godoc
// Combine gRPC and REST in a single server on the same port using chi
func RunServerWithChi() error {
	grpcServer := grpc.NewServer()

	helloworld.RegisterGreeterServiceServer(grpcServer, &server.GreeterServer{})
	payments.RegisterPaymentServiceServer(grpcServer, &server.PaymentServer{})

	reflection.Register(grpcServer)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	gwMux := runtime.NewServeMux()
	ctx := context.Background()

	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	err = helloworld.RegisterGreeterServiceHandler(ctx, gwMux, conn)
	if err != nil {
		return err
	}

	err = payments.RegisterPaymentServiceHandler(ctx, gwMux, conn)
	if err != nil {
		log.Fatalf("failed to start PaymentServer HTTP gateway: %v", err)
	}

	r.Mount("/v1/", gwMux)

	httpServer := &http.Server{
		Addr: ":8080",
		Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.ProtoMajor == 2 && req.Header.Get("Content-Type") == "application/grpc" {
				grpcServer.ServeHTTP(w, req)
			} else {
				r.ServeHTTP(w, req)
			}
		}), &http2.Server{}),
	}

	fmt.Println("Starting gRPC and REST (with chi) services on :8080...")
	if err := httpServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
