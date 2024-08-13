package payments

import (
	"context"
	"log"

	"github.com/javiertelioz/grpc-templates/proto/payments/v1"
)

// WithdrawService godoc
func WithdrawService(err error, paymentClient payments.PaymentServiceClient) {
	// Call to PaymentService.Withdraw
	withdrawReq := &payments.WithdrawRequest{
		UserId: "user123",
		Amount: 50.0,
	}

	withdrawRes, err := paymentClient.Withdraw(context.Background(), withdrawReq)
	if err != nil {
		log.Fatalf("Error calling Withdraw: %v", err)
	}

	log.Printf("Withdraw Response: %v", withdrawRes)
}
