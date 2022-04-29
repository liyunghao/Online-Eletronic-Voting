package manager

import (
	"net/http"
)

var ClusterManager Manager


type Manager interface {
	// State and http server in implement session
	Initialize(args ...interface{}) // Creation, Control, open http server, join c
	Run()							// Http server start
	CheckRoles() string				// Check if leader

	// Active
	HeartBeat()
	WriteSync()
	ElectForLeader()				// backups call -> elect new leader
	CatchUp()						// Primary call -> 


	// Http route handler definition
	HeartBeatHandler(w http.ResponseWriter, r *http.Request)
	WriteSyncHandler(w http.ResponseWriter, r *http.Request)
	DeclareLeaderHandler(w http.ResponseWriter, r *http.Request)
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
