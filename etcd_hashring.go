package etcdhashring

import (
	"context"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/serialx/hashring"
	"github.com/sirupsen/logrus"
	// "github.com/serialx/hashring"
)

//ETCDHashring main struct for lib ETCD Hashring
type ETCDHashring struct {
	etcdServicePath string
	etcdEndpoints   []string
	defaultTimeout  time.Duration
	CurrentNodeList []Node
	listeners       []Watcher
}

//Node info on ETCD
type Node struct {
	Name string
	Info map[string]string
}

//Watcher is the interface for nodes update listeners
type Watcher interface {
	NodesUpdated(nodes []Node)
}

//NewETCDHashring create a new ETCD Hashring instance. Call Start() to connect and start listening to ETCD changes
func NewETCDHashring(etcdEndpoints []string, etcdServicePath string, defaultTimeout time.Duration) (*ETCDHashring, error) {
	e := &ETCDHashring{}
	e.defaultTimeout = defaultTimeout
	e.etcdServicePath = etcdServicePath
	e.etcdEndpoints = etcdEndpoints
	e.listeners = make([]Watcher, 0)
	return e, nil
}

//AddWatcher add a new update listener
func (e *ETCDHashring) AddWatcher(listener Watcher) {
	e.listeners = append(e.listeners, listener)
}

//Connect to ETCD and start monitoring service node for changes
func (e *ETCDHashring) Connect() {
	cli, err := clientv3.New(clientv3.Config{Endpoints: e.etcdEndpoints, DialTimeout: 10 * time.Second})
	if err != nil {
		logrus.Errorf("Could not initialize ETCD client. err=%s", err)
		panic(err)
	}
	logrus.Debug("Etcd client initialized")

	nodesChan := make(chan []Node, 0)
	logrus.Debug("Starting to watch nodes at %s", e.etcdServicePath)
	go watchRegisteredNodes(cli, e.etcdServicePath, nodesChan)

	for {
		select {
		case nodes := <-nodesChan:
			logrus.Debugf("nodes list updated: %s", nodes)
			e.CurrentNodeList = nodes
			logrus.Debugf("notifying listeners")
			for _, l := range e.listeners {
				l.NodesUpdated(nodes)
			}
		}
	}
}

/*
GetNode get the target node according to the elementId.
For example, if the monitores Nodes in ETCD is a list of servers that accepts requests,
you can you can use GetNode("client-172.23.23.21") to identify to which server the client
would send its requests so that all servers will have a distributed set of clients targeting
requests to them, according to the idea of Consistent Hashing
*/
func (e *ETCDHashring) GetNode(elementID string) (Node, error) {
	ring := hashring.New(hashList(e.CurrentNodeList))
	nodeName, ok := ring.GetNode(stringSHA512(elementID))
	if !ok {
		return Node{}, fmt.Errorf("Could not get the node for this element id")
	}
	for _, n := range e.CurrentNodeList {
		if n.Name == nodeName {
			return n, nil
		}
	}
	return Node{}, fmt.Errorf("Elected node with name '%s' could not be found among current node list", nodeName)
}

func watchRegisteredNodes(cli *clientv3.Client, servicePath string, nodesChan chan []Node) {
	watchChan := cli.Watch(context.TODO(), servicePath, clientv3.WithPrefix())
	for {
		logrus.Debugf("Registered nodes updated")
		rsp, err0 := cli.Get(context.TODO(), servicePath, clientv3.WithPrefix())
		if err0 != nil {
			logrus.Warnf("Error retrieving service nodes. err=%s", err0)
		}

		if len(rsp.Kvs) == 0 {
			logrus.Debugf("No services nodes were found under %s", servicePath)

		} else {
			nodes := make([]Node, 0)
			for _, kv := range rsp.Kvs {
				node := Node{}
				node.Name = string(kv.Key)
				node.Info = decode(kv.Value)
				nodes = append(nodes, node)
			}
			nodesChan <- nodes
			logrus.Debugf("Found new registered nodes %s", nodes)
		}
		<-watchChan
	}
}

func decode(ds []byte) map[string]string {
	if ds != nil && len(ds) > 0 {
		var s map[string]string
		json.Unmarshal(ds, &s)
		return s
	}
	return nil
}

func stringSHA512(str string) string {
	hashedByte := sha512.Sum512([]byte(str))
	hashedString := string(hashedByte[:])

	return hashedString
}

func hashList(list []Node) []string {
	hashedList := make([]string, 0)
	for _, item := range list {
		hash := stringSHA512(item.Name)
		hashedList = append(hashedList, hash)
	}

	return hashedList
}
