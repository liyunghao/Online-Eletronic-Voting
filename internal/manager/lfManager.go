package manager

import (
	"os"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

type replicas struct {
	node Node
	cluster []Cluster
}

type LfManager struct {
	leader		bool					// if this node is currently leader
	primary		bool					// if this node is primary node 
	replicas							// record other replicas' info for broadcast?
}

func (l *LfManager) Initialize(args ...interface{}) error {
	// args[0] -> config filename 
	l.node, l.cluster = parseConfig(args[0].(string))
	// primary node's id is 1 as default
	if l.node.Id == 1 {
		l.primary = true
		l.leader = true
	}

	return nil
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

func (l *LfManager) Start(notifyStop chan bool) error {

	return nil
}

func (l *LfManager) GetRoles() bool {
	return l.primary
}
