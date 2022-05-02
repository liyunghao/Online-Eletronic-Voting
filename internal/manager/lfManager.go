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
	"strconv"
	"time"

	"github.com/gorilla/mux"
	st "github.com/liyunghao/Online-Eletronic-Voting/internal/storage"
)

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
	go lf.Start(args[1].(int))
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

	var cmd int
	json.Unmarshal(cmdData, &cmd)
	var payload string
	json.Unmarshal(payloadData, &payload)

	if err := st.DataStorage.(*st.ReplicaLogWrapper).SynctoStorage(cmd, payload, true); err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		// error status
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (lf *LfManager) CatchUpHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	data := make(map[string]int)
	json.Unmarshal(body, &data)
	logIdData, _ := json.Marshal(data["log_id"])
	var logId int
	json.Unmarshal(logIdData, &logId)
	if logs, err := st.DataStorage.(*st.ReplicaLogWrapper).CatchUp(logId); err == nil {
		w.Header().Set("Content-Type", "application/json")
		data, _ := json.Marshal(logs)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (lf *LfManager) Start(controlPort int) error {
	router := mux.NewRouter().StrictSlash(true)

	// Only handle POSTS request
	router.HandleFunc("/heartbeat", lf.HeartBeatHandler).Methods("POST")
	router.HandleFunc("/writesync", lf.WriteSyncHandler).Methods("POST")
	router.HandleFunc("/catch_up", lf.CatchUpHandler).Methods("POST")

	lf.server = &http.Server{Addr: ":" + strconv.Itoa(controlPort), Handler: router} // using self-defined router instead of DefaultServeMux
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
