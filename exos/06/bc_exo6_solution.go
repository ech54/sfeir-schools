package main

import (
	"fmt"
	"strconv"
	"strings"
)

type MessageType string
type ErrorCode string

const (
	OK  ErrorCode = "ok"
	NOK ErrorCode = "nok"
)

const (
	SAVE   MessageType = "save"
	ACQ    MessageType = "acknowledge"
	REQ    MessageType = "request"
	NEW    MessageType = "new"
	REMOVE MessageType = "remove"
)

type Message struct {
	Mtype MessageType
	Body  string
}

type Listener struct {
	Destination Replica
	Ltype       MessageType
	OnMessage   func(r Replica, m chan Message)
}

type ReplicaService struct {
	Network Network

	Start   func(n Network)
	Execute func(n Network)
	Close   func(n Network)
}

func (service ReplicaService) Do() {
	service.Network.Listeners = make(map[string][]Listener)
	defer service.Close(service.Network)
	service.Start(service.Network)
	service.Execute(service.Network)
}

type Network struct {
	Quorum    int
	Replica   []Replica
	SendEvent func(m Message, listeners map[string][]Listener)
	Listeners map[string][]Listener
}

type Replica struct {
	Id   int
	Name string
	Root Node
}

type Node struct {
	Id int
	Path     string
	Children map[string]Node
}

func (r *Replica) Min() int {
	status, node := r.Root.Get("/election")
	if status == NOK {
		return 0
	}

}

func (r *Replica) Min(n Node) (int, ErrorCode) {

	status, node := r.Root.Get("/election")
	if status == NOK {
		return 0, NOK
	}
	if len(node.Children) == 0 {
		return node.
	}
	return 0, NOK
}

func (r *Replica) Count() int {
	status, node := r.Root.Get("/election")
	if status == NOK {
		return 0
	}
	count := Count(node[0])
	return count
}

func Count(node Node) int {
	count := 0
	for _, child := range node.Children {
		count += 1
		if len(child.Children) > 0 {
			count += Count(child)
		}
	}
	return count
}

func (r *Replica) Add(path string) {
	key := path
	if !strings.HasPrefix("/", path) {
		path = "/" + path
	}
	fmt.Println(path)
	if r.Root.Path == key {
		return
	}
	if r.Root.Children == nil {
		r.Root.Children = make(map[string]Node)
	}
	if _, ok := r.Root.Children[path]; ok == false {
		// case when path is a sub path of existing paths
		if status, res := contains(path, r.Root.Children); status == "ok" && len(res) > 0 {
			res[0].Children[key] = createNode(res[0].Path + path)
		} else { // else add the path, it's a new one.
			r.Root.Children[key] = createNode(r.Root.Path + path)
		}
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
		Name: "replica-" + strconv.Itoa(id),
		Root: createNode(strconv.Itoa(id)),
	}
}

func createNode(path string) Node {
	return Node{Path: "/election", Children: make(map[string]Node)}
}

func createListeners(r Replica) []Listener {
	return []Listener{Listener{
		Destination: r,
		Ltype:       NEW,
		OnMessage: func(r Replica, c chan Message) {
			m := <-c
			min := r.Count()
			if i, _ := strconv.Atoi(m.Body); min == i {
				fmt.Println(i)
			}
			r.Add(m.Body)
		},
	},
	}
}

func OnMessage(listener Listener, m chan Message) {
	go listener.OnMessage(listener.Destination, m)
}

func main() {

	network := Network{
		Quorum:  5,
		Replica: make([]Replica, 5),
		SendEvent: func(m Message, listeners map[string][]Listener) {
			channel := make(chan Message)
			for _, v := range listeners {
				for _, e := range v {
					if e.Ltype == m.Mtype {
						OnMessage(e, channel)
					}
				}
			}
		}}

	ReplicaService{
		Network: network,
		Start: func(n Network) {
			for i := 0; i < n.Quorum; i++ {
				n.Replica[i] = createReplica(i)
				n.Listeners[n.Replica[i].Name] = createListeners(n.Replica[i])
			}
		},
		Execute: func(n Network) {
			for {

				fmt.Println("iteration ... execute")
			}
		},
		Close: func(n Network) {
			// stop
		},
	}.Do()

	fmt.Println(network)

}
