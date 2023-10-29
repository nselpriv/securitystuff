package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"math/rand"
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
	shares map[string]int
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
		shares: map[string]int{"Alice": -9999, "Bob": -9999,"Charlie":-9999},
	}


	
	// Create listener tcp on port ownPort
	go SetupPeerConnection(client)
	go SetupHospitalConnection(client)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		log.Printf("client input:  %s\n", input)

		if input == "hospital" {
			sum := getSumOfShares(client.shares)
				message, err := client.hospitalConnection.SendPersonalInfo(context.Background(), 
			&proto.PersonalInfo{
				Name: client.name,
				Value: sum,
			})
			if err != nil {
				log.Printf("error is %s" , err.Error())
			} else {
				if (message.Success) {
					log.Print("Hospital returned: Success!\n all values received and verified :)\n")
				} else {
					log.Print("Hospital returned: still waiting for values\n")
				}
			}
		} else {
			
			conv, inputerror := strconv.Atoi(input)
			if inputerror != nil {
				fmt.Printf("You need to share a number with the other clients \n")
				continue
			}
			n1,n2,n3 := MPCScramble(conv)
			client.shares[client.name] = n1  //setting own secret share
			log.Printf("Scrables from %v are first %v second %v third %v \n", conv, n1,n2,n3)
			count := 0 
			for _, c := range client.clients {
				var share int64
				if(count == 0){
					share = int64(n2)
				} else {share = int64(n3)}
				message, err := c.Share(context.Background(), 
			&proto.ShareInfo{
				Share: share,
				Name: client.name,
			})
			if err != nil {
				log.Printf("error is %s" , err.Error())
			} else {
				log.Printf("%s\n" ,message.Status)
			}
			count++
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

	log.Printf("%s sent %d \n", in.Name, in.Share)
	c.shares[in.Name] = int(in.Share)
	v := strconv.Itoa(int(in.Share))
	msg := c.name + " set " + v + " for " + in.Name
		return &proto.Reply{
		Status: msg,
	}, nil
}

func MPCScramble (number int) (first,second,third int) {
first = rand.Intn(number-3)
if(randBool()){
	first = first *-1
}

for {
gen := rand.Intn(number)
if(gen + first < number){
	second = gen 
	if(randBool()){
		second=second*-1
		}
	third = number - first - second 
	return
	}
}

}

//generates a random bool 
func randBool() bool{
return rand.Intn(2) == 0
}

func getSumOfShares(m map[string]int) (sum int64){
	sum = int64((m["Alice"]) + (m["Bob"]) + (m["Charlie"]))
	log.Printf("Summing values from:\n Alice: %d\n Bob: %d\n Charlie: %d \n With sum %v \n", m["Alice"], m["Bob"], m["Charlie"], sum)
	return
}