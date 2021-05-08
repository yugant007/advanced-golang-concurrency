package main

import (
	"testing"
	"time"
)

//The DoWork function is a pretty simple generator that converts the numbers we pass
//in to a stream on the channel it returns. Let’s try testing this function. Here’s an example of a bad test:

func TestDoWork_GeneratesAllNumbers(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	_, results := DoWork(done, intSlice...)

	for i, expected := range intSlice {
		select {
		case r := <-results:
			if r != expected {
				t.Errorf(
					"index %v: expected %v, but received %v,",
					i,
					expected,
					r,
				)
			}
		case <-time.After(1 * time.Second):
			t.Fatal("test timed out")
		}
	}
}

func TestDoWork_GeneratesAllNumbers1(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	heartbeat, results := DoWork(done, intSlice...)

	<-heartbeat //Here we wait for the goroutine to signal that it’s beginning to process an iteration.

	i := 0
	for r := range results {
		if expected := intSlice[i]; r != expected {
			t.Errorf("index %v: expected %v, but received %v,", i, expected, r)
		}
		i++
	}
}

func TestDoWork_GeneratesAllNumbers2(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	const timeout = 2*time.Second
	heartbeat, results := DoWork1(done, timeout/2, intSlice...)

	<-heartbeat //We still wait for the first heartbeat to occur to indicate we’ve entered the goroutine’s loop.

	i := 0
	for {
		select {
		case r, ok := <-results:
			if ok == false {
				return
			} else if expected := intSlice[i]; r != expected {
				t.Errorf(
					"index %v: expected %v, but received %v,",
					i,
					expected,
					r,
				)
			}
			i++
		case <-heartbeat: //We also select on the heartbeat here to keep the timeout from occuring.
		case <-time.After(timeout):
			t.Fatal("test timed out")
		}
	}
}
