package payments

import (
	"context"
	"log"

	pbPaymentsService "github.com/javiertelioz/grpc-templates/proto/payments/v1"
)

// DepositService godoc
func DepositService(err error, paymentClient pbPaymentsService.PaymentServiceClient) {
	// Call to PaymentService.Deposit
	depositReq := &pbPaymentsService.DepositRequest{
		UserId: "user123",
		Amount: 100.0,
	}

	depositRes, err := paymentClient.Deposit(context.Background(), depositReq)
	if err != nil {
		log.Fatalf("Error calling Deposit: %v", err)
	}

	log.Printf("Deposit Response: %v", depositRes)
}
