package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/silkeholmebonnen/disysExamPrep/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var name = flag.String("name", "DefaultName", "client name")
var id = flag.Int("id", 0, "id")
var frontendPort = flag.Int("frontendPort", 8000, "frontend port")

func main() {
	logfile, err := os.OpenFile("../log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.SetOutput(logfile)
	log.SetFlags(2)
	flag.Parse()
	go startClient()
	for {
		time.Sleep(100 * time.Second)
	}
}

func startClient() {
	serverConnection := getServerConnection()
	sendMessage(serverConnection)
}

func getServerConnection() proto.FrontendClient {
	conn, err := grpc.Dial(":"+strconv.Itoa(*frontendPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Could not dial server")
	}
	return proto.NewFrontendClient(conn)
}

func sendMessage(serverConnection proto.FrontendClient) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if strings.HasPrefix(input, "bid") {
			postNewBid(serverConnection, input, ctx)
		} else if input == "result" {
			logAuctionResult(serverConnection, ctx)
		} else if input == "start" {
			startAuction(serverConnection, ctx)
		} else {
			log.Println("Client: Unknown command: ", input)
		}
	}
}

func postNewBid(serverConnection proto.FrontendClient, input string, ctx context.Context) {
	split := strings.Split(input, " ")
	amount, _ := strconv.ParseInt(split[1], 10, 32)

	ans, err := serverConnection.Bid(ctx, &proto.BidRequest{
		Amount:    int32(amount),
		Name:      *name,
		ProcessID: int32(*id),
	})
	if err != nil {
		log.Fatal("Some error occured")
	} else {
		log.Println("Client: " + *name + " got response: " + ans.Ack)
	}
}

func logAuctionResult(serverConnection proto.FrontendClient, ctx context.Context) {
	ans, err := serverConnection.Result(ctx, &proto.Void{})
	if err != nil {
		log.Fatal("Some error occured")
	} else {
		log.Println("Client: " + ans.AuctionStatus + " The highest bid is " + strconv.Itoa(int(ans.Amount)) + " by " + ans.Name)
	}
}

func startAuction(serverConnection proto.FrontendClient, ctx context.Context) {
	ans, err := serverConnection.StartAuction(ctx, &proto.Void{})
	if err != nil {
		log.Fatal("Some error occured")
	} else {
		log.Println("Client: " + *name + " got response: " + ans.Ack)
	}
}
