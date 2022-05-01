package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	// Replica Log File Name
	replicaLogFileName = "replica.log"
)

type ReplicaLogWrapper struct {
	engine Storage

	// This is where log entries are stored
	logFile *os.File
}

// 1. Code here is actually not fully-implemented write-ahead-logging
// For the cleanness of the code, I decided to leave it as is QQ...
//
// 2. Logging here doesn't support rotation. But it does support start in clean
// mode, which will remove the log file instead of recovering from it.

// First arg should be the pointer to the storage object which being proxy
func (r *ReplicaLogWrapper) Initialize(args ...interface{}) error {
	var err error

	// Register the engine
	r.engine = args[0].(Storage)

	// Clean mode
	perm := os.O_APPEND | os.O_CREATE | os.O_RDWR
	if args[1].(bool) {
		perm |= os.O_TRUNC
	}
	r.logFile, err = os.OpenFile(replicaLogFileName, perm, 0644)
	if err != nil {
		return err
	}

	// Recovering
	err = r.recover()
	if err != nil {
		return err
	}
	log.Println("Log Recovering Done")

	return nil
}

func (r *ReplicaLogWrapper) recover() error {
	reader := bufio.NewScanner(r.logFile)

	for reader.Scan() {
		line := reader.Text()

		// Parse
		switch line[0] {
		case '1':
			var user User
			err := json.Unmarshal([]byte(line[2:]), &user)
			if err != nil {
				return err
			}
			_ = r.engine.CreateUser(user.Name, user.Group, user.PublicKey)
		case '2':
			var param struct {
				Name string `json:"name"`
			}
			err := json.Unmarshal([]byte(line[2:]), &param)
			if err != nil {
				return err
			}
			_ = r.engine.RemoveUser(param.Name)
		case '3':
			var election Election
			err := json.Unmarshal([]byte(line[2:]), &election)
			if err != nil {
				return err
			}
			_ = r.engine.CreateElection(election.Name, election.Groups, election.Choices, election.EndDate)
		case '4':
			var vote struct {
				ElectionName string `json:"election_name"`
				VoterName    string `json:"voter_name"`
				Choice       string `json:"choice"`
			}
			err := json.Unmarshal([]byte(line[2:]), &vote)
			if err != nil {
				return err
			}
			_ = r.engine.VoteElection(vote.ElectionName, vote.VoterName, vote.Choice)
		default:
			return fmt.Errorf("Invalid log entry: %s", line)
		}
	}

	return reader.Err()
}

func (r *ReplicaLogWrapper) log(t int, payload string) error {
	// Generate Logs
	entry := strconv.Itoa(t) + "|" + payload + "\n"

	_, err := r.logFile.Write([]byte(entry))

	return err
}

func (r *ReplicaLogWrapper) CreateUser(name string, group string, publicKey string) error {
	err := r.engine.CreateUser(name, group, publicKey)
	if err != nil {
		return err
	}

	// Prepare Payload
	p, err := json.Marshal(User{
		Name:      name,
		Group:     group,
		PublicKey: publicKey,
	})

	if err != nil {
		return err
	}

	return r.log(1, string(p))
}

func (r *ReplicaLogWrapper) FetchUser(name string) (User, error) {
	return r.engine.FetchUser(name)
}

func (r *ReplicaLogWrapper) RemoveUser(name string) error {
	err := r.engine.RemoveUser(name)
	if err != nil {
		return err
	}

	p, err := json.Marshal(struct {
		Name string `json:"name"`
	}{
		Name: name,
	})
	if err != nil {
		return err
	}

	return r.log(2, string(p))
}

func (r *ReplicaLogWrapper) CreateElection(name string, groups []string, choices []string, endDate time.Time) error {
	err := r.engine.CreateElection(name, groups, choices, endDate)
	if err != nil {
		return err
	}

	p, err := json.Marshal(Election{
		Name:    name,
		Groups:  groups,
		Choices: choices,
		EndDate: endDate,
	})
	if err != nil {
		return err
	}

	return r.log(3, string(p))
}

func (r *ReplicaLogWrapper) FetchElection(name string) (Election, error) {
	return r.engine.FetchElection(name)
}

func (r *ReplicaLogWrapper) VoteElection(electionName string, voterName string, choice string) error {
	err := r.engine.VoteElection(electionName, voterName, choice)
	if err != nil {
		return err
	}

	p, err := json.Marshal(struct {
		ElectionName string `json:"election_name"`
		VoterName    string `json:"voter_name"`
		Choice       string `json:"choice"`
	}{
		ElectionName: electionName,
		VoterName:    voterName,
		Choice:       choice,
	})
	if err != nil {
		return nil
	}

	return r.log(4, string(p))
}

func (r *ReplicaLogWrapper) FetchElectionResults(electionName string) (ElectionResults, error) {
	return r.engine.FetchElectionResults(electionName)
}
