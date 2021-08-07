package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/AlexDiru/grpc-course/greet/greetpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
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

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet function was invoked with a streaming request")

	result := ""

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		firstName := req.GetGreeting().FirstName
		result += "Hello " + firstName + "! "
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("GreetEveryone function was invoked with a streaming request")

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}

		firstName := req.Greeting.FirstName

		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: "Hello " + firstName,
		})

		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", sendErr)
			return sendErr
		}

	}
}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Println("GreetWithDeadline function was invoked")

	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("The client cancelled the request")
			return nil, status.Error(codes.DeadlineExceeded, "The client cancelled the request")
		}
		time.Sleep(time.Second * 1)
	}

	firstName := req.GetGreeting().FirstName

	result := "Hello " + firstName

	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}

	return res, nil
}

func main() {
	fmt.Println("Hello")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	defer lis.Close()

	certFile := "ssl/server.crt"
	keyFile := "ssl/server.pem"
	creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)

	if sslErr != nil {
		log.Fatalf("Fatal error loading SSL certs")
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))

	greetpb.RegisterGreetServiceServer(grpcServer, &server{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
