package main

import (
	"flag"
	"log"
	"net"

	srv "github.com/liyunghao/Online-Eletronic-Voting/internal/server/services"
	pb "github.com/liyunghao/Online-Eletronic-Voting/internal/voting"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8080, "Specify which port should gRPC server listen on")
)

func main() {
	flag.Parse()
	tcp_listner, err := net.ListenTCP("tcp", &net.TCPAddr{IP: nil, Port: *port})
	if err != nil {
		log.Fatalf("Create TCP listner failed. Something WRONG: %v\n", err)
	}

	// EVoting Service
	var eVotingSrv srv.Service_eVoting

	grpcServer := grpc.NewServer()
	pb.RegisterEVotingServer(grpcServer, &eVotingSrv)

	// Start Server
	log.Printf("gRPC server start to listen at %d\n", *port)
	err = grpcServer.Serve(tcp_listner)
	if err != nil {
		log.Fatalf("Create TCP listner failed. Something WRONG: %v\n", err)
	}
}
