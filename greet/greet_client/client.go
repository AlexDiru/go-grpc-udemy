package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/AlexDiru/grpc-course/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Client started")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer conn.Close()

	client := greetpb.NewGreetServiceClient(conn)

	// doUnary(client)
	// doServerStreaming(client)
	// doClientStreaming(client)
	doBiDiStreaming(client)
}

func doUnary(client greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Alex",
			LastName:  "Spedding",
		},
	}
	res, err := client.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling Greet: %v", err)
	}

	log.Printf("Response from Greet: %v", res)
}

func doServerStreaming(client greetpb.GreetServiceClient) {

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Alex",
			LastName:  "Spedding",
		},
	}

	resStream, err := client.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes: %v", err)
	}

	for {
		msg, err := resStream.Recv()

		if err == io.EOF {
			// End of stream
			break
		} else if err != nil {
			log.Fatalf("Error while reading GreetManyTimes stream: %v", err)
		}

		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}

func doClientStreaming(client greetpb.GreetServiceClient) {

	stream, err := client.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("Error while calling LongGreet: %v", err)
	}

	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Alex",
				LastName:  "Spedding",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Tom",
				LastName:  "Jones",
			},
		},
	}

	for _, req := range requests {
		stream.Send(req)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while receiving response from LongGreet: %v", err)
	}

	fmt.Printf("LongGreet response: %v\n", res)
}

func doBiDiStreaming(client greetpb.GreetServiceClient) {

	stream, err := client.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("Error while calling GreetEveryone: %v", err)
		return
	}

	waitc := make(chan struct{})

	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Alex",
				LastName:  "Spedding",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Tom",
				LastName:  "Jones",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
				LastName:  "Carmack",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
				LastName:  "Romero",
			},
		},
	}

	// Send
	go func() {
		defer stream.CloseSend()

		for _, req := range requests {
			fmt.Printf("Sending message %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	go func() {
		defer close(waitc)

		for {
			res, err := stream.Recv()

			if err == io.EOF {
				return
			} else if err != nil {
				log.Fatalf("Error while reading GreetManyTimes stream: %v", err)
				return
			}

			fmt.Printf("Received %v\n", res.Result)
		}
	}()

	<-waitc

}
