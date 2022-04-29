package manager

import (
	"net/http"
)

var ClusterManager Manager

type Manager interface {
	// State and http server in implement session
	// Creation, Control, setup http server.
	// May need to handle process of joining the cluster.
	// Parse Config
	Initialize(args ...interface{}) error

	// will be invoke in an new go routine
	Start(notifyStop chan bool) error

	// Getter
	GetRoles() string // Retrieve current roles

	// Active
	// Leaders Capabilities
	BroadcastHeartBeat() error
	WriteSync(storageCmd string, payload string) error

	// Followers Capabilities
	ElectForLeader() // backups call -> elect new leader
	CatchUp()        // Primary call ->

	// Http route handler definition
	// Follower's Capabilities
	HeartBeatHandler(w http.ResponseWriter, r *http.Request)
	WriteSyncHandler(w http.ResponseWriter, r *http.Request)
	DeclareLeaderHandler(w http.ResponseWriter, r *http.Request)

	// Leader's Capabilities
	CatchUpHandler(w http.ResponseWriter, r *http.Request)
	RecvElectHandler(w http.ResponseWriter, r *http.Request)
}

// TODO:
// 1. STATE + INIT							--> TONY
// 2. REPLICA LOG  							--> ELVEN
// 3. ACTIVE								--> TSWANG
// 4. HTTP SERVER + HTTP HANDLER			--> CPC
// 5. COMBINE INTO CURRENT SERVER CODE

//type Replica struct {
//Name	[]string `json:"name"`
//Ip		[]string `json:"ip"`
//Id		[]string `json:"id"`
//Leader	[]bool	 `json:"leader"`
//}

//var Replicas []Replica
