package manager

import (
	"os"
	"os/signal"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"time"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"context"
)

type replicas struct {
	node Node
	cluster []Cluster
}

type LfManager struct {
	leader		bool					// if this node is currently leader
	primary		bool					// if this node is primary node 
	replicas							// record other replicas' info for broadcast?
	server		*http.Server
}

func (lf *LfManager) Initialize(args ...interface{}) error {

	lf.node, lf.cluster = parseConfig(args[0].(string)) // args[0] -> config filename 

	// primary node's id is 1 as default
	if lf.node.Id == 1 {
		lf.primary = true
		lf.leader = true
	}

	// run http server
	go lf.Start()
	// handshake?

	// start heartbeat
	ticker := time.NewTicker (30 * time.Second) // send heartbeat per 30sec 
	quit := make(chan struct{})					// backdoor to end this func by close(quit)
	go func() {
		for {
			select {
				case <- ticker.C:
					//lf.BroadcastHeartBeat()
				case <- quit:
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

func (lf *LfManager) Start() error {
	router := mux.NewRouter().StrictSlash(true)

	// Only handle POSTS request
	router.HandleFunc("/heartbeat", lf.HeartBeatHandler).Methods("POST")
	router.HandleFunc("/writesync", lf.WriteSyncHandler).Methods("POST")
	router.HandleFunc("/declare_capability", lf.DeclareLeaderHandler).Methods("POST")
	router.HandleFunc("/catch_up", lf.CatchUpHandler).Methods("POST")
	router.HandleFunc("/recv_elect", lf.RecvElectHandler).Methods("POST")

	lf.server = &http.Server{Addr: ":9000", Handler: router} // using self-defined router instead of DefaultServeMux
	if err := lf.server.ListenAndServe(); err != nil {
		log.Fatal("Http Server start error")
	}
	return nil
}

func (lf *LfManager) GetRoles() bool {
	return lf.primary
}

func parseConfig(filename string) (Node, []Cluster) {
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

func (lf *LfManager) HeartBeatHandler(w http.ResponseWriter, r *http.Request) {
	return
}

func (lf *LfManager) WriteSyncHandler(w http.ResponseWriter, r *http.Request) {
	return
}

func (lf *LfManager) DeclareLeaderHandler(w http.ResponseWriter, r *http.Request) {
	return
}

func (lf *LfManager) CatchUpHandler(w http.ResponseWriter, r *http.Request) {
	return
}

func (lf *LfManager) RecvElectHandler(w http.ResponseWriter, r *http.Request) {
	return
}


