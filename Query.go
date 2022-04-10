package main

import (
	"fmt"
	"sync"
	"time"
)

type Query struct {
	key                                string
	nodeResponse                       map[string]int
	nodeResponseMutex                  *sync.RWMutex
	nodeResponseSizeList               []int
	chatterSize                        int
	maxProcessFrequency                int
	statusCheckFrequencyInMilliSeconds int
	printInfo                          bool
}

func NewQuery(key string, chatterSize, maxProcessFrequency, statusCheckFrequencyInMilliSeconds int, printInfo bool) *Query {
	query := &Query{
		key:                                key,
		nodeResponse:                       make(map[string]int),
		nodeResponseMutex:                  &sync.RWMutex{},
		nodeResponseSizeList:               []int{},
		chatterSize:                        chatterSize,
		maxProcessFrequency:                maxProcessFrequency,
		statusCheckFrequencyInMilliSeconds: statusCheckFrequencyInMilliSeconds,
		printInfo:                          printInfo,
	}

	query.statusDaemon()

	return query
}

func (query *Query) UpdateNodeResponse(node *Node) bool {

	if _, ok := query.nodeResponse[node.name]; ok {
		return false
	}

	query.nodeResponseMutex.Lock()
	defer query.nodeResponseMutex.Unlock()

	if _, ok := query.nodeResponse[node.name]; !ok {
		query.nodeResponse[node.name] = node.Get(query.key)
		query.nodeResponseSizeList = append(query.nodeResponseSizeList, len(query.nodeResponse))
		return true
	}

	return false
}

func (query *Query) statusDaemon() {
	go func() {
		for {
			time.Sleep(time.Millisecond * time.Duration(query.statusCheckFrequencyInMilliSeconds))

			if query.printInfo {
				fmt.Println(time.Now().String(), "|", "query of key:", query.key, "response:", query.nodeResponse)
			}

			fmt.Println(time.Now().String(), "|", "aggregrate of key:", query.key, "val:", query.getAggregrate())

			if query.hasConverged() {
				return
			}
		}
	}()
}

func (query *Query) hasConverged() bool {
	return false
}

func (query *Query) getAggregrate() int {
	query.nodeResponseMutex.RLock()
	defer query.nodeResponseMutex.RUnlock()

	val := 0

	for _, res := range query.nodeResponse {
		val += res
	}

	return val
}

/*

a
b
c
d
e

a -> c, d {a}

c -> b, e {a, c}

d -> a, b {a, c, d}

b -> a, d {a, c, d, b}

e -> b, c {a, c, d, b, e}






*/
