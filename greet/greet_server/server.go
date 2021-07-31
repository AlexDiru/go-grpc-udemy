package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/AlexDiru/grpc-course/greet/greetpb"

	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Println("Greet function was invoked with %v", req)

	firstName := req.GetGreeting().FirstName

	result := "Hello " + firstName

	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Println("GreetManyTimes function was invoked with %v", req)

	firstName := req.GetGreeting().FirstName

	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyTimesResponse{
			Result: "Hello " + firstName + " number " + strconv.Itoa(i),
		}

		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}

	return nil
}

func main() {
	fmt.Println("Hello")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	grpcServer := grpc.NewServer()

	greetpb.RegisterGreetServiceServer(grpcServer, &server{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
