package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	jwt "github.com/liyunghao/Online-Eletronic-Voting/internal/server/jwt"
	srv "github.com/liyunghao/Online-Eletronic-Voting/internal/server/services"
	st "github.com/liyunghao/Online-Eletronic-Voting/internal/storage"
	pb "github.com/liyunghao/Online-Eletronic-Voting/internal/voting"
	"google.golang.org/grpc"
)

var (
	port           = flag.Int("port", 8080, "Specify which port should gRPC server listen on")
	storage_type   = flag.String("storage", "memory", "Specify which storage type should be used")
	sqlite3db_name = flag.String("sqlite3db", "./database.db", "Specify which sqlite3 database should be used")
)

func main() {
	flag.Parse()

	// Initialize Storage System (Currently only support memory storage)
	st.DataStorage = &st.MemoryStorage{}
	st.DataStorage.Initialize()

	// Initialize JWT
	jwt.InitJWT()

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
			err := st.DataStorage.CreateUser(name, group, public_key)
			if err != nil {
				fmt.Printf("Register failed. Something WRONG: %v\n", err)
			} else {
				fmt.Printf("Register success\n")
			}
		case "unregister":
			fmt.Printf("Enter name: ")
			stdin_scanner.Scan()
			name := stdin_scanner.Text()

			// Unregister voter
			err := st.DataStorage.RemoveUser(name)
			if err != nil {
				fmt.Printf("Unregister failed. Something WRONG: %v\n", err)
			} else {
				fmt.Printf("Unregister success\n")
			}
		case "exit":
			notifyStop <- true
			return
		default:
			log.Printf("Unknown command: %s\n", cmd)
		}
		fmt.Printf("> ")
	}
}
