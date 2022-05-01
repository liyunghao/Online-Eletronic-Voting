package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	st "github.com/liyunghao/Online-Eletronic-Voting/internal/storage"
)

type Config struct {
	Node     Node      `json:"node"`
	Clusters []Cluster `json:"clusters"`
}

type Node struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type Cluster struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Id   int    `json:"id"`
}

type Logs struct {
	Logs []LogString `json:"logs"`
}

type LogString struct {
	log string
}

type Log struct {
	Cmd     string `json:"storage_cmd"`
	Payload string `json:"payload"`
}

type Payload struct {
}

type LfManager struct {
	Config            Config               // config will be stored when Initiaize is called
	replicaLogWrapper st.ReplicaLogWrapper // this ReplicaLogWrapper may be defined in other place rather than in LfManager
}

func (m *LfManager) BroadcastHeartBeat() error {
	// iterate through nodes to send heartbeat
	for i := 0; i < len(m.Config.Clusters); i++ {
		if m.Config.Clusters[i].Id != 0 { // suppose id 0 is leader
			resp, err := http.PostForm("http://"+m.Config.Clusters[i].Ip, url.Values{})
			if resp.StatusCode != http.StatusOK {
				fmt.Println(err)
				return fmt.Errorf("Failed with status code %d", resp.StatusCode)
			}

		}
	}

	return nil
}

func (m *LfManager) WriteSync(storageCmd string, payload string) error {
	// iterate through nodes to send write sync
	for i := 0; i < len(m.Config.Clusters); i++ {
		if m.Config.Clusters[i].Id != 0 { // suppose id 0 is leader
			postBody, _ := json.Marshal(Log{
				Cmd:     storageCmd,
				Payload: payload,
			})
			respBody := bytes.NewBuffer(postBody)
			resp, err := http.Post("http://"+m.Config.Clusters[i].Ip, "application/json", respBody)
			if resp.StatusCode != http.StatusOK {
				fmt.Println(err)
				return fmt.Errorf("Failed with status code %d", resp.StatusCode)
			}
		}
	}

	return nil
}

func (m *LfManager) ElectForLeader() error {
	// iterate through all nodes to check if this is the highest priority node
	var i int
	for i = 0; i < len(m.Config.Clusters); i++ {
		if m.Config.Clusters[i].Id < m.Config.Node.Id {
			// not the highest priority
			break
		}
	}
	if i == len(m.Config.Clusters) {
		// this is the highest priority node to become leader
		for i = 0; i < len(m.Config.Clusters); i++ {
			if m.Config.Clusters[i].Id != m.Config.Node.Id {
				postBody, _ := json.Marshal(map[string]int{
					"node_idx": m.Config.Node.Id,
				})
				respBody := bytes.NewBuffer(postBody)
				resp, err := http.Post("http://"+m.Config.Clusters[i].Ip, "application/json", respBody)
				if resp.StatusCode != http.StatusOK {
					fmt.Println(err)
					return fmt.Errorf("Failed with status code %d", resp.StatusCode)
				}
			}
		}
	}
	return nil
}

func (m *LfManager) CatchUp() error {
	snapshot_id := 1 // need to know how snapshot id is stored

	idx := -1
	var ip string
	for i := 0; i < len(m.Config.Clusters); i++ {
		if m.Config.Clusters[i].Id < idx || idx < 0 {
			idx = m.Config.Clusters[i].Id
			ip = m.Config.Clusters[i].Ip
		}
	}
	postBody, _ := json.Marshal(map[string]int{
		"snapshot_id": snapshot_id,
	})
	resBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://"+ip, "application/json", resBody)
	if resp.StatusCode != http.StatusOK {
		fmt.Println(err)
		return fmt.Errorf("Failed with status code %d:", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var logs Logs
	json.Unmarshal(body, &logs)

	for i := 0; i < len(logs.Logs); i++ {
		// store logs.Logs[i].log into replicaLogWrapper
		m.replicaLogWrapper.logFile.Write([]byte(logs.Logs[i].log))
	}
	return nil
}
