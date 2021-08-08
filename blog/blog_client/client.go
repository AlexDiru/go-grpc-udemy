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

	// Create Blog
	createBlogRes, err := client.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "Alex",
			Title:    "My first blog",
			Content:  "Content of my first blog",
		},
	})

	if err != nil {
		log.Fatalf("Unexpected error \n\t%v\n", err)
	}

	fmt.Printf("Blog has been created \n\t%v\n", createBlogRes)

	// Read blog
	_, err2 := client.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: "frijirgnw",
	})

	if err2 != nil {
		fmt.Printf("Error happened while reading: \n\t%v\n", err2)
	}

	readBlogRes, err3 := client.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: createBlogRes.GetBlog().GetId(),
	})

	if err3 != nil {
		fmt.Printf("Error happened while reading: \n\t%v\n", err3)
	}

	fmt.Printf("Blog has been read:\n\t%v\n", readBlogRes)
}
