package manager

import (
	"encoding/json"
	"log"
	"net/http"
)

type LfManager struct {
}

var LFManger LfManager

func (lf *LfManager) Initialize(args ...interface{}) error {
	http.HandleFunc("/hearbeat", lf.HeartBeatHandler)
	http.HandleFunc("/writesync", lf.WriteSyncHandler)
	http.HandleFunc("/declare_capability", lf.DeclareLeaderHandler)
	http.HandleFunc("/catch_up", lf.CatchUpHandler)
	http.HandleFunc("/elect", lf.RecvElectHandler)
	log.Fatal(http.ListenAndServe(":9000", nil))
	return nil
}

func route(w http.ResponseWriter, r *http.Request) {
	// path := r.URL.Path

	switch {
	// case path == "/writesync":
	// WriteSyncHandler(w, r)

	}

}

func (lf *LfManager) HeartBeatHandler(w http.ResponseWriter, r *http.Request) {
	// r.ParseForm()
	storage_cmd := r.FormValue("storage_cmd")
	payload := r.FormValue("payload")
	json.Unmarshal()
}
func (lf *LfManager) WriteSyncHandler(w http.ResponseWriter, r *http.Request) {

}
func (lf *LfManager) DeclareLeaderHandler(w http.ResponseWriter, r *http.Request) {

}

func (lf *LfManager) CatchUpHandler(w http.ResponseWriter, r *http.Request) {

}
func (lf *LfManager) RecvElectHandler(w http.ResponseWriter, r *http.Request) {

}
