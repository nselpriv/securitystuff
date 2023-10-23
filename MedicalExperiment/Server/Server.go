package main

import (
	"context"
	"flag"
	proto "medic/Proto"
	"log"
	"net"
	"strconv"
	"google.golang.org/grpc"
)


type Server struct {
	proto.UnimplementedHospitalServer
	name string
	port int
}


var port = flag.Int("port", 0, "server port number")

func main () {

	flag.Parse()
	
	server := &Server{
		name: "Hospital",
		port: *port,
	}

	go startServer(server)

	for {
		
	}
}




func startServer (server *Server) {
	grpcServer := grpc.NewServer()

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

	log.Printf("client sent %s \n", in.Content)
	return &proto.ServerResponse{
		Success: false,
	}, nil
}