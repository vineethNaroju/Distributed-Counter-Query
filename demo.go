package main

import (
	"sync"
	"time"
)

func Demo() {
	tracker := NewTracker()

	bob, mary, lisa, jane := NewNode("bob"), NewNode("mary"), NewNode("lisa"), NewNode("jane")

	chatterSize, maxProcessFrequency := 2, 1

	tracker.AddNode(bob)
	time.Sleep(time.Second * 2)

	tracker.AddNode(mary)
	time.Sleep(time.Second * 2)

	tracker.AddNode(lisa)
	time.Sleep(time.Second * 2)

	tracker.AddNode(jane)
	time.Sleep(time.Second * 2)

	bob.Inc("a", 10)
	mary.Inc("a", 5)
	lisa.Inc("a", 20)
	jane.Inc("a", 2)

	time.Sleep(time.Second * 2)
	a := NewQuery("a", chatterSize, maxProcessFrequency, 3)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	bob.Query(a)

	wg.Wait()

}
