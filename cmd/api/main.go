package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	pbGreeter "github.com/javiertelioz/grpc-templates/proto/helloworld/v1"
	pbPayments "github.com/javiertelioz/grpc-templates/proto/payments/v1"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	paymentClient := pbPayments.NewPaymentServiceClient(conn)
	greeterClient := pbGreeter.NewGreeterServiceClient(conn)

	SayHelloService(err, greeterClient)              // Unary
	DepositService(err, paymentClient)               // Unary
	WithdrawService(err, paymentClient)              // Unary
	GetTransactionHistoryService(err, paymentClient) // Server Streaming
	UploadTransactionsService(err, paymentClient)    // Client Streaming
	RealTimeTransactionService(err, paymentClient)   // Bidirectional Streaming
}

// SayHelloService godoc
func SayHelloService(err error, greeterClient pbGreeter.GreeterServiceClient) {
	// Call to GreeterService.SayHello
	helloReq := &pbGreeter.GreeterServiceSayHelloRequest{
		Name: "Javier",
	}
	helloRes, err := greeterClient.SayHello(context.Background(), helloReq)
	if err != nil {
		log.Fatalf("Error calling SayHello: %v", err)
	}

	log.Printf("SayHello Response: %v", helloRes)
}

// WithdrawService godoc
func WithdrawService(err error, paymentClient pbPayments.PaymentServiceClient) {
	// Call to PaymentService.Withdraw
	withdrawReq := &pbPayments.WithdrawRequest{
		UserId: "user123",
		Amount: 50.0,
	}

	withdrawRes, err := paymentClient.Withdraw(context.Background(), withdrawReq)
	if err != nil {
		log.Fatalf("Error calling Withdraw: %v", err)
	}

	log.Printf("Withdraw Response: %v", withdrawRes)
}

// DepositService godoc
func DepositService(err error, paymentClient pbPayments.PaymentServiceClient) {
	// Call to PaymentService.Deposit
	depositReq := &pbPayments.DepositRequest{
		UserId: "user123",
		Amount: 100.0,
	}
	depositRes, err := paymentClient.Deposit(context.Background(), depositReq)
	if err != nil {
		log.Fatalf("Error calling Deposit: %v", err)
	}
	log.Printf("Deposit Response: %v", depositRes)
}

// GetTransactionHistoryService godoc
func GetTransactionHistoryService(err error, paymentClient pbPayments.PaymentServiceClient) {
	// Call to PaymentService.GetTransactionHistory
	historyReq := &pbPayments.TransactionHistoryRequest{
		UserId: "user123",
	}

	historyStream, err := paymentClient.GetTransactionHistory(context.Background(), historyReq)
	if err != nil {
		log.Fatalf("Error calling GetTransactionHistory: %v", err)
	}

	for {
		transaction, err := historyStream.Recv()
		if err != nil {
			break
		}
		log.Printf("Transaction: %v", transaction)
	}
}

// UploadTransactionsService godoc
func UploadTransactionsService(err error, paymentClient pbPayments.PaymentServiceClient) {
	// Call to PaymentService.UploadTransactions
	uploadStream, err := paymentClient.UploadTransactions(context.Background())
	if err != nil {
		log.Fatalf("Error calling UploadTransactions: %v", err)
	}
	for i := 0; i < 5; i++ {
		transaction := &pbPayments.Transaction{
			TransactionId: fmt.Sprintf("txn%d", i),
			UserId:        "user123",
			Amount:        float64(i * 10),
			Type:          "deposit",
			Status:        "SUCCESS",
			Timestamp:     time.Now().Format(time.RFC3339),
		}
		err := uploadStream.Send(&pbPayments.UploadTransactionsRequest{
			Transactions: []*pbPayments.Transaction{transaction},
		})
		if err != nil {
			log.Fatalf("Error sending transaction: %v", err)
		}
	}
	uploadRes, err := uploadStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving upload response: %v", err)
	}
	log.Printf("Upload Response: %v", uploadRes)
}

// RealTimeTransactionService godoc
func RealTimeTransactionService(err error, paymentClient pbPayments.PaymentServiceClient) {
	// Call to PaymentService.RealTimeTransaction
	realTimeStream, err := paymentClient.RealTimeTransaction(context.Background())
	if err != nil {
		log.Fatalf("Error calling RealTimeTransaction: %v", err)
	}
	go func() {
		for i := 0; i < 5; i++ {
			err := realTimeStream.Send(&pbPayments.Transaction{
				TransactionId: fmt.Sprintf("txn%d", i),
				UserId:        "user123",
				Amount:        float64(i * 20),
				Type:          "withdraw",
				Status:        "SUCCESS",
				Timestamp:     time.Now().Format(time.RFC3339),
			})
			if err != nil {
				log.Fatalf("Error sending transaction: %v", err)
			}
			time.Sleep(time.Second)
		}
		realTimeStream.CloseSend()
	}()
	go func() {
		for {
			transaction, err := realTimeStream.Recv()
			if err != nil {
				break
			}
			log.Printf("Real-Time Transaction: %v", transaction)
		}
	}()
	time.Sleep(time.Second * 10)
}
