package storage

import "time"

// Global Storage Variable
var DataStorage Storage

// Public Write API Definition
const (
	WriteAPI_CreateUser     = iota
	WriteAPI_RemoveUser     = iota
	WriteAPI_CreateElection = iota
	WriteAPI_VoteElection   = iota
)

// CommPayload, transfered between nodes in the control network or used as
// hardening logs that are stored in the log file.
type CommRemoveUserPayload struct {
	Name string `json:"name"`
}

type CommVotePayload struct {
	ElectionName string `json:"election_name"`
	VoterName    string `json:"voter_name"`
	Choice       string `json:"choice"`
}

type CommUserPayload User
type CommElectionPayload Election

// Code Component Interface
type User struct {
	Name      string `json:"name"`
	Group     string `json:"group"`
	PublicKey string `json:"public_key"`
}

type Election struct {
	Name    string    `json:"name"`
	Groups  []string  `json:"groups"`
	Choices []string  `json:"choices"`
	EndDate time.Time `json:"end_date"`
}

type ElectionResults map[string]int32

// Storage Interface Definition
type Storage interface {
	// Initialize the storage
	Initialize(args ...interface{}) error

	// Users
	CreateUser(name string, group string, publicKey string) error
	FetchUser(name string) (User, error)
	RemoveUser(name string) error

	// Elections
	CreateElection(name string, groups []string, choices []string, endDate time.Time) error
	FetchElection(name string) (Election, error)

	// Votes & Results
	VoteElection(electionName string, voterName string, choice string) error
	FetchElectionResults(electionName string) (ElectionResults, error)
}
