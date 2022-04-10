package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func Loadtest() {
	maxNodes := 10
	incReq, uniqKeys := 100000, 1000

	tracker := NewTracker()

	nodes := []*Node{}

	for i := 0; i < maxNodes; i++ {
		nd := NewNode(fmt.Sprint("node-", i))
		tracker.AddNode(nd)
		nodes = append(nodes, nd)
	}

	go func() {
		for i := 0; i < incReq; i++ {
			nodeIdx, keyIdx := rand.Intn(maxNodes), rand.Intn(uniqKeys)
			nodes[nodeIdx].Inc(fmt.Sprint("key-", keyIdx), keyIdx)
		}
	}()

	queryReq := 100
	chatterSize := 3                   // talks with 3 random nodes
	maxProcessFrequency := 3           // aging out in 3 repeatitions
	statusCheckFrequencyInSeconds := 5 // 5 seconds

	wg := &sync.WaitGroup{}
	wg.Add(1)

	for i := 0; i < queryReq; i++ {
		nodeIdx, keyIdx := rand.Intn(maxNodes), rand.Intn(uniqKeys)
		query := NewQuery(fmt.Sprint("key-", keyIdx), chatterSize, maxProcessFrequency, statusCheckFrequencyInSeconds)
		nodes[nodeIdx].Query(query)
	}

	wg.Wait()
}
