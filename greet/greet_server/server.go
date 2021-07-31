package main

import (
	"fmt"
	"log"
	"net"

	"github.com/AlexDiru/grpc-course/greet/greetpb"

	"google.golang.org/grpc"
)

type server struct{}

func main() {
	fmt.Println("Hello")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	server := grpc.NewServer()
	gss := greetpb.UnimplementedGreetServiceServer{}
	greetpb.RegisterGreetServiceServer(server, gss)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
