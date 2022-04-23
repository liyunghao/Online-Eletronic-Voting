package auth

import (
	"crypto/rand"
	"encoding/base64"
	"log"

	db "github.com/liyunghao/Online-Eletronic-Voting/internal/server/database"
	"github.com/liyunghao/Online-Eletronic-Voting/internal/server/jwt"
	pb "github.com/liyunghao/Online-Eletronic-Voting/internal/voting"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var challengeCache = make(map[string][]byte)

func PreAuth(name *pb.VoterName) (*pb.Challenge, error) {
	// Check if requested voter has record in database.
	// If not, return error.
	// If yes, return challenge.
	res, err := db.SqliteDB.Query("SELECT name FROM voters WHERE name = ?", name.Name)
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
	challengeCache[*name.Name] = challenge

	return &pb.Challenge{Value: challenge}, nil
}

func Auth(req *pb.AuthRequest) (*pb.AuthToken, error) {
	// Check if requested voter has record in database.
	// If not, return error.
	// If yes, check if challenge is correct.
	// If yes, return auth token.
	res, err := db.SqliteDB.Query("SELECT name, grouptype, public_key FROM voters WHERE name = ?", req.Name.Name)
	defer func() {
		err = res.Close()
		if err != nil {
			log.Fatalf("Close database failed. Something WRONG: %v\n", err)
		}
	}()
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error: "+err.Error())
	}

	var sign_pk_base64 string
	var name string
	var group string
	if !res.Next() {
		// Voter not found
		return nil, status.Error(codes.InvalidArgument, "Voter not found")
	} else {
		err = res.Scan(&name, &group, &sign_pk_base64)
		if err != nil {
			return nil, status.Error(codes.Internal, "Internal server error: "+err.Error())
		}
	}

	// Check if Voter name exist in Challenge Cache
	if _, ok := challengeCache[*req.Name.Name]; !ok {
		return nil, status.Error(codes.InvalidArgument, "Voter's corresponding challenge not found")
	}

	// Verify Signature
	sign_pk, err := base64.StdEncoding.DecodeString(sign_pk_base64)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error: "+err.Error())
	}

	if isValidSignature(sign_pk, challengeCache[*req.Name.Name], req.Response.Value) {
		// Generate Jwt token and send back to client
		token, err := jwt.GenerateToken(name, group)
		if err != nil {
			return nil, status.Error(codes.Internal, "Internal server error: "+err.Error())
		}

		return &pb.AuthToken{Value: token}, nil
	} else {
		return nil, status.Error(codes.PermissionDenied, "Invalid signature")
	}
}
