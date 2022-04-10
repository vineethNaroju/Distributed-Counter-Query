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
		nodeList[rand.Intn(len(nodeList))].Inc(key, rand.Intn(10000))
	}
}

func Loadtest() {
	// wtf to test ?
	// 1000 nodes

	printInfo := false // few print outs
	nodeCount := 1000  // cluster node size
	queryKey := "abc"
	chatterSize := 5                          // speaks with atmost 5 random nodes
	maxProcessFrequency := 3                  // a node will process this query atmost 3 times
	statusCheckFrequencyInMilliSeconds := 500 // will print status every 500 ms

	tracker := NewTracker()

	nodeList := []*Node{}

	for i := 0; i < nodeCount; i++ {
		nd := NewNode(fmt.Sprint("node-", i), printInfo)
		tracker.AddNode(nd)
		nodeList = append(nodeList, nd)
	}

	time.Sleep(time.Second * 2)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go generatePerpetualTraffic(nodeList, queryKey)

	query := NewQuery(queryKey, chatterSize, maxProcessFrequency, statusCheckFrequencyInMilliSeconds, printInfo)

	nodeList[nodeCount/2-1].Query(query)

	wg.Wait()
}
