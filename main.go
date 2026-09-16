package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/SimonasKalpokas/todo_server/greeter_protos"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(context context.Context, helloRequest *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Message")
	return &pb.HelloReply{Message: fmt.Sprintf("Hello, %s", helloRequest.Name)}, nil
}

func main() {
	// setup a listener on port 9001
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a new grpc server
	grpcServer := grpc.NewServer()

	// register our server struct as a handle for the GreeterService rpc calls that come in through grpcServer
	pb.RegisterGreeterServer(grpcServer, &server{})

	// Serve traffic
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
