package election

import (
	pb "github.com/liyunghao/Online-Eletronic-Voting/internal/voting"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"time"
)

type Election struct {
	name string
	groups []string
	choice []string
	t time
}

var elections []Election

func CreateElection(election *pb.Election) (*pb.Status, error) {
	if len(election.Groups) <= 0 || len(election.Choices) <= 0 {
		return &pb.Status{Code: 2}, status.Error(codes.InvalidArgument, "At least one group and one choice should be listed.")
	}
	new_elect := Election{}
	return &pb.Status{Code: 200}, nil
}

func CastVote(vote *pb.Vote) (*pb.Status, error) {
	return &pb.Status{Code: 200}, nil
}

func GetResult(elecName *pb.ElectionName) (*pb.ElectionResult, error) {
	return &pb.ElectionResult{
		Status: 200,
		Counts: []*pb.VoteCount{{ChoiceName: "Trump", Count: 1}},
	}, nil
}
