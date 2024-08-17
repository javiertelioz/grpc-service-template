package examples

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/javiertelioz/grpc-service-template/examples/server"
	"github.com/javiertelioz/grpc-service-template/proto/helloworld/v1"
	"github.com/javiertelioz/grpc-service-template/proto/payments/v1"
)

// RunServerWithGin godoc
// Combine gRPC and REST in a single server on the same port using Gin
func RunServerWithGin() error {
	grpcServer := grpc.NewServer()

	helloworld.RegisterGreeterServiceServer(grpcServer, &server.GreeterServer{})
	payments.RegisterPaymentServiceServer(grpcServer, &server.PaymentServer{})

	reflection.Register(grpcServer)

	// Crear un router Gin para manejar las solicitudes REST
	r := gin.Default()

	gwMux := runtime.NewServeMux()
	ctx := context.Background()

	// Conectar al servidor gRPC en el mismo puerto
	conn, err := grpc.DialContext(ctx, "localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	// Registrar los handlers del gateway REST en el router Gin
	err = helloworld.RegisterGreeterServiceHandler(ctx, gwMux, conn)
	if err != nil {
		return err
	}

	err = payments.RegisterPaymentServiceHandler(ctx, gwMux, conn)
	if err != nil {
		log.Fatalf("failed to start PaymentServer HTTP gateway: %v", err)
	}

	// Montar el mux REST en el router Gin
	r.Any("/v1/*any", gin.WrapH(gwMux))

	// Crear un servidor HTTP2 compatible con h2c (HTTP/2 Cleartext)
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

	fmt.Println("Starting gRPC and REST (with Gin) services on :8080...")
	if err := httpServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
