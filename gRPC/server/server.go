package gRPC

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/DarkLordOfDeadstiny/Mini-project-2/gRPC"

	"google.golang.org/grpc"
)

type Server struct {
	gRPC.UnimplementedGetXXXServer
}

func (s *Server) GetTime(ctx context.Context, in *gRPC.GetXXXRequest) (*gRPC.GetXXXReply, error) {
	fmt.Printf("Received XXX request")
	return &gRPC.GetXXXReply{Reply: "Your reply here"}, nil
}

func main() {
	// Create listener tcp on port 9080
	list, err := net.Listen("tcp", ":9080")
	if err != nil {
		log.Fatalf("Failed to listen on port 9080: %v", err)
	}
	grpcServer := grpc.NewServer()
	gRPC.RegisterXXXServer(grpcServer, &Server{})

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}
