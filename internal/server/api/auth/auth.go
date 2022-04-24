package auth

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/liyunghao/Online-Eletronic-Voting/internal/server/jwt"
	st "github.com/liyunghao/Online-Eletronic-Voting/internal/server/storage"
	pb "github.com/liyunghao/Online-Eletronic-Voting/internal/voting"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var challengeCache = make(map[string][]byte)

func PreAuth(name *pb.VoterName) (*pb.Challenge, error) {
	// Check if requested voter has record in database.
	// If not, return error.
	// If yes, return challenge.
	_, err := st.DataStorage.FetchUser(*name.Name)
	if err != nil && err.Error() == "user not found" {
		return nil, status.Error(codes.InvalidArgument, "Voter not found")
	}

	challenge := make([]byte, 32)
	_, err = rand.Read(challenge)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error: "+err.Error())
	}
	challengeCache[*name.Name] = challenge

	return &pb.Challenge{Value: challenge}, nil
}

func Auth(req *pb.AuthRequest) (*pb.AuthToken, error) {
	// Check if requested voter has record in database.
	// If not, return error.
	// If yes, check if challenge is correct.
	// If yes, return auth token.
	user, err := st.DataStorage.FetchUser(*req.Name.Name)
	if err != nil && err.Error() == "user not found" {
		return nil, status.Error(codes.InvalidArgument, "Voter not found")
	}

	// Check if Voter name exist in Challenge Cache
	if _, ok := challengeCache[*req.Name.Name]; !ok {
		return nil, status.Error(codes.InvalidArgument, "Voter's corresponding challenge not found")
	}

	// Verify Signature
	sign_pk, err := base64.StdEncoding.DecodeString(user.PublicKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error: "+err.Error())
	}

	if isValidSignature(sign_pk, challengeCache[*req.Name.Name], req.Response.Value) {
		// Generate Jwt token and send back to client
		token, err := jwt.GenerateToken(user.Name, user.Group)
		if err != nil {
			return nil, status.Error(codes.Internal, "Internal server error: "+err.Error())
		}

		return &pb.AuthToken{Value: token}, nil
	} else {
		return nil, status.Error(codes.PermissionDenied, "Invalid signature")
	}
}
