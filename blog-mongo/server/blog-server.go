package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	blog "github.com/bilalhaider789/go-check/blog-mongo/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var collection *mongo.Collection

type server struct {
	blog.UnimplementedBlogServiceServer
}

type blogItem struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	UserID  string             `bson:"userId"`
	Title   string             `bson:"title"`
	Content string             `bson:"content"`
}

func (s *server) CreateBlog(ctx context.Context, req *blog.CreateBlogRequest) (*blog.CreateBlogResponse, error) {
	reqBlog := req.GetBlog()
	data := blogItem{
		UserID:  reqBlog.GetUserId(),
		Title:   reqBlog.GetTitle(),
		Content: reqBlog.GetContent(),
	}
	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		fmt.Println(err)
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		fmt.Println("something wrong")
	}

	return &blog.CreateBlogResponse{
		Blog: &blog.Blog{
			Id:      oid.Hex(),
			UserId:  reqBlog.GetUserId(),
			Title:   reqBlog.GetTitle(),
			Content: reqBlog.GetContent(),
		},
	}, nil
}

func (s *server) UpdateBlog(ctx context.Context, req *blog.UpdateBlogRequest) (*blog.UpdateBlogResponse, error) {
	reqBlog := req.GetBlog()
	userId := reqBlog.GetUserId()

	_, err := collection.UpdateOne(context.Background(), bson.D{{Key: "userId", Value: userId}}, bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "title", Value: reqBlog.GetTitle()},
			{Key: "content", Value: reqBlog.GetContent()},
		}},
	})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res := collection.FindOne(context.Background(), bson.D{{Key: "userId", Value: userId}})
	data := &blogItem{}
	if err := res.Decode(data); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &blog.UpdateBlogResponse{
		Blog: &blog.Blog{
			Id:      data.ID.Hex(),
			UserId:  data.UserID,
			Title:   data.Title,
			Content: data.Content,
		},
	}, nil
}

func (s *server) GetBlog(ctx context.Context, req *blog.GetBlogRequest) (*blog.GetBlogResponse, error) {
	reqBlog := req.GetUserId()
	fmt.Println(reqBlog)
	filter := bson.D{{Key: "userId", Value: reqBlog}}
	res := collection.FindOne(context.Background(), filter)
	data := &blogItem{}
	if err := res.Decode(data); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &blog.GetBlogResponse{
		Blog: &blog.Blog{
			Id:      data.ID.Hex(),
			UserId:  data.UserID,
			Title:   data.Title,
			Content: data.Content,
		},
	}, nil
}

func (s *server) DeleteBlog(ctx context.Context, req *blog.DeleteBlogRequest) (*blog.DeleteBlogResponse, error) {
	userId := req.GetUserId()
	res := collection.FindOne(context.TODO(), bson.D{{Key: "userId", Value: userId}})
	if res.Err() != nil {
		return nil, res.Err()
	}
	return &blog.DeleteBlogResponse{UserId: userId}, nil
}

func (s *server) StreamBlog(_ *blog.StreamBlogRequest, stream blog.BlogService_StreamBlogServer) error {

	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		data := &blogItem{}
		if err := cursor.Decode(data); err != nil {
			return err
		}
		fmt.Println(data)
		stream.Send(&blog.StreamBlogResponse{Blog: &blog.Blog{
			Id:      data.ID.Hex(),
			UserId:  data.UserID,
			Title:   data.Title,
			Content: data.Content,
		}})
	}
	return nil
}

func main() {

	listener, portErr := net.Listen("tcp", ":8000")
	if portErr != nil {
		fmt.Println(portErr)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	client, dbErr := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:secret@localhost:27017"))
	if dbErr != nil {
		log.Fatalf("failed to connect to db: %v", dbErr)
	} else {
		fmt.Println("connected to db")
	}

	collection = client.Database("grpcBlog").Collection("blogs")

	opts := []grpc.ServerOption{}
	tls := false
	if tls {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Failed loading certificates: %v", sslErr)
			return
		}
		opts = append(opts, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(opts...)

	blog.RegisterBlogServiceServer(grpcServer, &server{})

	go func() {
		fmt.Println("server started")
		serverErr := grpcServer.Serve((listener))
		if serverErr != nil {
			log.Fatalf("failed to serve: %v", serverErr)
		}
	}()

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	<-channel
	fmt.Println("Server shutdown")
	grpcServer.Stop()
	fmt.Println("Stopping listner")
	listener.Close()
	client.Disconnect(context.TODO())
	fmt.Println("Normal Shutdown")

}
