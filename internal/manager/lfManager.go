package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	st "github.com/liyunghao/Online-Eletronic-Voting/internal/storage"
)

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

type replicas struct {
	node    Node
	cluster []Cluster
}

type LfManager struct {
	replicas      // record other replicas' info for broadcast?
	leader   bool // if this node is currently leader
	primary  bool // if this node is primary node
	server   *http.Server
}

func (lf *LfManager) Initialize(args ...interface{}) error {

	lf.node, lf.cluster = lf.ParseConfig(args[0].(string)) // args[0] -> config filename

	// primary node's id is 1 as default
	if lf.node.Id == 1 {
		lf.primary = true
		lf.leader = true
	}

	// run http server
	go lf.Start()
	// handshake?

	// start heartbeat
	ticker := time.NewTicker(30 * time.Second) // send heartbeat per 30sec
	quit := make(chan struct{})                // backdoor to end this func by close(quit)
	go func() {
		for {
			select {
			case <-ticker.C:
				//lf.BroadcastHeartBeat()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, shutdown := context.WithTimeout(context.Background(), 3*time.Second)
	defer shutdown()
	if err := lf.server.Shutdown(ctx); err != nil {
		log.Fatal("Http Server shutdown error")
	}

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
func (lf *LfManager) Start() error {
	router := mux.NewRouter().StrictSlash(true)

	// Only handle POSTS request
	router.HandleFunc("/heartbeat", lf.HeartBeatHandler).Methods("POST")
	router.HandleFunc("/writesync", lf.WriteSyncHandler).Methods("POST")
	router.HandleFunc("/catch_up", lf.CatchUpHandler).Methods("POST")

	lf.server = &http.Server{Addr: ":9000", Handler: router} // using self-defined router instead of DefaultServeMux
	if err := lf.server.ListenAndServe(); err != nil {
		log.Fatal("Http Server start error")
	}
	return nil
}

func (lf *LfManager) GetRoles() bool {
	return lf.primary
}

func (lf *LfManager) ParseConfig(filename string) (Node, []Cluster) {
	config, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer config.Close()
	var tmp replicas
	bytes, _ := ioutil.ReadAll(config)
	json.Unmarshal(bytes, &tmp)

	return tmp.node, tmp.cluster

}
