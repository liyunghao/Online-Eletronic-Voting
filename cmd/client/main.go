package main

import (
	context "context"
	flag "flag"
	fmt "fmt"
	insecure "google.golang.org/grpc/credentials/insecure"
	grpc "google.golang.org/grpc"
	pb "github.com/liyunghao/Online-Eletronic-Voting/internal/voting"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	var election *pb.Election
	var ctx context.Context
	var conn *grpc.ClientConn
	var name *pb.VoterName
	var authrequest *pb.AuthRequest
	var vote *pb.Vote
	var election_name *pb.ElectionName


	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	fmt.Println(err)

	client := pb.NewEVotingClient(conn)
	var e error
	var chal *pb.Challenge
	var token *pb.AuthToken
	var s1 *pb.Status
	var s2 *pb.Status
	var result *pb.ElectionResult
	// five function for eVoting
	chal, e = client.PreAuth(ctx, name)
	fmt.Println(chal, e)
	token, e  = client.Auth(ctx, authrequest)
	fmt.Println(token, e)
	s1, e  = client.CreateElection(ctx, election)
	fmt.Println(s1, e)
	s2, e = client.CastVote(ctx, vote)
	fmt.Println(s2, e)
	result, e = client.GetResult(ctx, election_name)
	fmt.Println(result, e)
	fmt.Println("Hello Client")
}
