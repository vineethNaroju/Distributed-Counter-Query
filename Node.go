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
	storeMutex          *sync.Mutex
	nodeListChannel     chan []*Node
	nodeList            []*Node
	nodeListMutex       *sync.Mutex
	queryChannel        chan *Query
	queryFrequency      map[*Query]int
	queryFrequencyMutex *sync.Mutex
	printInfo           bool
}

func NewNode(name string, printInfo bool) *Node {
	node := &Node{
		name:                name,
		store:               make(map[string]int),
		storeMutex:          &sync.Mutex{},
		nodeListChannel:     make(chan []*Node),
		nodeListMutex:       &sync.Mutex{},
		queryChannel:        make(chan *Query),
		queryFrequency:      make(map[*Query]int),
		queryFrequencyMutex: &sync.Mutex{},
		printInfo:           printInfo,
	}

	node.consumeNodeListDaemon()
	node.queryLoop()

	return node
}

func (node *Node) Get(key string) int {
	node.storeMutex.Lock()
	defer node.storeMutex.Unlock()

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

	if node.printInfo {
		fmt.Println(time.Now().String(), "|", node.name, "inc", key, "by", val, "results in", node.store[key])
	}
}

func (node *Node) consumeNodeListDaemon() {
	go func() {
		for {
			nl := <-node.nodeListChannel

			node.nodeListMutex.Lock()

			node.nodeList = nl

			if node.printInfo {
				fmt.Println(time.Now().String(), "|", node.name, "received node list", node.nodeNames(nl))
			}

			node.nodeListMutex.Unlock()

			time.Sleep(time.Millisecond * 10)
		}
	}()
}

func (node *Node) nodeNames(nodeList []*Node) []string {
	names := []string{}

	for _, val := range nodeList {
		names = append(names, val.name)
	}

	return names
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

			if node.printInfo {
				fmt.Println(time.Now().String(), "|", node.name, "query:", query.key, "nodeResponse:", query.nodeResponse, "node-query-freq:", val, "node-processed-query:", ok)
			}

			shouldProcess := true

			if ok && 1+val <= query.maxProcessFrequency {
				node.queryFrequency[query] += 1
			} else if !ok {
				node.queryFrequency[query] = 1
			} else {
				shouldProcess = false
			}

			node.queryFrequencyMutex.Unlock()

			if shouldProcess {
				query.UpdateNodeResponse(node)
				go node.publishQueryToRandomNodes(query)
			}

			time.Sleep(time.Millisecond * 10)
		}
	}()
}

func (node *Node) publishQueryToRandomNodes(query *Query) {

	node.nodeListMutex.Lock()

	if len(node.nodeList) <= 1 {
		node.nodeListMutex.Unlock()
		return
	}

	nodes := node.nodeList

	node.nodeListMutex.Unlock()

	notQueriedNodes := []*Node{}

	query.nodeResponseMutex.Lock()

	for _, val := range nodes {
		if _, ok := query.nodeResponse[val.name]; !ok {
			notQueriedNodes = append(notQueriedNodes, val)
		}
	}

	query.nodeResponseMutex.Unlock()

	cs := query.chatterSize

	queryNodes := []*Node{}

	for len(notQueriedNodes) > 0 && cs > 0 {

		cs--
		n := len(notQueriedNodes)
		rid := rand.Intn(n)

		queryNodes = append(queryNodes, notQueriedNodes[rid])

		notQueriedNodes[rid], notQueriedNodes[n-1] = notQueriedNodes[n-1], notQueriedNodes[rid]

		notQueriedNodes = notQueriedNodes[:n-1]
	}

	for i := 0; i < len(queryNodes); i++ {

		if node.printInfo {
			fmt.Println(time.Now().String(), "|", node.name, "publish query:", query.key, "to", queryNodes[i].name)
		}

		go func(nd *Node) {
			nd.queryChannel <- query
		}(queryNodes[i])
	}
}
