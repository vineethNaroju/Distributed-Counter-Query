package main

import (
	"sync"
	"time"
)

func Demo() {

	printInfo := true

	tracker := NewTracker()

	// nodes
	bob, mary, lisa, jane := NewNode("bob", printInfo), NewNode("mary", printInfo), NewNode("lisa", printInfo), NewNode("jane", printInfo)

	queryKey := "a"                           // inc / get for key "a"
	chatterSize := 2                          // each node publishes to atmost 2 random un-visited nodes
	maxProcessFrequency := 1                  // each node will process this query atmost 1 time
	statusCheckFrequencyInMilliSeconds := 500 // display results for every 500ms

	tracker.AddNode(bob)
	time.Sleep(time.Second * 2)

	tracker.AddNode(mary)
	time.Sleep(time.Second * 2)

	tracker.AddNode(lisa)
	time.Sleep(time.Second * 2)

	tracker.AddNode(jane)
	time.Sleep(time.Second * 2)

	bob.Inc(queryKey, 10)
	mary.Inc(queryKey, 5)
	lisa.Inc(queryKey, 20)
	jane.Inc(queryKey, 2)

	a := NewQuery(queryKey, chatterSize, maxProcessFrequency, statusCheckFrequencyInMilliSeconds, printInfo)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	bob.Query(a)

	wg.Wait() // to simply halt main go-routine

}
