package manager

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	Logs []Log `json: "logs"`
}

type Log struct {
	Cmd     string  `json: "storage_cmd"`
	Payload Payload `json: "payload"`
}

type Payload struct {
}

type LfManager struct {
}

func (m *LfManager) BroadcastHeartBeat() error {
	// read config file
	configFile, err := os.Open("config file name")
	if err != nil {
		// failed to read config file
		return status.Error(codes.Internal, "Internal server error: "+err.Error())
	}
	defer configFile.Close()
	byteValue, _ := ioutil.ReadAll(configFile)
	var config Config
	json.Unmarshal(byteValue, &config)

	// iterate through nodes to send heartbeat
	for i := 0; i < len(config.Clusters); i++ {
		if config.Clusters[i].Id != 0 { // suppose id 0 is leader
			postBody, _ := json.Marshal(map[string]int{
				"status": 200,
			})
			respBody := bytes.NewBuffer(postBody)
			resp, err := http.Post(config.Clusters[i].Ip, "application/json", respBody)
			if resp.StatusCode != 200 {
				return status.Error(codes.Internal, "Internal server error: "+err.Error())
			}

		}
	}

	return nil
}

func (m *LfManager) WriteSync(storageCmd string, payload string) error {
	// read config file
	configFile, err := os.Open("config file name")
	if err != nil {
		// failed to read config file
		return status.Error(codes.Internal, "Internal server error: "+err.Error())
	}
	defer configFile.Close()
	byteValue, _ := ioutil.ReadAll(configFile)
	var config Config
	json.Unmarshal(byteValue, &config)

	// iterate through nodes to send write sync
	for i := 0; i < len(config.Clusters); i++ {
		if config.Clusters[i].Id != 0 { // suppose id 0 is leader
			postBody, _ := json.Marshal(payload)
			respBody := bytes.NewBuffer(postBody)
			resp, err := http.Post(config.Clusters[i].Ip, "application/json", respBody)
			if resp.StatusCode != 200 {
				return status.Error(codes.Internal, "Internal server error: "+err.Error())
			}
		}
	}

	return nil
}

func (m *LfManager) ElectForLeader() error {
	// read config file
	configFile, err := os.Open("config file name")
	if err != nil {
		// failed to read config file
	}
	defer configFile.Close()
	byteValue, _ := ioutil.ReadAll(configFile)
	var config Config
	json.Unmarshal(byteValue, &config)

	// iterate through all nodes to check if this is the highest priority node
	var i int
	for i = 0; i < len(config.Clusters); i++ {
		if config.Clusters[i].Id < config.Node.Id {
			// not the highest priority
			break
		}
	}
	if i == len(config.Clusters) {
		// this is the highest priority node to become leader
		for i = 0; i < len(config.Clusters); i++ {
			if config.Clusters[i].Id != config.Node.Id {
				postBody, _ := json.Marshal(map[string]int{
					"node_idx": config.Node.Id,
				})
				respBody := bytes.NewBuffer(postBody)
				resp, err := http.Post(config.Clusters[i].Ip, "application/json", respBody)
				if resp.StatusCode != 200 {
					return status.Error(codes.Internal, "Internal server error: "+err.Error())
				}
			}
		}
	}
	return nil
}

func (m *LfManager) CatchUp() error {
	// read config file
	configFile, err := os.Open("config file name")
	if err != nil {
		// failed to read config file
	}
	defer configFile.Close()
	byteValue, _ := ioutil.ReadAll(configFile)
	var config Config
	json.Unmarshal(byteValue, &config)

	snapshot_id := 1 // need to know how snapshot id is stored

	idx := -1
	var ip string
	for i := 0; i < len(config.Clusters); i++ {
		if config.Clusters[i].Id < idx || idx < 0 {
			idx = config.Clusters[i].Id
			ip = config.Clusters[i].Ip
		}
	}
	postBody, _ := json.Marshal(map[string]int{
		"snapshot_id": snapshot_id,
	})
	resBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(ip, "application/json", resBody)
	if resp.StatusCode != 200 {
		return status.Error(codes.Internal, "Internal server error: "+err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var logs Logs
	json.Unmarshal(body, &logs)
	for i := 0; i < len(logs.Logs); i++ {
		// store logs.Logs[i].Payload into logs
		// need to know how the logs are stored
	}
	return nil
}
