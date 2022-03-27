package election

import (
	pb "github.com/liyunghao/Online-Eletronic-Voting/internal/voting"
)

func CreateElection(election *pb.Election) (*pb.Status, error) {
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
