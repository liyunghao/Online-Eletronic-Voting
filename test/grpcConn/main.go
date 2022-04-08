package main

import (
	"context"
	"flag"
	"log"
	"strconv"

	pb "github.com/liyunghao/Online-Eletronic-Voting/internal/voting"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port = flag.Int("port", 8080, "Localhost server's port")
)

func main() {
	log.Println("Start to test if gRPC server works correctly against each routes...")

	// Dial Options
	var opt []grpc.DialOption
	opt = append(opt, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// Connect to Server
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(*port), opt...)
	if err != nil {
		log.Fatalf("SOMETHING WRONG, test failed: %v\n", err)
	}
	defer conn.Close()

	// Get service
	votingSrv := pb.NewEVotingClient(conn)

	// Test against each route.
	_, err = votingSrv.PreAuth(context.Background(), &pb.VoterName{})
	if err != nil {
		log.Fatalf("Test failed at PreAuth: %v", err)
	}

	_, err = votingSrv.Auth(context.Background(), &pb.AuthRequest{})
	if err != nil {
		log.Fatalf("Test failed at Auth: %v", err)
	}

	_, err = votingSrv.CreateElection(context.Background(), &pb.Election{})
	if err != nil {
		log.Fatalf("Test failed at CreateElection: %v", err)
	}

	_, err = votingSrv.CastVote(context.Background(), &pb.Vote{})
	if err != nil {
		log.Fatalf("Test failed at CastVote: %v", err)
	}

	_, err = votingSrv.GetResult(context.Background(), &pb.ElectionName{})
	if err != nil {
		log.Fatalf("Test failed at GetResult: %v", err)
	}

	log.Println("All test pass, WE ARE GOOD TO GO!!!")
}
