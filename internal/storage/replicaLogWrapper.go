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

// CommPayload
type removeUserPayload struct {
	Name string `json:"name"`
}

type votePayload struct {
	ElectionName string `json:"election_name"`
	VoterName    string `json:"voter_name"`
	Choice       string `json:"choice"`
}

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

// Control Interface
func (r *ReplicaLogWrapper) SynctoStorage(cmd int, payload string) error {
	switch cmd {
	case WriteAPI_CreateUser:
		var user User
		err := json.Unmarshal([]byte(payload), &user)
		if err != nil {
			return err
		}
		_ = r.engine.CreateUser(user.Name, user.Group, user.PublicKey)
	case WriteAPI_RemoveUser:
		var param removeUserPayload
		err := json.Unmarshal([]byte(payload), &param)
		if err != nil {
			return err
		}
		_ = r.engine.RemoveUser(param.Name)
	case WriteAPI_CreateElection:
		var election Election
		err := json.Unmarshal([]byte(payload), &election)
		if err != nil {
			return err
		}
		_ = r.engine.CreateElection(election.Name, election.Groups, election.Choices, election.EndDate)
	case WriteAPI_VoteElection:
		var vote votePayload
		err := json.Unmarshal([]byte(payload), &vote)
		if err != nil {
			return err
		}
		_ = r.engine.VoteElection(vote.ElectionName, vote.VoterName, vote.Choice)
	default:
		return fmt.Errorf("Invalid entry: %d -> %s", cmd, payload)
	}
	return nil
}

// Storage Interface Implementation
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
		logType, err := strconv.Atoi(line[:1])
		if err != nil {
			return err
		}
		err = r.SynctoStorage(logType, line[2:])
		if err != nil {
			return err
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

	p, err := json.Marshal(removeUserPayload{
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

	p, err := json.Marshal(votePayload{
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
