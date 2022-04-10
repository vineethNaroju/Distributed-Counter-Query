package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Node struct {
	name                string
	store               map[string]int
	storeMutex          *sync.RWMutex
	nodeListChannel     chan []*Node
	nodeList            []*Node
	nodeListMutex       *sync.RWMutex
	queryChannel        chan *Query
	queryFrequency      map[*Query]int
	queryFrequencyMutex *sync.RWMutex
}

func NewNode(name string) *Node {
	node := &Node{
		name:                name,
		store:               make(map[string]int),
		storeMutex:          &sync.RWMutex{},
		nodeListChannel:     make(chan []*Node),
		nodeListMutex:       &sync.RWMutex{},
		queryChannel:        make(chan *Query),
		queryFrequency:      make(map[*Query]int),
		queryFrequencyMutex: &sync.RWMutex{},
	}

	node.consumeNodeListDaemon()
	node.queryLoop()

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

func (node *Node) Query(query *Query) {
	go func() {
		node.queryChannel <- query
	}()
}

func (node *Node) queryLoop() {
	go func() {
		for {
			query := <-node.queryChannel

			node.queryFrequencyMutex.Lock()

			val, ok := node.queryFrequency[query]

			if ok {
				if 1+val > query.maxProcessFrequency {
					node.queryFrequencyMutex.Unlock()
					continue
				} else {
					node.queryFrequency[query] += 1
				}
			} else {
				node.queryFrequency[query] = 1
			}

			node.queryFrequencyMutex.Unlock()

			go query.UpdateNodeResponse(node)
			go node.publishQueryToRandomNodes(query)
		}
	}()
}

func (node *Node) publishQueryToRandomNodes(query *Query) {

	if len(node.nodeList) <= 1 {
		return
	}

	node.nodeListMutex.RLock()

	nodes := node.nodeList

	node.nodeListMutex.RUnlock()

	for i := 0; i < query.chatterSize; i++ {
		go func(nd *Node) {
			nd.queryChannel <- query
		}(nodes[rand.Intn(len(nodes))])
	}
}
