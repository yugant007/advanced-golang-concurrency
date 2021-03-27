package main

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Livelocks are programs that are actively performing concurrent operations,
// but these operations do nothing to move the state of the program forward.

// Livelocks are a subset of a larger set of problems called starvation.

// example -> Have you ever been in a hallway walking toward another person?
// She moves to one side to let you pass, but you’ve just done the same. So
// you move to the other side, but she’s also done the same. Imagine this
// going on forever, and you understand livelocks.
var cadence = sync.NewCond(&sync.Mutex{})

var peopleInHallway sync.WaitGroup

// For the example to demonstrate a livelock, each person must move at
// the same rate of speed, or cadence. takeStep simulates a constant cadence between all parties

var takeStep = func() {
	// shared lock
	cadence.L.Lock()
	cadence.Wait()
	cadence.L.Unlock()
}

var tryDir = func(dirName string, dir *int32, out *bytes.Buffer) bool {
	_, _ = fmt.Fprintf(out, " %v", dirName)
	atomic.AddInt32(dir, 1)
	takeStep()
	if atomic.LoadInt32(dir) == 1 {
		_, _ = fmt.Fprint(out, ". Success!")
		return true
	}
	takeStep()
	atomic.AddInt32(dir, -1)
	return false
}

var left, right int32
var tryLeft = func(out *bytes.Buffer) bool {
	return tryDir("left", &left, out)
}

var tryRight = func(out *bytes.Buffer) bool {
	return tryDir("right", &right, out)
}

var walk = func(walking *sync.WaitGroup, name string) {
	var out bytes.Buffer
	defer func() { fmt.Println(out.String()) }()
	defer walking.Done()
	_, _ = fmt.Fprintf(&out, "%v is trying to scoot:", name)
	for i := 0; i < 5; i++ {
		if tryLeft(&out) || tryRight(&out) {
			return
		}
	}
	_, _ = fmt.Fprintf(&out, "\n%v tosses her hands up in exasperation!", name)
}

func main() {
	go func() {
		for range time.Tick(1 * time.Millisecond) {
			cadence.Broadcast()
		}
	}()

	peopleInHallway.Add(2)
	go walk(&peopleInHallway, "Alice")
	go walk(&peopleInHallway, "Barbara")
	peopleInHallway.Wait()
}
