package manager

import (
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

func (lf *LfManager) HeartBeatHandler(w http.ResponseWriter, r *http.Request) {
	// r.ParseForm()
	// storage_cmd := r.FormValue("storage_cmd")
	// body, _ := ioutil.ReadAll(r.Body)
	// var storage_cmd string
	// json.Unmarshal(r.FormValue("storage_cmd"), storage_cmd)
}

func (lf *LfManager) WriteSyncHandler(w http.ResponseWriter, r *http.Request) {

}

func (lf *LfManager) DeclareLeaderHandler(w http.ResponseWriter, r *http.Request) {

}

func (lf *LfManager) CatchUpHandler(w http.ResponseWriter, r *http.Request) {

}

func (lf *LfManager) RecvElectHandler(w http.ResponseWriter, r *http.Request) {

}
