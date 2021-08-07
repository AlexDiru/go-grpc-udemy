package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/AlexDiru/grpc-course/calculator/calculatorpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Client started")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer conn.Close()

	client := calculatorpb.NewCalculatorServiceClient(conn)

	doSum(client)
	doPrimeNumberDecomposition(client)
	doAverage(client)
	doFindMaximum(client)
	doSquareRootError(client)

	fmt.Printf("Created client %f", client)
}

func doSum(client calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.SumRequest{
		X: 10,
		Y: 3,
	}
	res, err := client.Sum(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling Sum: %v", err)
	}

	log.Printf("Response from Sum: %v", res)
}

func doPrimeNumberDecomposition(client calculatorpb.CalculatorServiceClient) {

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		N: 120,
	}

	resStream, err := client.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling PrimeNumberDecomposition: %v", err)
	}

	for {
		msg, err := resStream.Recv()

		if err == io.EOF {
			// End of stream
			break
		} else if err != nil {
			log.Fatalf("Error while reading PrimeNumberDecomposition stream: %v", err)
		}

		log.Printf("Response from PrimeNumberDecomposition: %v", msg.GetFactor())
	}

}

func doAverage(client calculatorpb.CalculatorServiceClient) {

	stream, err := client.ComputeAverage(context.Background())

	if err != nil {
		log.Fatalf("Error while calling ComputeAverage: %v", err)
	}

	requests := []*calculatorpb.ComputeAverageRequest{
		{
			Value: 1,
		},
		{
			Value: 2,
		},
		{
			Value: 3,
		},
		{
			Value: 4,
		},
	}

	for _, req := range requests {
		stream.Send(req)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while receiving response from ComputeAverage: %v", err)
	}

	fmt.Printf("ComputeAverage response: %v\n", res)
}

func doFindMaximum(client calculatorpb.CalculatorServiceClient) {

	stream, err := client.FindMaximum(context.Background())

	if err != nil {
		log.Fatalf("Error while calling FindMaximum: %v", err)
		return
	}

	waitc := make(chan struct{})

	requests := []*calculatorpb.FindMaximumRequest{
		{
			Value: 1,
		},
		{
			Value: 5,
		},
		{
			Value: 3,
		},
		{
			Value: 6,
		},
		{
			Value: 2,
		},
		{
			Value: 20,
		},
	}

	// Send
	go func() {
		defer stream.CloseSend()

		for _, req := range requests {
			fmt.Printf("Sending message %v\n", req)
			stream.Send(req)
		}
	}()

	go func() {
		defer close(waitc)

		for {
			res, err := stream.Recv()

			if err == io.EOF {
				return
			} else if err != nil {
				log.Fatalf("Error while reading FindMaximum stream: %v", err)
				return
			}

			fmt.Printf("New Maximum found %v\n", res.Maximum)
		}
	}()

	<-waitc

}

func doSquareRootError(client calculatorpb.CalculatorServiceClient) {
	res, err := client.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{
		Number: -3,
	})

	if err != nil {
		resErr, ok := status.FromError(err)

		if ok {
			// Actual error from gRPC
			fmt.Println(resErr.Message())
			fmt.Println(resErr.Code())

			if resErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number")
			}
		} else {
			log.Fatalf("Custom Error while calling SquareRoot: %v", err)
		}

		return
	}

	fmt.Printf("Square root response: %v\n", res)
}
