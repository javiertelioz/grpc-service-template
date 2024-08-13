package payments

import (
	"context"
	"log"

	pbPaymentsService "github.com/javiertelioz/grpc-templates/proto/payments/v1"
)

// GetTransactionHistoryService godoc
func GetTransactionHistoryService(err error, paymentClient pbPaymentsService.PaymentServiceClient) {
	// Call to PaymentService.GetTransactionHistory
	historyReq := &pbPaymentsService.TransactionHistoryRequest{
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
