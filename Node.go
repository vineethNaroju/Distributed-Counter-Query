package main

import (
	"fmt"
	"sync"
	"time"
)

type Node struct {
	name            string
	store           map[string]int
	storeMutex      *sync.RWMutex
	nodeListChannel chan []*Node
	nodeList        []*Node
	nodeListMutex   *sync.RWMutex
	queryChannel    chan *Query
}

func NewNode(name string) *Node {
	node := &Node{
		name:            name,
		store:           make(map[string]int),
		storeMutex:      &sync.RWMutex{},
		nodeListChannel: make(chan []*Node),
		nodeListMutex:   &sync.RWMutex{},
		queryChannel:    make(chan *Query),
	}

	node.consumeNodeListDaemon()

	return node
}

func (node *Node) Get(key string) int {
	node.storeMutex.RLock()
	defer node.storeMutex.RUnlock()

	if val, ok := node.store[key]; ok {
		return val
	}

	return 0
}

func (node *Node) Inc(key string, val int) {
	node.storeMutex.Lock()
	defer node.storeMutex.Unlock()

	if _, ok := node.store[key]; ok {
		node.store[key] += val
	} else {
		node.store[key] = val
	}
}

func (node *Node) consumeNodeListDaemon() {
	go func() {
		for {
			nl := <-node.nodeListChannel

			node.nodeListMutex.Lock()

			node.nodeList = nl

			fmt.Println(time.Now().String() + "| Node " + node.name + " received node list of size " + fmt.Sprint(len(nl)))

			node.nodeListMutex.Unlock()
		}
	}()
}

// TODO
func (node *Node) Query(query *Query) {

}
