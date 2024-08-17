package payments

import (
	"context"
	"fmt"
	"log"
	"time"

	pbPaymentsService "github.com/javiertelioz/grpc-service-template/proto/payments/v1"
)

// RealTimeTransactionService godoc
func RealTimeTransactionService(err error, paymentClient pbPaymentsService.PaymentServiceClient) {
	// Call to PaymentService.RealTimeTransaction
	realTimeStream, err := paymentClient.RealTimeTransaction(context.Background())
	if err != nil {
		log.Fatalf("Error calling RealTimeTransaction: %v", err)
	}

	go func() {
		for i := 0; i < 5; i++ {
			err := realTimeStream.Send(&pbPaymentsService.Transaction{
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
