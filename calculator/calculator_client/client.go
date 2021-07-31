package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/AlexDiru/grpc-course/calculator/calculatorpb"

	"google.golang.org/grpc"
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
