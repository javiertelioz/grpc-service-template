package greeter

import (
	"context"
	"log"

	"github.com/javiertelioz/grpc-templates/proto/helloworld/v1"
)

// SayHelloService godoc
func SayHelloService(err error, greeterClient helloworld.GreeterServiceClient) {
	// Call to GreeterService.SayHello
	helloReq := &helloworld.GreeterServiceSayHelloRequest{
		Name: "Javier",
	}

	helloRes, err := greeterClient.SayHello(context.Background(), helloReq)

	if err != nil {
		log.Fatalf("Error calling SayHello: %v", err)
	}

	log.Printf("SayHello Response: %v", helloRes)
}
