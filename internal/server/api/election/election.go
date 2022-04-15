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
	choices []string
    votes []int
	t time.Time
}

elections := make(map[string]Election)

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
    var status int32
    var counts []*pb.VoteCount
    if elecName, ok := elections[elecName]; ok {
        now := time.Now()
        if elections[elecName].t.Before(now) {
            status = 0
            for i := range {
                choiceName = elections[elecName].choices[i]
                count = elections[elecName].votes[i]
                counts = append(count, *pd.VoteCount{ChioceName: , Count: count}) 
            }
        } else {
            status = 2
        }
    } else {
        status = 1
    }

	return &pb.ElectionResult{
		Status: status,
        Count: counts
		// Counts: []*pb.VoteCount{ChoiceName: "Trump", Count: 1},
	}, nil
}
