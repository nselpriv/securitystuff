package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	proto "medic/Proto"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)


type Server struct {
	proto.UnimplementedHospitalServer
	name string
	port int
	values []int
	check map [string]bool
}


var port = flag.Int("port", 0, "server port number")

func main () {

	flag.Parse()
	
	server := &Server{
		name: "Hospital",
		port: *port,
		values: make([]int, 3),
		check: map[string]bool{"Alice": false, "Bob": false,"Charlie":false},
	}

	go startServer(server)

	for {
		
	}
}




func startServer (server *Server) {
	grpcServer := grpc.NewServer(grpc.Creds(loadCerts()))
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(server.port))

	if err != nil {
		log.Fatalf("could not create the server %v", err)
	}
	log.Printf("Started server at port: %d\n", server.port)
	proto.RegisterHospitalServer(grpcServer, server)
	serveError := grpcServer.Serve(listener)

	if serveError != nil {
		log.Fatalf("could not serve listener")
	}
}

func (c *Server) SendPersonalInfo(ctx context.Context, in *proto.PersonalInfo) (*proto.ServerResponse, error) {
	var count int
	switch in.Name {
	case "Alice":
		count = 0 
	case "Bob":
		count = 1
	case "Charlie":
		count = 2
	}
	c.values[count] = int(in.Value)

	log.Printf("Received value from %s: %d\n", in.Name, in.Value)
	c.check[in.Name] = true

	for _, v := range c.check {
		if !v {
			return &proto.ServerResponse{
				Success: false,
			}, nil
		}
	}
	final := c.values[0] + c.values[1] + c.values[2]
	log.Printf("\n🔥🔥🔥🔥\nAll values received!\nFinal aggregate value is %v \n 🔥🔥🔥🔥", final)
	return &proto.ServerResponse{
		Success: true,
	}, nil
	
}

func loadCerts() credentials.TransportCredentials {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
    if err != nil {
        log.Fatalf("Failed to load certificates: %v", err)
    }
	// Create a gRPC server with TLS configuration
	creds := credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{cert}})
	return creds
}

