package auth

import (
	"crypto/rand"
	"log"

	db "github.com/liyunghao/Online-Eletronic-Voting/internal/server/database"
	pb "github.com/liyunghao/Online-Eletronic-Voting/internal/voting"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ChallengeCache = make(map[string][]byte)

func PreAuth(name *pb.VoterName) (*pb.Challenge, error) {
	// Check if requested voter has record in database.
	// If not, return error.
	// If yes, return challenge.
	res, err := db.SqliteDB.Query("SELECT * FROM voters WHERE name = ?", name.Name)
	defer func() {
		err = res.Close()
		if err != nil {
			log.Fatalf("Close database failed. Something WRONG: %v\n", err)
		}
	}()
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error: "+err.Error())
	}
	if !res.Next() {
		// Voter not found
		return nil, status.Error(codes.InvalidArgument, "Voter not found")
	}

	challenge := make([]byte, 32)
	_, err = rand.Read(challenge)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error: "+err.Error())
	}
	ChallengeCache[name.Name] = challenge

	return &pb.Challenge{Value: challenge}, nil
}

func Auth(name *pb.AuthRequest) (*pb.AuthToken, error) {
	return &pb.AuthToken{Value: []byte("maybeBase64Token")}, nil
}
