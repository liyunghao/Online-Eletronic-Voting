package auth

import (
	pb "github.com/liyunghao/Online-Eletronic-Voting/internal/voting"
)

func PreAuth(name *pb.VoterName) (*pb.Challenge, error) {
	return &pb.Challenge{Value: []byte("Approve")}, nil
}

func Auth(name *pb.AuthRequest) (*pb.AuthToken, error) {
	return &pb.AuthToken{Value: []byte("maybeBase64Token")}, nil
}
