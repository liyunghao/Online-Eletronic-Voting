package services

import (
	"context"

	mt "github.com/liyunghao/Online-Eletronic-Voting/internal/manager"
	"github.com/liyunghao/Online-Eletronic-Voting/internal/server/api/auth"
	api_elec "github.com/liyunghao/Online-Eletronic-Voting/internal/server/api/election"
	st "github.com/liyunghao/Online-Eletronic-Voting/internal/storage"
	pb "github.com/liyunghao/Online-Eletronic-Voting/internal/voting"
)

type Service_eVoting struct {
	pb.UnimplementedEVotingServer
}

// Authentication Related
func (s *Service_eVoting) PreAuth(ctx context.Context, name *pb.VoterName) (*pb.Challenge, error) {
	chal, err := auth.PreAuth(name)

	return chal, err
}

func (s *Service_eVoting) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthToken, error) {
	token, err := auth.Auth(req)

	return token, err
}

// Election
func (s *Service_eVoting) CreateElection(ctx context.Context, election *pb.Election) (*pb.Status, error) {
	if mt.ClusterManager.GetRoles() {
		status, err := api_elec.CreateElection(election)
		if err != nil || *status.Code != 0 {
			return status, err
		}
		latestLog := st.DataStorage.(*st.ReplicaLogWrapper).RetrieveLatestLog()
		_ = mt.ClusterManager.WriteSync(latestLog.T, latestLog.Value)
		return status, err
	} else {
		var code int32 = 100
		return &pb.Status{
			Code: &code,
		}, nil
	}
}

func (s *Service_eVoting) CastVote(ctx context.Context, vote *pb.Vote) (*pb.Status, error) {
	if mt.ClusterManager.GetRoles() {
		status, err := api_elec.CastVote(vote)
		if err != nil || *status.Code != 0 {
			return status, err
		}
		latestLog := st.DataStorage.(*st.ReplicaLogWrapper).RetrieveLatestLog()
		_ = mt.ClusterManager.WriteSync(latestLog.T, latestLog.Value)
		return status, err
	} else {
		var code int32 = 100
		return &pb.Status{
			Code: &code,
		}, nil
	}
}

func (s *Service_eVoting) GetResult(ctx context.Context, name *pb.ElectionName) (*pb.ElectionResult, error) {
	res, err := api_elec.GetResult(name)

	return res, err
}
