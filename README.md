# disysExamPrep
The project is copied from MadsRoager/AuctionReplication assignment 5

## How to start the program

Open at least 1 client terminal, 1 frontend terminal and 3 server terminals.
In the server terminals, navigate to the server folder and type the following commands in the 3 different terminals:

- Server Terminal 1: go run . -port 8080
- Server Terminal 2: go run . -port 8081
- Server Terminal 3: go run . -port 8082

In the Client terminal navigate to the Client folder and type:

- go run . 

In the Frontend terminal, navigate to the Frontend folder, and type:

- go run .

If there is more than one frontend, then the command is followed by a port. E.g. go run . -port [insert unique port here]

And the Client is told what frontend to connect to with
go run . -frontendPort [insert frontend port here] -name [insert name of client]

In the Client terminal you can write the 3 following commands (case sensitive)

1. To start an auction: start
2. To bid: bid [insert any number] e.g. bid 12
3. To get the result: result

To simulate a crash you can write crash (case sensitive) in one of the Server terminals.

### Commands for renewing proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/proto.proto

go mod tidy
