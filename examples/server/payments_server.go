package server

import (
	"context"
	"io"

	pb "github.com/javiertelioz/grpc-service-template/proto/payments/v1"
)

type PaymentServer struct {
	pb.UnimplementedPaymentServiceServer
}

func (s *PaymentServer) Deposit(ctx context.Context, req *pb.DepositRequest) (*pb.DepositResponse, error) {
	// Implement your business logic for deposit here
	return &pb.DepositResponse{
		TransactionId: "12345",
		Status:        "SUCCESS",
		Message:       "Deposit successful",
	}, nil
}

func (s *PaymentServer) Withdraw(ctx context.Context, req *pb.WithdrawRequest) (*pb.WithdrawResponse, error) {
	// Implement your business logic for withdrawal here
	return &pb.WithdrawResponse{
		TransactionId: "67890",
		Status:        "SUCCESS",
		Message:       "Withdrawal successful",
	}, nil
}

func (s *PaymentServer) GetTransactionHistory(req *pb.TransactionHistoryRequest, stream pb.PaymentService_GetTransactionHistoryServer) error {
	// Implement your business logic to stream transaction history here
	for i := 0; i < 10; i++ {
		stream.Send(&pb.Transaction{
			TransactionId: "12345",
			UserId:        req.UserId,
			Amount:        100.0,
			Type:          "deposit",
			Status:        "SUCCESS",
			Timestamp:     "2024-08-06T00:00:00Z",
		})
	}
	return nil
}

func (s *PaymentServer) UploadTransactions(stream pb.PaymentService_UploadTransactionsServer) error {
	var successCount, failureCount int32
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.UploadTransactionsResponse{
				SuccessCount: successCount,
				FailureCount: failureCount,
			})
		}
		if err != nil {
			return err
		}

		// Implement your business logic to handle each transaction here
		successCount++
	}
}

func (s *PaymentServer) RealTimeTransaction(stream pb.PaymentService_RealTimeTransactionServer) error {
	for {
		transaction, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// Implement your business logic to handle each transaction in real-time here
		stream.Send(transaction)
	}
}
