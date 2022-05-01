package manager

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	st "github.com/liyunghao/Online-Eletronic-Voting/internal/storage"
)

type Node struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type Cluster struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Id   int    `json:"id"`
}

type Payload struct {
	Name         string    `json:"name"`
	Group        string    `json:"group"`
	PublicKey    string    `json:"public_key"`
	Groups       []string  `json:"groups"`
	Choices      []string  `json:"choices"`
	EndDate      time.Time `json:"end_date"`
	ElectionName string    `json:"election_name"`
	VoterName    string    `json:"voter_name"`
	Choice       string    `json:"choice"`
}

type Log struct {
	Cmd     string  `json:"storage_cmd"`
	Payload Payload `json:"payload"`
}

var logs []Log

// vote & result
type vote struct {
	ElectionName string `json:"election_name"`
	VoterName    string `json:"voter_name"`
	Choice       string `json:"choice"`
}

type LfManager struct {
}

func (lf *LfManager) HeartBeatHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (lf *LfManager) WriteSyncHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	data := make(map[string]interface{})
	json.Unmarshal(body, &data)
	cmdData, _ := json.Marshal(data["storage_cmd"])
	payloadData, _ := json.Marshal(data["payload"])

	var cmd string
	json.Unmarshal(cmdData, &cmd)

	var logPayload Payload

	switch cmd {
	case "CreateUser":
		var user st.User
		json.Unmarshal(payloadData, &user)
		st.DataStorage.CreateUser(user.Name, user.Group, user.PublicKey)
		logPayload.Name = user.Name
		logPayload.Group = user.Group
		logPayload.PublicKey = user.Group
	// case "FetchUser":
	//     var payload st.User
	//     json.Unmarshal(payloadData, payload)
	//     st.DataStorage.FetchUser(payload.Name)
	case "RemoveUser":
		var user st.User
		json.Unmarshal(payloadData, &user)
		st.DataStorage.RemoveUser(user.Name)
		logPayload.Name = user.Name

	case "CreateElection":
		var election st.Election
		json.Unmarshal(payloadData, &election)
		st.DataStorage.CreateElection(election.Name, election.Groups, election.Choices, election.EndDate)
		logPayload.Name = election.Name
		logPayload.Groups = election.Groups
		logPayload.Choices = election.Choices
		logPayload.EndDate = election.EndDate

	// case "FetchElection":
	//     var payload st.Election
	//     json.Unmarshal(payloadData, payload)
	//     st.DataStorage.FetchElection(payload.Name)

	case "VoteElection":
		var election vote
		json.Unmarshal(payloadData, &election)
		st.DataStorage.VoteElection(election.ElectionName, election.VoterName, election.Choice)
		logPayload.ElectionName = election.ElectionName
		logPayload.VoterName = election.VoterName
		logPayload.Choice = election.Choice
	}
	w.WriteHeader(http.StatusOK)
	logs = append(logs, Log{cmd, logPayload})
}

func (lf *LfManager) DeclareLeaderHandler(w http.ResponseWriter, r *http.Request) {
}

func (lf *LfManager) CatchUpHandler(w http.ResponseWriter, r *http.Request) {
	// body, _ := ioutil.ReadAll(r.Body)
	w.Header().Set("Content-Type", "application/json")
	logsJson := make(map[string][]Log)
	logsJson["logs"] = logs
	data, _ := json.Marshal(logsJson)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (lf *LfManager) RecvElectHandler(w http.ResponseWriter, r *http.Request) {

}
