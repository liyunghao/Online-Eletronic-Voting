package manager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	st "github.com/liyunghao/Online-Eletronic-Voting/internal/storage"
)

type replicas struct {
	Node    Node      `json:"node"`
	Cluster []Cluster `json:"clusters"`
}

type LfManager struct {
	replicas               // record other replicas' info for broadcast?
	peerPort               int
	leader                 bool // if this node is currently leader
	primary                bool // if this node is primary node
	Server                 *http.Server
	CheckPrimaryAliveTimer *time.Timer
}

func (lf *LfManager) Initialize(args ...interface{}) error {

	lf.Node, lf.Cluster = lf.ParseConfig(args[0].(string)) // args[0] -> config filename
	// primary node's id is 1 as default
	lf.primary = false
	lf.leader = false
	lf.peerPort = lf.Cluster[0].ControlPort
	if lf.Node.Id == 1 {
		lf.primary = true
		lf.leader = true
		lf.peerPort = lf.Cluster[1].ControlPort
	}

	// start heartbeat
	ticker := time.NewTicker(5 * time.Second) // send heartbeat per 30sec
	lf.CheckPrimaryAliveTimer = time.AfterFunc(20*time.Second, func() {})
	go func() {
		for {
			select {
			case <-ticker.C:
				if lf.leader {
					lf.BroadcastHeartBeat()
				}
			}
		}
	}()

	// stop := make(chan os.Signal, 1)
	// signal.Notify(stop, os.Interrupt)

	// <-stop

	lf.CatchUp()

	return nil
}

func (lf *LfManager) BroadcastHeartBeat() error {
	// iterate through nodes to send heartbeat
	client := http.Client{
		Timeout: time.Second * 2,
	}

	for i := 0; i < len(lf.Cluster); i++ {
		if lf.Cluster[i].Id != lf.Node.Id { // suppose id 0 is leader
			resp, err := client.Post("http://"+lf.Cluster[i].Ip+":"+strconv.Itoa(lf.peerPort)+"/heartbeat", "application/json", strings.NewReader(""))
			if err != nil {
				return err
			}
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("Failed with status code %d: %v", resp.StatusCode, err)
			}

		}
	}

	return nil
}

func (lf *LfManager) WriteSync(storageCmd int, payload string) error {
	// iterate through nodes to send write sync
	for i := 0; i < len(lf.Cluster); i++ {
		if lf.Cluster[i].Id != lf.Node.Id { // suppose id 0 is leader
			postBody, _ := json.Marshal(st.WriteSyncLog{
				T:     storageCmd,
				Value: payload,
			})
			resp, err := http.Post("http://"+lf.Cluster[i].Ip+":"+strconv.Itoa(lf.peerPort)+"/writesync", "application/json", strings.NewReader(string(postBody)))
			if err != nil {
				return err
			}
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("Failed with status code %d: %v", resp.StatusCode, err)
			}
		}
	}
	return nil
}

func (lf *LfManager) CatchUp() error {
	log_id := st.DataStorage.(*st.ReplicaLogWrapper).GetNewestLogIndex() // need to know how snapshot id is stored

	var ip string
	if lf.Node.Id == 1 {
		ip = lf.Cluster[1].Ip
	} else {
		ip = lf.Cluster[0].Ip
	}
	// postBody, _ := json.Marshal(map[string]int{
	// 	"snapshot_id": snapshot_id,
	// })
	// resBody := bytes.NewBuffer(postBody)
	payload_string := "{\"log_id\": " + strconv.Itoa(log_id) + "}"
	postBody := strings.NewReader(payload_string)
	resp, err := http.Post("http://"+ip+":"+strconv.Itoa(lf.peerPort)+"/catch_up", "application/json", postBody)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed with status code %d: %v", resp.StatusCode, err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var logs []st.WriteSyncLog
	json.Unmarshal(body, &logs)

	for i := 0; i < len(logs); i++ {
		// store logs.Logs[i].log into replicaLogWrapper
		st.DataStorage.(*st.ReplicaLogWrapper).SynctoStorage(logs[i].T, logs[i].Value, true)
	}
	return nil
}

func (lf *LfManager) HeartBeatHandler(w http.ResponseWriter, r *http.Request) {
	lf.CheckPrimaryAliveTimer.Stop()
	lf.CheckPrimaryAliveTimer = time.AfterFunc(20*time.Second, func() {
		lf.leader = true
	})
	w.WriteHeader(http.StatusOK)
}

func (lf *LfManager) WriteSyncHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var log st.WriteSyncLog
	json.Unmarshal(body, &log)

	if err := st.DataStorage.(*st.ReplicaLogWrapper).SynctoStorage(log.T, log.Value, true); err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		// error status
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (lf *LfManager) CatchUpHandler(w http.ResponseWriter, r *http.Request) {
	if lf.primary == false {
		lf.leader = false
	}
	body, _ := ioutil.ReadAll(r.Body)
	var data struct {
		Log_id int `json:"log_id"`
	}
	_ = json.Unmarshal(body, &data)
	if logs, err := st.DataStorage.(*st.ReplicaLogWrapper).CatchUp(data.Log_id); err == nil {
		w.Header().Set("Content-Type", "application/json")
		data, _ := json.Marshal(logs)
		w.Write(data)
		// w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (lf *LfManager) Start() {
	router := mux.NewRouter().StrictSlash(true)

	// Only handle POSTS request
	router.HandleFunc("/heartbeat", lf.HeartBeatHandler).Methods("POST")
	router.HandleFunc("/writesync", lf.WriteSyncHandler).Methods("POST")
	router.HandleFunc("/catch_up", lf.CatchUpHandler).Methods("POST")

	lf.Server = &http.Server{Addr: ":" + strconv.Itoa(lf.Node.ControlPort), Handler: router} // using self-defined router instead of DefaultServeMux
	_ = lf.Server.ListenAndServe()
}

func (lf *LfManager) GetRoles() bool {
	return lf.leader
}

func (lf *LfManager) ParseConfig(filename string) (Node, []Cluster) {
	config, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer config.Close()
	var tmp replicas
	bytes, _ := ioutil.ReadAll(config)
	_ = json.Unmarshal(bytes, &tmp)

	return tmp.Node, tmp.Cluster
}
