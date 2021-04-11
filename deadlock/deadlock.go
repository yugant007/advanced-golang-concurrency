package main

import (
	"fmt"
	"sync"
	"time"
)

// A deadlocked program is one in which all concurrent processes are waiting on one another.
// In this state, the program will never recover without outside intervention.

type value struct {
	mu    sync.Mutex
	value int
}

var wg sync.WaitGroup

var printSum = func(v1, v2 *value) {
	defer wg.Done()
	v1.mu.Lock() // 1
	defer v1.mu.Unlock() // 2

	time.Sleep(2 * time.Second) // 3
	v2.mu.Lock()
	defer v2.mu.Unlock()

	fmt.Printf("sum=%v\n", v1.value+v2.value)
}


// -------------------------------
//    thread 1  ---> a -------   <---
//    thread 2  ---> b <-----|  -----|
// -------------------------------

func main() {
	// when we run this program we will get error saying
	// fatal error: all goroutines are asleep - deadlock!
	var a, b value
	wg.Add(2)

	// Essentially, we have created two gears that cannot turn together: our first
	// call to printSum locks a and then attempts to lock b, but in the meantime our second call to printSum
	// has locked b and has attempted to lock a. Both goroutines wait infinitely on each other.
	go printSum(&a, &b)
	go printSum(&b, &a)
	wg.Wait()
}
