package main

import (
	"google.golang.org/grpc"
	"log"

	"github.com/javiertelioz/grpc-service-template/examples/client/greeter"
	"github.com/javiertelioz/grpc-service-template/examples/client/payments"
	pbGreeter "github.com/javiertelioz/grpc-service-template/proto/helloworld/v1"
	pbPayments "github.com/javiertelioz/grpc-service-template/proto/payments/v1"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	greeterClient := pbGreeter.NewGreeterServiceClient(conn)
	paymentClient := pbPayments.NewPaymentServiceClient(conn)

	greeter.SayHelloService(err, greeterClient)               // Unary
	payments.DepositService(err, paymentClient)               // Unary
	payments.WithdrawService(err, paymentClient)              // Unary
	payments.GetTransactionHistoryService(err, paymentClient) // Server Streaming
	payments.UploadTransactionsService(err, paymentClient)    // Client Streaming
	payments.RealTimeTransactionService(err, paymentClient)   // Bidirectional Streaming
}
