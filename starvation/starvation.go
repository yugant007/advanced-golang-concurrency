package main

import (
	"fmt"
	"sync"
	"time"
)

// Starvation is any situation where a concurrent process cannot get all the resources it needs to perform work.

// More broadly, starvation usually implies that there are one or more greedy
// concurrent process that are unfairly preventing one or more concurrent
// processes from accomplishing work as efficiently as possible, or maybe at all.

// this program outputs
// - Greedy worker was able to execute 426768 work loops
// - Polite worker was able to execute 356949 work loops.

var wg sync.WaitGroup
var sharedLock sync.Mutex

const runtime = 1 * time.Second

// The greedy worker greedily holds onto the shared lock for the entirety of its
// work loop, whereas the polite worker attempts to only lock when it needs to.
// Both workers do the same amount of simulated work (sleeping for three nanoseconds),
// but as you can see in the same amount of time, the greedy worker got almost twice the amount of work done!

var greedyWorker = func() {
	defer wg.Done()

	var count int
	for begin := time.Now(); time.Since(begin) <= runtime; {
		sharedLock.Lock()
		time.Sleep(3 * time.Nanosecond)
		sharedLock.Unlock()
		count++
	}

	fmt.Printf("Greedy worker was able to execute %v work loops\n", count)
}

var politeWorker = func() {
	defer wg.Done()

	var count int
	for begin := time.Now(); time.Since(begin) <= runtime; {
		sharedLock.Lock()
		time.Sleep(1 * time.Nanosecond)
		sharedLock.Unlock()

		sharedLock.Lock()
		time.Sleep(1 * time.Nanosecond)
		sharedLock.Unlock()

		sharedLock.Lock()
		time.Sleep(1 * time.Nanosecond)
		sharedLock.Unlock()

		count++
	}

	fmt.Printf("Polite worker was able to execute %v work loops.\n", count)
}

func main() {
	wg.Add(2)
	go greedyWorker()
	go politeWorker()

	wg.Wait()
}
