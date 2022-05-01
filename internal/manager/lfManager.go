package manager

import (
	"io/ioutil"
	"log"
	"net/http"
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

}

func (lf *LfManager) DeclareLeaderHandler(w http.ResponseWriter, r *http.Request) {
}

func (lf *LfManager) CatchUpHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

}

func (lf *LfManager) RecvElectHandler(w http.ResponseWriter, r *http.Request) {

}
