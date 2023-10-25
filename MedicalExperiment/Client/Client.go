package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	proto "medic/Proto"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Client struct {
	proto.UnimplementedPersonServer
	id         int
	name 	 string
	ownPort int
	clients     map[int32]proto.PersonClient
	hospitalConnection proto.HospitalClient
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
		hospitalConnection: nil,
	}

	
	// Create listener tcp on port ownPort
	go SetupPeerConnection(client)
	go SetupHospitalConnection(client)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		log.Printf("client sent %s\n", input)

		if input == "hospital" {
				message, err := client.hospitalConnection.SendPersonalInfo(context.Background(), 
			&proto.PersonalInfo{Content: input})
			if err != nil {
				log.Printf("error is %s" , err.Error())
			} else {
				log.Printf("Server returned %t" , message.Success)
			}
		} else {
			for _, c := range client.clients {
				conv, inputerror := strconv.Atoi(input)
				if inputerror != nil {
					fmt.Printf("You need to share a number with the other clients \n")
					break
				}
				message, err := c.Share(context.Background(), 
			&proto.ShareInfo{
				Id: int32(client.id),
				Timestamp: int32(conv),
				})
			if err != nil {
				log.Printf("error is %s" , err.Error())
			} else {
				log.Printf("Server returned %v" , message.Timestamp)
			}
			}
		}
	}


		
	for {

	}

}

func loadCertsServer() credentials.TransportCredentials {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
    if err != nil {
        log.Fatalf("Failed to load certificates: %v", err)
    }
	 // Create a gRPC server with TLS configuration
	creds := credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{cert}})

	return creds
}

func loadCertsClient() *x509.CertPool {
	certPool := x509.NewCertPool()
caCert, err := os.ReadFile("cert.pem")
if err != nil {
    log.Fatalf("Failed to load server certificate: %v", err)
}
certPool.AppendCertsFromPEM(caCert)
return certPool
}


func SetupPeerConnection(client *Client) {
	
	//setup own connection 
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", client.ownPort))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(loadCertsServer()))
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
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(loadCertsClient(), "")))
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		c := proto.NewPersonClient(conn)
		client.clients[port] = c
		log.Printf("Connected to: %v\n", port)
	}
}


func SetupHospitalConnection(client *Client) {

	conn, err := grpc.Dial(fmt.Sprintf(":%v", *serverPort), grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(loadCertsClient(), "")))
	if err != nil {
		log.Fatalf("Could not connect to port %d", *serverPort)
	} else {
		log.Printf("connected to the server at port %d \n", *serverPort)
	}
	client.hospitalConnection = proto.NewHospitalClient(conn)
}

func (c *Client) Share(ctx context.Context, in *proto.ShareInfo) (*proto.Reply, error) {

	log.Printf("client sent %d \n", in.Timestamp)
	return &proto.Reply{
		Id: in.Id,
		Timestamp: in.Timestamp,
	}, nil
}
