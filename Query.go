package main

import (
	"fmt"
	"sync"
	"time"
)

type Query struct {
	key                           string
	nodeResponse                  map[string]int
	nodeResponseMutex             *sync.RWMutex
	nodeResponseSizeList          []int
	chatterSize                   int
	maxProcessFrequency           int
	statusCheckFrequencyInSeconds int
}

func NewQuery(key string, chatterSize, maxProcessFrequency, statusCheckFrequencyInSeconds int) *Query {
	query := &Query{
		key:                           key,
		nodeResponse:                  make(map[string]int),
		nodeResponseMutex:             &sync.RWMutex{},
		nodeResponseSizeList:          []int{},
		chatterSize:                   chatterSize,
		maxProcessFrequency:           maxProcessFrequency,
		statusCheckFrequencyInSeconds: statusCheckFrequencyInSeconds,
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
			time.Sleep(time.Second * time.Duration(query.statusCheckFrequencyInSeconds))

			fmt.Print(query.nodeResponse)

			// for _, sz := range query.nodeResponseSizeList {
			// 	fmt.Print(sz)
			// 	fmt.Print(" ")
			// }

			fmt.Println()
		}
	}()
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
