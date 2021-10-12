package gRPC

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

func main() {
	// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
	defer conn.Close()

	//  Create new Client from generated gRPC code from proto
	c := gRPC.NewGetXXXClient(conn)

	SendRequest(c)
}

func SendRequest(c gRPC.XXXClient) {
	// Between the curly brackets are nothing, because the .proto file expects no input.
	message := gRPC.somethingsomethingRequest{}

	response, err := c.somethingsomething(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling XXX: %s", err)
	}

	fmt.Printf("Response from the Server: %s \n", response.Reply)
}
