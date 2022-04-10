package main

import (
	"fmt"
	"sync"
	"time"
)

type Tracker struct {
	nodeStore      map[string]*Node
	nodeStoreMutex *sync.RWMutex
}

func NewTracker() *Tracker {

	tracker := &Tracker{
		nodeStore:      make(map[string]*Node),
		nodeStoreMutex: &sync.RWMutex{},
	}

	tracker.publishNodeListDaemon()

	return tracker
}

func (tracker *Tracker) AddNode(node *Node) error {
	tracker.nodeStoreMutex.Lock()
	defer tracker.nodeStoreMutex.Unlock()

	if _, ok := tracker.nodeStore[node.name]; ok {
		return fmt.Errorf("duplicate node name " + node.name)
	}

	tracker.nodeStore[node.name] = node
	return nil
}

func (tracker *Tracker) publishNodeListDaemon() {
	go func() {
		for {
			time.Sleep(time.Second * 1)

			nodeList := []*Node{}

			tracker.nodeStoreMutex.RLock()

			for _, val := range tracker.nodeStore {
				nodeList = append(nodeList, val)
			}

			tracker.nodeStoreMutex.RUnlock()

			for _, val := range nodeList {
				go func(nd *Node) {
					nd.nodeListChannel <- nodeList
				}(val)
			}
		}
	}()
}
