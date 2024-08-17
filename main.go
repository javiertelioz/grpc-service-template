package main

import (
	"log"

	"github.com/javiertelioz/grpc-service-template/examples"
)

func main() {
	// Run on different port
	// examples.SetupGRPCServer()

	// Run with same port
	// if err := examples.RunServerWithSamePort(); err != nil {
	//	log.Fatalf("Failed to run server: %v", err)
	// }

	// Run with same port on chi
	// if err := examples.RunServerWithChi(); err != nil {
	//	log.Fatalf("Failed to run server: %v", err)
	// }

	// Run with same port on gin
	if err := examples.RunServerWithGin(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
