package main

import (
	"fmt"
	"sync"
	"time"
)

type Query struct {
	key                                string
	nodeResponse                       map[string]int
	nodeResponseMutex                  *sync.Mutex
	chatterSize                        int
	maxProcessFrequency                int
	statusCheckFrequencyInMilliSeconds int
	printInfo                          bool
}

func NewQuery(key string, chatterSize, maxProcessFrequency, statusCheckFrequencyInMilliSeconds int, printInfo bool) *Query {
	query := &Query{
		key:                                key,
		nodeResponse:                       make(map[string]int),
		nodeResponseMutex:                  &sync.Mutex{},
		chatterSize:                        chatterSize,
		maxProcessFrequency:                maxProcessFrequency,
		statusCheckFrequencyInMilliSeconds: statusCheckFrequencyInMilliSeconds,
		printInfo:                          printInfo,
	}

	query.statusDaemon()

	return query
}

func (query *Query) UpdateNodeResponse(node *Node) bool {
	query.nodeResponseMutex.Lock()
	defer query.nodeResponseMutex.Unlock()

	if _, ok := query.nodeResponse[node.name]; !ok {
		query.nodeResponse[node.name] = node.Get(query.key)
		return true
	}

	return false
}

func (query *Query) statusDaemon() {

	fmt.Println(time.Now().String(), "|", "starting query daemon for key:", query.key)

	go func() {
		for {
			if query.printInfo {
				fmt.Println(time.Now().String(), "|", "query of key:", query.key, "response:", query.nodeResponse)
			}

			fmt.Println(time.Now().String(), "|", "aggregrate of key:", query.key, "val:", query.getAggregrate())

			if query.hasConverged() {
				return
			}

			time.Sleep(time.Millisecond * time.Duration(query.statusCheckFrequencyInMilliSeconds))
		}
	}()
}

func (query *Query) hasConverged() bool {
	return false
}

func (query *Query) getAggregrate() int {
	query.nodeResponseMutex.Lock()
	defer query.nodeResponseMutex.Unlock()

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
