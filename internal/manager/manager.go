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

	ParseConfig(string) (Node, []Cluster)
	// will be invoke in an new go routine
	Start(controlPort int)

	// Getter
	GetRoles() bool // Retrieve current roles; return 1 if Primary else 0

	// Active
	// Leaders Capabilities
	BroadcastHeartBeat() error
	WriteSync(storageCmd int, payload string) error

	// Followers Capabilities
	CatchUp() error // Primary call ->

	// Http route handler definition
	// Follower's Capabilities
	HeartBeatHandler(w http.ResponseWriter, r *http.Request)
	WriteSyncHandler(w http.ResponseWriter, r *http.Request)

	// Leader's Capabilities
	CatchUpHandler(w http.ResponseWriter, r *http.Request)
}

// TODO:
// 1. STATE + INIT							--> TONY
// 2. REPLICA LOG							--> ELVEN
// 3. ACTIVE								--> TSWANG
// 4. HTTP SERVER + HTTP HANDLER			--> CPC
// 5. COMBINE INTO CURRENT SERVER CODE

type Node struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type Cluster struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Id   int    `json:"id"`
}
