package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/AlexDiru/grpc-course/calculator/calculatorpb"

	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	res := calculatorpb.SumResponse{
		Result: req.X + req.Y,
	}

	return &res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	k := (int64)(2)
	n := req.N

	for n > 1 {
		if n%k == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				Factor: k,
			})
			n /= k
		} else {
			k++
		}
	}

	return nil
}

func main() {
	fmt.Println("Calculator Service")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	grpcServer := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(grpcServer, &server{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
