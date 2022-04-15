package election

import (
	pb "github.com/liyunghao/Online-Eletronic-Voting/internal/voting"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"time"
	"github.com/liyunghao/Online-Eletronic-Voting/internal/server/jwt"
)

type Election struct {
	name string
	groups []string
	t time.Time
	choices map([string]int)
	voted map([string]bool)
}

var elections = make(map[string]Election)

func CreateElection(election *pb.Election) (*pb.Status, error) {
	if len(election.Groups) <= 0 || len(election.Choices) <= 0 {
		return &pb.Status{Code: 2}, status.Error(codes.InvalidArgument, "At least one group and one choice should be listed.")
	}
	new_elect := Election{election.Name, election.Groups, election.Choices, election.EndDate.AsTime()}
	value, isExist := elections[election.Name]
	if !isExist {
		return &pb.Status{Code: 3}, status.Error(codes.InvalidArgument, "Election already exists.")
	}
	elections[election.Name] = new_elect
	return &pb.Status{Code: 0}, nil
}

func CastVote(vote *pb.Vote) (*pb.Status, error) {
	tokenstring := vote.Token
	election, ok := elections[vote.ElectionName]
	name, err := VerifyToken(tokenstring)

	if err != nil {
		// Invalid token 
		return &pb.Status{Code: 1}, nil
	} else if !ok {
		// Invalid election name
		return &pb.Status{Code: 2}, nil
	} else if false {
		// check if user group
		return &pb.Status{Code: 3}, nil
	} else if _, ok := election.voted[name]; ok {
		// already votes
		return &pb.Status{Code: 4}, nil
	} else {
		// Invalid choice
		if _, found := election.choices[vote.ChoiceName]; !found {
			return &pb.Status{Code: 5}, nil
		} else {
			elections[vote.ElectionName].choices[vote.ChoiceName] += 1
			elections[vote.ElectionName].voted[name] = true
			return &pb.Status{Code: 0}, nil
		}
	}
}

func GetResult(elecName *pb.ElectionName) (*pb.ElectionResult, error) {
    var status int32
    var counts []*pb.VoteCount
    if elecName, ok := elections[elecName]; ok {
        now := time.Now()
        if elections[elecName].t.Before(now) {
            status = 0
			for choiceName, count := choices {
                counts = append(counts, *pd.VoteCount{ChoiceName: choiceName, Count: count})
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
	}, nil
}
