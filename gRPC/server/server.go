package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/DarkLordOfDeadstiny/Mini-project-2/gRPC"

	"google.golang.org/grpc"
)

type Server struct {
	gRPC.UnimplementedMessageServiceServer
	messageChan []chan *gRPC.Message
}

func main() {
	fmt.Println(".:server is starting:.")

	// Create listener tcp on port 5400
	list, err := net.Listen("tcp", "localhost:5400")
	if err != nil {
		log.Fatalf("Failed to listen on port 5400: %v", err)
	}

	//Clears the log.txt file when a new server is started
	if err := os.Truncate("log.txt", 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}

	f, err := os.OpenFile("log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	
	log.SetOutput(f)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	gRPC.RegisterMessageServiceServer(grpcServer, newServer())
	grpcServer.Serve(list)

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}

func newServer() *Server {
	s := &Server{}
	fmt.Println(s)
	return s
}

func (s *Server) Join(ch *gRPC.JoinRequest, msgStream gRPC.MessageService_JoinServer) error {

	msgChannel := make(chan *gRPC.Message)
	s.messageChan = append(s.messageChan, msgChannel)

	s.sendMessage(&gRPC.Message{
		Message: ch.SendersName + " entered the chat",
		Sender: "Server",
	})

	// doing this never closes the stream
	for {
		select {
		case <-msgStream.Context().Done():
			s.sendMessage(&gRPC.Message{
					Message: ch.SendersName + " left the chat",
					Sender: "Server",
				})
			return nil
		case msg := <-msgChannel:
			fmt.Printf("GO ROUTINE (got message): %v \n", msg)
			msgStream.Send(msg)
		}
	}
}

func (s *Server) Send(msgStream gRPC.MessageService_SendServer) error {
	msg, err := msgStream.Recv()

	if err == io.EOF {
		return nil
	}

	if err != nil {
		return err
	}

	ack := gRPC.MessageAck{Status: "SENT"}
	msgStream.SendAndClose(&ack)

	go s.sendMessage(msg)

	return nil
}

func (s* Server) Leave(ctx context.Context, ch *gRPC.LeaveRequest) (*gRPC.LeaveResponse, error) {
	log.Printf("LEAVE_MESSAGE")

	response := gRPC.LeaveResponse{
		Status : "LEAVE_RESPONSE",
		}

	return &response, nil
}

func (s *Server) sendMessage(msg *gRPC.Message) {
	for _, msgChan := range s.messageChan {
		select {
		case msgChan <-msg:
		default:
		}
	}
	log.Printf("User: %s sent message \"%s\"",msg.Sender, msg.GetMessage())
}