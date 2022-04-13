package auth

import (
	"crypto/rand"
	"encoding/base64"
	"log"

	db "github.com/liyunghao/Online-Eletronic-Voting/internal/server/database"
	pb "github.com/liyunghao/Online-Eletronic-Voting/internal/voting"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*
#cgo LDFLAGS: -lsodium
#include <sodium.h>
*/
import "C"

var ChallengeCache = make(map[string][]byte)

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
	ChallengeCache[name.Name] = challenge

	return &pb.Challenge{Value: challenge}, nil
}

func Auth(req *pb.AuthRequest) (*pb.AuthToken, error) {
	// Check if requested voter has record in database.
	// If not, return error.
	// If yes, check if challenge is correct.
	// If yes, return auth token.
	res, err := db.SqliteDB.Query("SELECT public_key FROM voters WHERE name = ?", req.Name.Name)
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
	if !res.Next() {
		// Voter not found
		return nil, status.Error(codes.InvalidArgument, "Voter not found")
	} else {
		err = res.Scan(&sign_pk_base64)
		if err != nil {
			return nil, status.Error(codes.Internal, "Internal server error: "+err.Error())
		}
	}

	// Check if Voter name exist in Challenge Cache
	if _, ok := ChallengeCache[req.Name.Name]; !ok {
		return nil, status.Error(codes.InvalidArgument, "Voter's corresponding challenge not found")
	}

	// Verify Signature
	sign_pk, err := base64.StdEncoding.DecodeString(sign_pk_base64)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error: "+err.Error())
	}

	if C.crypto_sign_verify_detached(
		(*C.uchar)(&req.Response.Value[0]),
		(*C.uchar)(&ChallengeCache[req.Name.Name][0]),
		C.ulonglong(len(ChallengeCache[req.Name.Name])),
		(*C.uchar)(&sign_pk[0]),
	) == 0 {
		return &pb.AuthToken{Value: []byte("Verify")}, nil
	} else {
		return nil, status.Error(codes.PermissionDenied, "Invalid signature")
	}
}
