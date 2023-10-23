package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	proto "medic/Proto"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	proto.UnimplementedPersonServer
	id         int
	name 	 string
	ownPort int
	clients     map[int32]proto.PersonClient
	ctx         context.Context
}

var amount_of_peers = 3

type State int


var (
	cid = flag.Int("id", 0, "client ID")
	serverPort = flag.Int("sport", 0, "server port")

)

func main () {
	flag.Parse()

	cp := int(*cid)+5001

	var name string
	switch *cid {
	case 0:
		name = "Alice"
	case 1:
		name = "Bob"
	case 2:
		name = "Charlie"
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()


	client := &Client{
		name: name,
		id: *cid,
		ownPort: cp,
		ctx: ctx,
		clients: make(map[int32]proto.PersonClient),
	}

	// Create listener tcp on port ownPort
	go SetupPeerConnection(client)
	go SetupHospitalConnection(client)

	

	for {

	}

}

func SetupPeerConnection(client *Client) {

	//setup own connection 
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", client.ownPort))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterPersonServer(grpcServer, client)

	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v", err)
		}
	}()

	// Dial all other peers
	for i := 0; i < amount_of_peers; i++ {
		port := int32(5001) + int32(i)

		if port == int32(client.ownPort) {
			continue
		}

		var conn *grpc.ClientConn
		log.Printf("Trying to dial: %v\n", port)
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		defer conn.Close()
		c := proto.NewPersonClient(conn)
		client.clients[port] = c
	}
}


func SetupHospitalConnection(client *Client) {
	serverConnection, _ := connectToHospital()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		log.Printf("client asked for the with input %s\n", input)

		message, err := serverConnection.SendPersonalInfo(context.Background(), 
		&proto.PersonalInfo{Content: input})


		if err != nil {
			log.Printf("error is %s" , err.Error())
		} else {
			log.Printf("Server returned %t" , message.Success)
		}
	}

}


func connectToHospital() (proto.HospitalClient, error) {
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(*serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to port %d", *serverPort)
	} else {
		log.Printf("connected to the server at port %d \n", *serverPort)
	}
	return proto.NewHospitalClient(conn), nil
}