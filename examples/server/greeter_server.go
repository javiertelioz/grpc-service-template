package server

import (
	"context"
	"fmt"

	pb "github.com/javiertelioz/grpc-templates/proto/helloworld/v1"
)

type GreeterServer struct {
	pb.UnimplementedGreeterServiceServer
}

func (s *GreeterServer) SayHello(ctx context.Context, req *pb.GreeterServiceSayHelloRequest) (*pb.GreeterServiceSayHelloResponse, error) {
	return &pb.GreeterServiceSayHelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.Name),
	}, nil
}
