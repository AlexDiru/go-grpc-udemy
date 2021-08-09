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
	blogId := createBlogRes.GetBlog().GetId()
	ReadBlog(client, "frijirgnw")
	ReadBlog(client, blogId)

	// Update blog
	newBlog := &blogpb.Blog{
		Id:       blogId,
		AuthorId: "Alex (Edited)",
		Title:    "My first blog (Edited)",
		Content:  "Content of my first blog (Edited)",
	}

	UpdateBlog(client, newBlog)

}

func UpdateBlog(client blogpb.BlogServiceClient, blog *blogpb.Blog) {

	req := &blogpb.UpdateBlogRequest{
		Blog: blog,
	}

	updateRes, err := client.UpdateBlog(context.Background(), req)

	if err != nil {
		fmt.Printf("Error happened while updating: \n\t%v\n", err)
	}

	fmt.Printf("Blog was updated:\n\t%v\n", updateRes.GetBlog())

}

func ReadBlog(client blogpb.BlogServiceClient, blogId string) *blogpb.Blog {
	readBlogRes, err := client.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: blogId,
	})

	if err != nil {
		fmt.Printf("Error happened while reading: \n\t%v\n", err)
	}

	fmt.Printf("Blog has been read:\n\t%v\n", readBlogRes)

	return readBlogRes.GetBlog()
}
