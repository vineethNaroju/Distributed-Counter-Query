package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func generatePerpetualTraffic(nodeList []*Node, key string) {
	for {
		time.Sleep(time.Millisecond * 1)
		id := rand.Intn(len(nodeList))
		nodeList[id].Inc(key, id)
	}
}

func Loadtest() {
	// wtf to test ?
	// 1000 nodes

	printInfo := false // few print outs
	nodeCount := 1000  // cluster node size
	queryKey := "abc"
	chatterSize := 10                       // speaks with atmost 5 random nodes
	maxProcessFrequency := 5                // a node will process this query atmost 5 times
	statusCheckFrequencyInMilliSeconds := 5 // will print status every 10 ms

	tracker := NewTracker()

	nodeList := []*Node{}

	for i := 0; i < nodeCount; i++ {
		nd := NewNode(fmt.Sprint("node-", i), printInfo)
		tracker.AddNode(nd)
		nodeList = append(nodeList, nd)
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go generatePerpetualTraffic(nodeList, queryKey)

	time.Sleep(time.Second * 2)

	query := NewQuery(queryKey, chatterSize, maxProcessFrequency, statusCheckFrequencyInMilliSeconds, printInfo)

	nodeList[nodeCount/2-1].Query(query)

	wg.Wait()
}
