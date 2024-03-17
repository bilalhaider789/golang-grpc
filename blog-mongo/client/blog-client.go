package main

import (
	"context"
	"fmt"

	"io"
	"log"

	blog "github.com/bilalhaider789/golang-grpc/blog-mongo/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)


func main() {
	tls := false
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	if tls {
		certFile := "ssl/ca.crt" 
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificate: %v", sslErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	client, err := grpc.Dial("localhost:8000", opts)
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}
	defer client.Close()

	clientService := blog.NewBlogServiceClient(client)

	// for creating a new blog uncomment this code
	// ----------------------------------------------
	newblog := &blog.Blog{
		UserId:  "4",
		Title:   "Blog 4",
		Content: "Blog 4 Content",
	}
	response, err := clientService.CreateBlog(context.Background(), &blog.CreateBlogRequest{Blog: newblog})
	if err != nil{
		fmt.Println("resonse error", err)
	}
	fmt.Println(response)
	//  ------------------------------------------------------

	// for reading a blog using userId uncomment this code
	// --------------------------------------------------
	// response, err2 := clientService.GetBlog(context.Background(), &blog.GetBlogRequest{UserId: "1"})
	// if err2 != nil {
	// 	fmt.Println("resonse error", err2)
	// } else {
	// 	fmt.Println(response)
	// }
	//  ------------------------------------------------------

	// for reading a blog using userId uncomment this code
	// --------------------------------------------------
	// newblog := &blog.Blog{
	// 	UserId:  "1",
	// 	Title:   "Blog 1 updated",
	// 	Content: "Blog 1 Content updated",
	// }
	// response, err := clientService.UpdateBlog(context.Background(), &blog.UpdateBlogRequest{Blog : newblog})
	// if err != nil {
	// 	fmt.Println("resonse error", err)
	// }else{
	// 	fmt.Println(response)
	// }
	//  ------------------------------------------------------

	// for deleting first blog using userId uncomment this code
	// --------------------------------------------------

	// response, err := clientService.DeleteBlog(context.Background(), &blog.DeleteBlogRequest{UserId : "1"})
	// if err != nil {
	// 	fmt.Println("resonse error", err)
	// }else{
	// 	fmt.Println("First blog deleted of userId : ",response)
	// }
	//  ------------------------------------------------------

	// for Streaming of all blogs uncomment this code
	// --------------------------------------------------

	// stream, err := clientService.StreamBlog(context.Background(), &blog.StreamBlogRequest{})
	// if err != nil {
	// 	fmt.Println("resonse error", err)
	// }
	// for {
	// 	res, err := stream.Recv()
	// 	if err == io.EOF {
	//         break
	//     }
	// 	if err != nil {
	// 		log.Fatalf("Something happened: %v", err)
	// 	}
	// 	fmt.Println("\n\n Blog Title : ",res.GetBlog().GetTitle())
	// 	fmt.Println("\n Blog Content : ",res.GetBlog().GetContent())

	// }
	//  ------------------------------------------------------

}
