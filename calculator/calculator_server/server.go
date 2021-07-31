package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
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

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {

	inputs := []int64{}

	for {
		req, err := stream.Recv()

		if err == io.EOF {

			sum := int64(0)
			for _, val := range inputs {
				sum += val
			}

			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: float64(sum) / float64(len(inputs)),
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		inputs = append(inputs, req.Value)
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {

	max := int64(math.MinInt64)

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}

		if req.Value > max {
			max = req.Value

			sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
				Maximum: max,
			})

			if sendErr != nil {
				log.Fatalf("Error while sending data to client: %v", sendErr)
				return sendErr
			}
		}

	}
}

func main() {
	fmt.Println("Calculator Service")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	defer lis.Close()

	grpcServer := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(grpcServer, &server{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
