package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type MessageType string
type Status string
type Role string
type ErrorCode string

const (
	OK  ErrorCode = "ok"
	NOK ErrorCode = "nok"
)

const (
	ACTIF   Status = "actif"
	INACTIF Status = "inactif"
)

const (
	LEADER   Role = "leader"
	FOLLOWER Role = "follower"
)

const (
	SAVE   MessageType = "save"
	REQ    MessageType = "request"
	ACQ    MessageType = "acknowledge"
	NEW    MessageType = "new"
	REMOVE MessageType = "remove"
)

type Message struct {
	Mtype MessageType
	Body  string
}

type Listener struct {
	ReplicaId int
	Ltype     MessageType
	OnMessage func(r Replica, m Message, expected chan int)
}

type ReplicaService struct {
	Network Network
}

func (s *ReplicaService) Register(r *Replica) {
	fmt.Println("Register: ", r.Name)
	actifReplica := s.Network.filterBy(func(r Replica) bool { return (r.Status == ACTIF) })
	if len(actifReplica) == 0 {
		r.Role = LEADER
	} else {
		r.Role = FOLLOWER
	}
	r.Status = ACTIF
	//	r.Add("election", "/")
	r.Add("states", "/")
	s.Network.SendEvent(Message{Mtype: NEW, Body: strconv.Itoa(r.Id)}, s.Network)
}

func (s *ReplicaService) Broadcast(r *Replica) {
	fmt.Println("Broadcast: ", r.Name)
	actifReplica := s.Network.filterBy(func(r Replica) bool { return (r.Status == ACTIF) })
	if len(actifReplica) > 0 {
		for _, actif := range actifReplica {
			r.Add(strconv.Itoa(actif.Id), "/election")
		}
	}
}

func (s *ReplicaService) UnRegister(r *Replica) {
	fmt.Println("UnRegister: ", r.Name)
	r.Role = ""
	r.Status = INACTIF
	r.Root.Children = make(map[string]Node)
	s.Network.SendEvent(Message{Mtype: REMOVE, Body: strconv.Itoa(r.Id)}, s.Network)
}

func (s *ReplicaService) AddState(r *Replica, k string, value string) {
	fmt.Println("AddState: ", r.Name)
	result := s.Network.filterBy(func(r Replica) bool { return (r.Role == LEADER) && r.Status == ACTIF })
	if len(result) == 1 {
		//leader := result[0]
		data := map[string]string{"k": k, "v": "value", "status": "req"}
		conversion, _ := json.Marshal(data)
		synchronizedResults := s.Network.SendEvent(Message{Mtype: REQ, Body: string(conversion)}, s.Network)
		fmt.Println(synchronizedResults)
	}
}

type ReplicationScenario struct {
	Network Network
	Start   func(n *Replica, s ReplicaService)
	Execute func(n Network, s ReplicaService)
	Close   func(n Replica, s ReplicaService)
}

func (scenario ReplicationScenario) Do(waitTime int) {

	fmt.Println("start scenario -----------")
	service := ReplicaService{}
	service.Network = scenario.Network
	service.Network.Listeners = make(map[int][]Listener)
	defer func() {
		for _, r := range service.Network.Replica {
			scenario.Close(*r, service)
		}
	}()
	for i := 0; i < service.Network.Quorum; i++ {
		r := &Replica{Id: i, Name: createName(i), Root: createNode(i, "/"), Status: INACTIF}
		scenario.Network.Replica[i] = r
		service.Network.Listeners[i] = createListeners(*r)
	}
	for i := 0; i < service.Network.Quorum; i++ {
		scenario.Start(service.Network.Replica[i], service)
	}
	scenario.Execute(service.Network, service)
	time.Sleep(time.Duration(waitTime) * time.Second)
}

type Network struct {
	Quorum    int
	Replica   []*Replica
	SendEvent func(m Message, n Network) int
	Listeners map[int][]Listener
}

func (n *Network) filterBy(match func(r Replica) bool) []Replica {
	results := []Replica{}
	for _, v := range n.Replica {
		if match(*v) {
			results = append(results, *v)
		}
	}
	return results
}

func (n *Network) isActif(id int) bool {
	result := n.filterBy(func(r Replica) bool { return (r.Id == id) && r.Status == ACTIF })
	return len(result) == 1
}

type Replica struct {
	Id     int
	Status Status
	Role   Role
	Name   string
	Root   Node
}

func (r *Replica) Min() int {
	status, node := r.Root.Get("/election")
	if status == NOK {
		return 0
	}
	return Min(node[0])
}

func (r *Replica) Add(path string, rootPath string) {
	status, node := r.Root.Get(rootPath)
	if status == NOK {
		return
	}
	key := path
	if (!strings.HasPrefix("/", path)) && (rootPath != "/") {
		path = "/" + path
	}
	if r.Root.Path == key {
		return
	}
	if node[0].Children == nil {
		node[0].Children = make(map[string]Node)
	}
	if _, ok := node[0].Children[path]; ok == false {
		// case when path is a sub path of existing paths
		if status, res := contains(path, node[0].Children); status == "ok" && len(res) > 0 {
			res[0].Children[key] = createNode(getIdFromPath(key), res[0].Path+path)
		} else { // else add the path, it's a new one.
			node[0].Children[key] = createNode(getIdFromPath(key), node[0].Path+path)
		}
	}
}

func (r *Replica) Remove(path string, rootPath string) {
	status, node := r.Root.Get(rootPath)
	if status == NOK {
		return
	}
	key := path
	if !strings.HasPrefix("/", path) {
		path = "/" + path
	}
	if node[0].Path == key {
		return
	}
	if _, ok := node[0].Children[key]; ok == true {
		delete(node[0].Children, key)
	}
}

func (r *Replica) Get(path string) (ErrorCode, [1]Node) {
	if r.Root.Path == path || strings.Contains(r.Root.Path, path) {
		return OK, [1]Node{r.Root}
	}
	if r.Root.Children != nil {
		for _, v := range r.Root.Children {
			return v.Get(path)
		}
	}
	return NOK, [1]Node{}
}

type Node struct {
	Id       int
	Path     string
	Children map[string]Node
}

func Min(node Node) int {
	//fmt.Println("node.Path ", node.Path)
	min := getIdFromPath(node.Path)
	if len(node.Children) == 0 {
		return min
	}
	for _, child := range node.Children {
		//	fmt.Println("node.Path ", child.Path)
		if childMin := getIdFromPath(child.Path); min > childMin {
			min = childMin
			if childMin := Min(child); childMin < min {
				min = childMin
			}
		}
	}
	return min
}

func (r *Node) Get(path string) (ErrorCode, [1]Node) {
	if r.Path == path || strings.Contains(r.Path, path) {
		return OK, [1]Node{*r}
	}
	if r.Children != nil {
		for _, v := range r.Children {
			return v.Get(path)
		}
	}
	return NOK, [1]Node{}
}

func contains(path string, paths map[string]Node) (ErrorCode, [1]Node) {
	for k, v := range paths {
		if strings.Contains(k, path) {
			return OK, [1]Node{0: v}
		}
	}
	return NOK, [1]Node{}
}

func createReplica(id int) Replica {
	return Replica{
		Id:   id,
		Name: createName(id),
		Root: createNode(id, createPath(id)),
	}
}

func createName(id int) string {
	return "replica-" + strconv.Itoa(id)
}

func createPath(id int) string {
	return "/election/" + strconv.Itoa(id)
}

func getIdFromPath(path string) int {
	if (path == "") || (strings.HasPrefix(path, "/") || strings.Contains(path, "election") || strings.Contains(path, "states")) {
		return 0
	}
	val, err := strconv.Atoi(strings.TrimPrefix(path, "/election/"))
	return func() int {
		if err != nil {
			fmt.Printf("Error ", err)
			return 0
		} else {
			return val
		}
	}()
}

func createNode(id int, path string) Node {
	return Node{Id: id, Path: path, Children: make(map[string]Node)}
}

func createListeners(r Replica) []Listener {
	return []Listener{
		Listener{ReplicaId: r.Id, Ltype: NEW, OnMessage: func(r Replica, m Message, expected chan int) {
			r.Add(m.Body, "/election")
			expected <- r.Id
			//val[r.Id] = ACQ
		}},
		Listener{ReplicaId: r.Id, Ltype: REMOVE, OnMessage: func(r Replica, m Message, expected chan int) {
			r.Remove(m.Body, "/election")
			r.Remove(m.Body, "/states")
			expected <- r.Id
			// val := <-expected
			// val[r.Id] = ACQ
		}},
		Listener{ReplicaId: r.Id, Ltype: REQ, OnMessage: func(r Replica, m Message, expected chan int) {
			r.Add(m.Body, "/states")
			expected <- r.Id
			// val := <-expected
			// val[r.Id] = ACQ
		}},
	}
}

func OnMessage(replica Replica, listener Listener, m Message, expected chan int) {
	go func() { listener.OnMessage(replica, m, expected) }()
}

func createNetwork(quorum int) Network {
	return Network{
		Quorum:  quorum,
		Replica: make([]*Replica, quorum),
		SendEvent: func(m Message, n Network) int {
			expected := make(chan int, quorum)
			for _, v := range n.Listeners {
				for _, e := range v {
					result := n.filterBy(func(r Replica) bool { return r.Id == e.ReplicaId })
					if (len(result) == 1) && (n.isActif(result[0].Id)) && (e.Ltype == m.Mtype) {
						OnMessage(result[0], e, m, expected)
					}
				}
			}
			//			fmt.Println(<-expected)
			//			a := <-expected
			//	fmt.Println(a)
			return 0
		}}

}

func main() {
	network := createNetwork(3)

	// First scenario: start a simple broad cast and add state on replica
	ReplicationScenario{
		Network: network,
		Start: func(r *Replica, service ReplicaService) {
			fmt.Println("Start for ", r.Name)
			service.Broadcast(r)
			service.Register(r)
		},
		Execute: func(n Network, service ReplicaService) {
			//service.UnRegister(n.Replica[1])

			//service.Register(n.Replica[1])
		},
		Close: func(n Replica, service ReplicaService) {
			// stop
		},
	}.Do(10)

	for _, r := range network.Replica {
		fmt.Println("replica -> ", r)
	}

}
