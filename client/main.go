package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/SimonasKalpokas/todo_server/greeter_protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Create a new grpc client
	conn, err := grpc.NewClient("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server at localhost:9001: %v", err)
	}
	// dont forget to close it
	defer conn.Close()

	// create a new coffee shop client from our generated code and pass in the connection created above
	c := pb.NewGreeterClient(conn)

	// give us a context that we can cancel, but also a timeout just to illustrate a point
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	helloRequest := pb.HelloRequest{
		Name: "World",
	}
	helloReply, err := c.SayHello(ctx, &helloRequest)

	fmt.Println(helloReply)
}
