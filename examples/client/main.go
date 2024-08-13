package main

import (
	"github.com/javiertelioz/grpc-templates/examples/client/greeter"
	"github.com/javiertelioz/grpc-templates/examples/client/payments"
	"google.golang.org/grpc"
	"log"

	pbGreeter "github.com/javiertelioz/grpc-templates/proto/helloworld/v1"
	pbPayments "github.com/javiertelioz/grpc-templates/proto/payments/v1"
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
