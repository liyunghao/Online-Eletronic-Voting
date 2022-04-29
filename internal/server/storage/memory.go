package storage

import (
	"fmt"
	"time"
)

// Memory Storage Implementation
// This provide a simple in-memory storage implementation.
// It is not intended to be used in production, but rather as a
// simple way to test the server.
// It can also serve as a simple implementation of the storage interface
// , which allowed us to focus on fault tolerance implementation
// rather than worring about basic storage operation

type MemoryStorage struct {
	users     map[string]User
	elections map[string]struct {
		Election
		choices_res map[string]int

		// Performance Issue
		votedUsers map[string]struct{}
	}
}

func (m *MemoryStorage) Initialize(args ...interface{}) error {
	m.users = make(map[string]User)
	m.elections = make(map[string]struct {
		Election
		choices_res map[string]int
		votedUsers  map[string]struct{}
	})

	return nil
}

// Users
func (m *MemoryStorage) CreateUser(name string, group string, publicKey string) error {
	_, ok := m.users[name]
	if ok {
		return fmt.Errorf("user already exists")
	}
	m.users[name] = User{
		Name:      name,
		Group:     group,
		PublicKey: publicKey,
	}
	return nil
}

func (m *MemoryStorage) FetchUser(name string) (User, error) {
	user, ok := m.users[name]
	if !ok {
		return User{}, fmt.Errorf("user not found")
	}
	return user, nil
}

func (m *MemoryStorage) RemoveUser(name string) error {
	delete(m.users, name)
	return nil
}

// Election
func (m *MemoryStorage) CreateElection(name string, groups []string, choices []string, endDate time.Time) error {
	if _, ok := m.elections[name]; ok {
		return fmt.Errorf("election already exists")
	}

	m.elections[name] = struct {
		Election
		choices_res map[string]int
		votedUsers  map[string]struct{}
	}{
		Election: Election{
			Name:    name,
			Groups:  make([]string, len(groups)),
			Choices: make([]string, len(choices)),
			EndDate: endDate,
		},
		choices_res: make(map[string]int),
		votedUsers:  make(map[string]struct{}),
	}
	// Clean copy
	copy(m.elections[name].Groups, groups)
	copy(m.elections[name].Choices, choices)

	// Initialize the choices
	for _, choice := range choices {
		m.elections[name].choices_res[choice] = 0
	}

	return nil
}

func (m *MemoryStorage) FetchElection(name string) (Election, error) {
	election, ok := m.elections[name]
	if !ok {
		return Election{}, fmt.Errorf("election not found")
	}
	return election.Election, nil
}

// Votes & Results
func (m *MemoryStorage) VoteElection(electionName string, voterName string, choice string) error {
	election, ok := m.elections[electionName]
	if !ok {
		return fmt.Errorf("election not found")
	}
	if _, ok := election.choices_res[choice]; !ok {
		return fmt.Errorf("invalid choice")
	}
	_, ok = election.votedUsers[voterName]
	if ok {
		return fmt.Errorf("voter had already voted")
	}
	// Check if user has permission (groups) to vote
	hasPermission := false
	for _, group := range election.Groups {
		if group == m.users[voterName].Group {
			hasPermission = true
			break
		}
	}
	if !hasPermission {
		return fmt.Errorf("voter does not have permission to vote")
	}

	election.choices_res[choice]++
	election.votedUsers[voterName] = struct{}{}
	return nil
}

func (m *MemoryStorage) FetchElectionResults(electionName string) (ElectionResults, error) {
	election, ok := m.elections[electionName]
	if !ok {
		return ElectionResults{}, fmt.Errorf("election not found")
	}
	if election.EndDate.After(time.Now()) {
		return ElectionResults{}, fmt.Errorf("election is still ongoing")
	}
	results := make(ElectionResults)
	for choice, votes := range election.choices_res {
		results[choice] = int32(votes)
	}
	return results, nil
}
