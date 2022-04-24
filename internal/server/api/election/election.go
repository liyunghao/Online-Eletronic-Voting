package election

import (
	"github.com/liyunghao/Online-Eletronic-Voting/internal/server/jwt"
	st "github.com/liyunghao/Online-Eletronic-Voting/internal/server/storage"
	pb "github.com/liyunghao/Online-Eletronic-Voting/internal/voting"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
	// Create Election
	err = st.DataStorage.CreateElection(*election.Name, election.Groups, election.Choices, election.EndDate.AsTime())
	if err != nil {
		if err.Error() == "election already exists" {
			ret_status = 3
		} else {
			return nil, status.Error(codes.Internal, "Internal server error: "+err.Error())
		}
	}
	return &pb.Status{Code: &ret_status}, nil
}

func CastVote(vote *pb.Vote) (*pb.Status, error) {
	var ret_status int32 = 0
	name, err := jwt.VerifyToken(string(vote.Token.Value))
	if err != nil {
		// Invalid token
		ret_status = 1
		return &pb.Status{Code: &ret_status}, nil
	}

	// Check if election exist
	err = st.DataStorage.VoteElection(*vote.ElectionName, name, *vote.ChoiceName)

	if err != nil {
		if err.Error() == "election not found" {
			ret_status = 2
		} else if err.Error() == "voter does not have permission to vote" {
			ret_status = 3
		} else if err.Error() == "voter had already voted" {
			ret_status = 4
		} else if err.Error() == "invalid choice" {
			ret_status = 5
		}
	}

	return &pb.Status{Code: &ret_status}, nil
}

func GetResult(elecName *pb.ElectionName) (*pb.ElectionResult, error) {
	var ret_status int32 = 0
	res, err := st.DataStorage.FetchElectionResults(*elecName.Name)
	if err != nil {
		if err.Error() == "election not found" {
			ret_status = 1
		} else if err.Error() == "election is still ongoing" {
			ret_status = 2
		}
		return &pb.ElectionResult{
			Status: &ret_status,
			Counts: []*pb.VoteCount{},
		}, nil
	}
	res_count := make([]*pb.VoteCount, 0)
	for k, v := range res {
		k := k
		v := v
		res_count = append(res_count, &pb.VoteCount{
			ChoiceName: &k,
			Count:      &v,
		})
	}

	return &pb.ElectionResult{
		Status: &ret_status,
		Counts: res_count,
	}, nil
}
