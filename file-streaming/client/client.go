package main

import (
	"context"
	"fmt"
	proto "github.com/bilalhaider789/go-check/file-streaming/proto"
	"io"
	"log"
	"os"

	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client proto.StreamUploadClient

func main() {
	conn, err := grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client = proto.NewStreamUploadClient(conn)

	mb := 1024 * 1024 * 2
	uploadStramFile("./1GB.bin", mb)
}

func uploadStramFile(path string, batchSize int) {
	t := time.Now()
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	// setting up buffer size
	buf := make([]byte, batchSize)
	batchNumber := 1
	stream, err := client.Upload(context.TODO())
	if err != nil {
		panic(err)
	}

	for {
		num, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		chunk := buf[:num]

		if err := stream.Send(&proto.UploadRequest{FilePath: path, Chunks: chunk}); err != nil {
			fmt.Println(err)
			return
		}
		log.Printf("Sent - batch #%v - size - %v\n", batchNumber, len(chunk))
		batchNumber += 1
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(time.Since(t))
	log.Printf("Sent - %v bytes - %s\n", res.GetFileSize(), res.GetMessage())
}
