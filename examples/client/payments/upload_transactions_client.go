package payments

import (
	"context"
	"fmt"
	"log"
	"time"

	pbPaymentsService "github.com/javiertelioz/grpc-templates/proto/payments/v1"
)

// UploadTransactionsService godoc
func UploadTransactionsService(err error, paymentClient pbPaymentsService.PaymentServiceClient) {
	// Call to PaymentService.UploadTransactions
	uploadStream, err := paymentClient.UploadTransactions(context.Background())
	if err != nil {
		log.Fatalf("Error calling UploadTransactions: %v", err)
	}
	for i := 0; i < 5; i++ {
		transaction := &pbPaymentsService.Transaction{
			TransactionId: fmt.Sprintf("txn%d", i),
			UserId:        "user123",
			Amount:        float64(i * 10),
			Type:          "deposit",
			Status:        "SUCCESS",
			Timestamp:     time.Now().Format(time.RFC3339),
		}
		err := uploadStream.Send(&pbPaymentsService.UploadTransactionsRequest{
			Transactions: []*pbPaymentsService.Transaction{transaction},
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
