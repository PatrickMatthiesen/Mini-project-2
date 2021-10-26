package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/DarkLordOfDeadstiny/Mini-project-2/gRPC"
	"google.golang.org/grpc"
)

var channelname = flag.String("Channel", "default", "Shitty-Chat")
var sendername = flag.String("sender", "default", "Senders name")
var tcpServer = flag.String("server", ":5400", "Tcp server")

func main() {
	flag.Parse()

	fmt.Println("--- CLIENT APP ---")

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithInsecure())

	conn, err := grpc.Dial(*tcpServer, opts...)
	if err != nil {
		log.Fatalf("something went wrong, big sad 👁👄👁 :c : %v", err)
	}

	defer conn.Close()

	ctx := context.Background()
	client := gRPC.NewMessageServiceClient(conn)

	fmt.Println("--- join channel ---")
	go joinChannel(ctx, client)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		go sendMessage(ctx, client, scanner.Text())
	}

}

func sendMessage(ctx context.Context, client gRPC.MessageServiceClient, message string) {
	stream, err := client.SendMessage(ctx)
	if err != nil {
		log.Printf("Unable to send message😡: error: %v", err)
	}
	msg := gRPC.Message{
		Channel: &gRPC.Channel{
			ChanName:    *channelname,
			SendersName: *sendername},
		Message: message,
		Sender:  *sendername,
	}
	stream.Send(&msg)

	ack, err := stream.CloseAndRecv()
	fmt.Printf("Message has been sent: %v \n", ack)

}

func joinChannel(ctx context.Context, client gRPC.MessageServiceClient) {

	channel := gRPC.Channel{ChanName: *channelname, SendersName: *sendername}
	stream, err := client.JoinChannel(ctx, &channel)
	if err != nil {
		log.Fatalf("client.JoinChannel(ctx, &channel) throws: %v", err)
	}

	fmt.Printf("Joined channel: %v \n", *&channelname)

	waitc := make(chan struct{})

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive message from channel joining. \nErr: %v", err)
			}

			if *sendername != in.Sender {
				fmt.Printf("MESSAGE: (%v) -> %v \n", in.Sender, in.Message)
			}
		}
	}()

	<-waitc
}
