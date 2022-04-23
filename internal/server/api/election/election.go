package election

import (
	"time"

	"github.com/liyunghao/Online-Eletronic-Voting/internal/server/jwt"
	pb "github.com/liyunghao/Online-Eletronic-Voting/internal/voting"
)

type election_rec struct {
	Name    string
	Groups  []string
	EndDate time.Time
	Choices map[string]int
	Voted   map[string]bool
}

var elections = make(map[string]election_rec)

func CreateElection(election *pb.Election) (*pb.Status, error) {
	_, err := jwt.VerifyToken(string(election.Token.Value))
	var ret_status int32 = 0

	if err != nil {
		ret_status = 1
		return &pb.Status{Code: &ret_status}, nil
	}
	if len(election.Groups) <= 0 || len(election.Choices) <= 0 {
		ret_status = 2
		return &pb.Status{Code: &ret_status}, nil
	}
	// Check if election exist
	if _, ok := elections[*election.Name]; ok {
		ret_status = 3
		return &pb.Status{Code: &ret_status}, nil
	}

	// Initialize Choices
	choices := make(map[string]int)
	for _, choice := range election.Choices {
		choices[choice] = 0
	}

	new_elect := election_rec{*election.Name, election.Groups, election.EndDate.AsTime(), choices, make(map[string]bool)}
	elections[*election.Name] = new_elect

	return &pb.Status{Code: &ret_status}, nil
}

func CastVote(vote *pb.Vote) (*pb.Status, error) {
	name, err := jwt.VerifyToken(string(vote.Token.Value))
	election, ok := elections[*vote.ElectionName]
	var ret_status int32 = 0

	if err != nil {
		// Invalid token
		ret_status = 1
		return &pb.Status{Code: &ret_status}, nil
	} else if !ok {
		// Invalid election name
		ret_status = 2
		return &pb.Status{Code: &ret_status}, nil
	} else if false {
		// check if user group
		ret_status = 3
		return &pb.Status{Code: &ret_status}, nil
	} else if _, ok := election.Voted[name]; ok {
		// already votes
		ret_status = 4
		return &pb.Status{Code: &ret_status}, nil
	} else {
		// Invalid choice
		if _, found := election.Choices[*vote.ChoiceName]; !found {
			ret_status = 5
			return &pb.Status{Code: &ret_status}, nil
		} else {
			elections[*vote.ElectionName].Choices[*vote.ChoiceName] += 1
			elections[*vote.ElectionName].Voted[name] = true
			return &pb.Status{Code: &ret_status}, nil
		}
	}
}

func GetResult(elecName *pb.ElectionName) (*pb.ElectionResult, error) {
	var status int32
	var counts []*pb.VoteCount
	if _, ok := elections[*elecName.Name]; ok {
		now := time.Now()
		if elections[*elecName.Name].EndDate.Before(now) {
			status = 0
			for choiceName, cnt := range elections[*elecName.Name].Choices {
				cnt := int32(cnt)
				counts = append(counts, &pb.VoteCount{ChoiceName: &choiceName, Count: &cnt})
			}
		} else {
			status = 2
		}
	} else {
		status = 1
	}

	return &pb.ElectionResult{
		Status: &status,
		Counts: counts,
	}, nil
}
