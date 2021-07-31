package main

import (
	"context"
	"fmt"
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
