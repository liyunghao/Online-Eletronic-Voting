package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	db "github.com/liyunghao/Online-Eletronic-Voting/internal/server/database"
	srv "github.com/liyunghao/Online-Eletronic-Voting/internal/server/services"
	pb "github.com/liyunghao/Online-Eletronic-Voting/internal/voting"
	"google.golang.org/grpc"
)

var (
	port           = flag.Int("port", 8080, "Specify which port should gRPC server listen on")
	sqlite3db_name = flag.String("sqlite3db", "./database.db", "Specify which sqlite3 database should be used")
)

func main() {
	flag.Parse()

	// Initialize Database
	db.Initialize(*sqlite3db_name)

	tcp_listner, err := net.ListenTCP("tcp", &net.TCPAddr{IP: nil, Port: *port})
	if err != nil {
		log.Fatalf("Create TCP listner failed. Something WRONG: %v\n", err)
	}

	// EVoting Service
	var eVotingSrv srv.Service_eVoting

	grpcServer := grpc.NewServer()
	pb.RegisterEVotingServer(grpcServer, &eVotingSrv)

	// Start cli goroutine
	notifyStop := make(chan bool)
	go cli(notifyStop)

	// Start Server
	go func() {
		log.Printf("gRPC server start to listen at %d\n", *port)
		err = grpcServer.Serve(tcp_listner)
		if err != nil {
			log.Fatalf("Create TCP listner failed. Something WRONG: %v\n", err)
		}
	}()
	defer func() {
		grpcServer.GracefulStop()
		log.Printf("gRPC server stop\n")
	}()

	<-notifyStop
}

func cli(notifyStop chan bool) {
	// Readline from stdin
	stdin_scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("> ")

	for stdin_scanner.Scan() {
		// Scan command
		cmd := stdin_scanner.Text()

		// Execute command
		switch cmd {
		case "register":
			fmt.Printf("Enter name: ")
			stdin_scanner.Scan()
			name := stdin_scanner.Text()

			fmt.Printf("Enter group: ")
			stdin_scanner.Scan()
			group := stdin_scanner.Text()

			fmt.Printf("Enter public key: ")
			stdin_scanner.Scan()
			public_key := stdin_scanner.Text()

			// Register voter
			RegisterVoter(name, group, public_key)
		case "unregister":
			fmt.Printf("Enter name: ")
			stdin_scanner.Scan()
			name := stdin_scanner.Text()

			// Unregister voter
			UnregisterVoter(name)
		case "exit":
			notifyStop <- true
			return
		default:
			log.Printf("Unknown command: %s\n", cmd)
		}
		fmt.Printf("> ")
	}
}
