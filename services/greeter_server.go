package services

import (
	"context"
	"fmt"

	pb "github.com/javiertelioz/grpc-templates/proto/helloworld/v1"
)

type GreeterService struct {
	pb.UnimplementedGreeterServiceServer
}

func (s *GreeterService) SayHello(ctx context.Context, req *pb.GreeterServiceSayHelloRequest) (*pb.GreeterServiceSayHelloResponse, error) {
	return &pb.GreeterServiceSayHelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.Name),
	}, nil
}
