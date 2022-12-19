package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc/status"

	"google.golang.org/grpc/credentials/insecure"

	proto "proto/proto"

	"google.golang.org/grpc"
)

type Frontend struct {
	proto.UnimplementedFrontendServer
	serverConnection []proto.ServerClient
}

var frontendPort = flag.Int("port", 8000, "server port number")

func main() {
	logfile, err := os.OpenFile("../log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.SetOutput(logfile)
	log.SetFlags(2)
	flag.Parse()
	log.Println("Frontend: Frontend started")
	fmt.Println("Frontend: Frontend started")
	frontend := &Frontend{
		serverConnection: getServerConnection(),
	}
	go startFrontEnd(frontend)

	for {
		time.Sleep(100 * time.Second)
	}
}

func startFrontEnd(frontend *Frontend) {
	grpcServer := grpc.NewServer()
	lister, err := net.Listen("tcp", ":"+strconv.Itoa(*frontendPort))
	if err != nil {
		log.Fatalln("Could not start listener")
	}
	proto.RegisterFrontendServer(grpcServer, frontend)
	serverError := grpcServer.Serve(lister)
	if serverError != nil {
		log.Fatalln("Could not start server")
	}
}

func getServerConnection() []proto.ServerClient {
	conns := make([]proto.ServerClient, 3)
	for i := 0; i < 3; i++ {
		port := 8080 + i
		conn, err := grpc.Dial(":"+strconv.Itoa(port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Println("Could not dial server")
		}
		log.Printf("Frontend: Dialed server %d\n", port)
		fmt.Printf("Frontend: Dialed server %d\n", port)
		conns[i] = proto.NewServerClient(conn)
	}
	return conns
}

func (frontend *Frontend) Result(ctx context.Context, in *proto.Void) (*proto.BidResult, error) {
	counter := 0
	acks := make([]*proto.BidResult, 3)
	for i := 0; i < 3; i++ {
		ack, err := frontend.serverConnection[i].GetHighestBid(ctx, in)

		if err == nil {
			log.Printf("Frontend: Request result in node %d\n", 8080+i)
			fmt.Printf("Frontend: Request result in node %d\n", 8080+i)
			acks[counter] = ack
			counter++
		}
	}
	if len(acks) == 0 {
		return nil, status.Errorf(1, "error")
	}
	return acks[0], nil
}

func (frontend *Frontend) Bid(ctx context.Context, in *proto.BidRequest) (*proto.Ack, error) {
	counter := 0
	acks := make([]*proto.Ack, 3)
	log.Printf("Frontend: Received bid of %d by %s", in.Amount, in.Name)
	fmt.Printf("Frontend: Received bid of %d by %s", in.Amount, in.Name)
	for i := 0; i < 3; i++ {
		conn := frontend.serverConnection[i]
		ack, err := conn.UpdateHighestBid(ctx, in)

		if err == nil {
			log.Printf("Frontend: Sends updatehighest bid to node %d\n", 8080+i)
			fmt.Printf("Frontend: Sends updatehighest bid to node %d\n", 8080+i)
			acks[counter] = ack
			counter++
		}
	}
	if len(acks) == 0 {
		return nil, status.Errorf(1, "error")
	}
	return acks[0], nil
}

func (frontend *Frontend) StartAuction(ctx context.Context, in *proto.Void) (*proto.Ack, error) {
	counter := 0
	acks := make([]*proto.Ack, 3)
	for i := 0; i < 3; i++ {

		ack, err := frontend.serverConnection[i].StartAuction(ctx, in)
		if err == nil {
			log.Printf("Frontend: " + ack.Ack)
			fmt.Printf("Frontend: " + ack.Ack)
			acks[counter] = ack
			counter++
		}
	}
	if len(acks) == 0 {
		return nil, status.Errorf(1, "error")
	}
	return acks[0], nil
}
