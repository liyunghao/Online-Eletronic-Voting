package manager

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

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

// vote & result
type vote struct {
	electionName string `json:"election_name"`
	voterName    string `json:"voter_name"`
	choice       string `json:"choice"`
}

type LfManager struct {
}

var LFManger LfManager

func (lf *LfManager) Initialize(args ...interface{}) error {
	http.HandleFunc("/hearbeat", lf.HeartBeatHandler)
	http.HandleFunc("/writesync", lf.WriteSyncHandler)
	http.HandleFunc("/declare_leader", lf.DeclareLeaderHandler)
	http.HandleFunc("/catch_up", lf.CatchUpHandler)
	http.HandleFunc("/recv_elect", lf.RecvElectHandler)
	log.Fatal(http.ListenAndServe(":9000", nil))
	return nil
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

	switch cmd {
	case "CreateUser":
		var payload st.User
		json.Unmarshal(payloadData, payload)
		st.DataStorage.CreateUser(payload.Name, payload.Group, payload.PublicKey)
	// case "FetchUser":
	//     var payload st.User
	//     json.Unmarshal(payloadData, payload)
	//     st.DataStorage.FetchUser(payload.Name)
	case "RemoveUser":
		var payload st.User
		json.Unmarshal(payloadData, payload)
		st.DataStorage.RemoveUser(payload.Name)

	case "CreateElection":
		var payload st.Election
		json.Unmarshal(payloadData, payload)
		st.DataStorage.CreateElection(payload.Name, payload.Groups, payload.Choices, payload.EndDate)
	// case "FetchElection":
	//     var payload st.Election
	//     json.Unmarshal(payloadData, payload)
	//     st.DataStorage.FetchElection(payload.Name)

	case "VoteElection":
		var payload vote
		json.Unmarshal(payloadData, payload)
		st.DataStorage.VoteElection(payload.electionName, payload.voterName, payload.electionName)

	}
}

func (lf *LfManager) DeclareLeaderHandler(w http.ResponseWriter, r *http.Request) {
}

func (lf *LfManager) CatchUpHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

}

func (lf *LfManager) RecvElectHandler(w http.ResponseWriter, r *http.Request) {

}
