package main

import (
	"context"
	"fmt"
	"log"

	"github.com/AlexDiru/grpc-course/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Client started")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer conn.Close()

	client := blogpb.NewBlogServiceClient(conn)

	createBlogRes, err := client.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "Alex",
			Title:    "My first blog",
			Content:  "Content of my first blog",
		},
	})

	if err != nil {
		log.Fatalf("Unexpected error %v", err)
	}

	fmt.Printf("Blog has been created %v", createBlogRes)
}
