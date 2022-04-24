package storage

import "time"

// Global Storage Variable
var DataStorage Storage

// Public Communication Interface
type User struct {
	Name      string
	Group     string
	PublicKey string
}

type Election struct {
	Name    string
	Groups  []string
	Choices []string
	EndDate time.Time
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
