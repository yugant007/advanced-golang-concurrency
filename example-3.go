package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	//The for-select Loop
	// Something you’ll see over and over again in Go programs is the for-select loop. It’s nothing more
	//than something like this:
	//for { // Either loop infinitely or range over something
	//	select {
	//	// Do some work with channels
	//	}
	//}
	//There are a couple of different scenarios where you’ll see this pattern pop up.
	//1. Sending iteration variables out on a channel
	// Oftentimes you’ll want to convert something that can be iterated over into values on a channel
	//. This is nothing fancy, and usually looks something like this:
	//for _, s := range []string{"a", "b", "c"} {
	//	select {
	//	case <-done:
	//		return
	//	case stringStream <- s:
	//	}
	//}
	//2. Looping infinitely waiting to be stopped
	// It’s very common to create goroutines that loop infinitely until they’re stopped. There are
	//a couple variations of this one. Which one you choose is purely a stylistic preference.
	//The first variation keeps the select statement as short as possible:
	//for {
	//	select {
	//	case <-done:
	//		return
	//	default:
	//	}
	//
	//	// Do non-preemptable work
	//}
	//If the done channel isn’t closed, we’ll exit the select statement and continue on to the rest
	//of our for loop’s body.
	//The second variation embeds the work in a default clause of the select statement:
	//for {
	//	select {
	//	case <-done:
	//		return
	//	default:
	//		// Do non-preemptable work
	//	}
	//}

	//Preventing Goroutine Leaks
	//we know goroutines are cheap and easy to create; it’s one of the things that makes Go
	//such a productive language. The runtime handles multiplexing the goroutines onto any
	//number of operating system threads so that we don’t often have to worry about that
	//level of abstraction. But they do cost resources, and goroutines are not garbage
	//collected by the runtime, so regardless of how small their memory footprint is,
	//we don’t want to leave them lying about our process. So how do we go about ensuring they’re cleaned up?
	//Let’s start from the beginning and think about this step by step: why would a goroutine
	//exist? we know that goroutines represent units of work that may or
	//may not run in parallel with each other. The goroutine has a few paths to termination:
	//1. When it has completed its work.
	//2. When it cannot continue its work due to an unrecoverable error.
	//3. When it’s told to stop working.
	//We get the first two paths for free—these paths are your algorithm—but what about
	//work cancellation? This turns out to be the most important bit because of the
	//network effect: if you’ve begun a goroutine, it’s most likely cooperating with
	//several other goroutines in some sort of organized fashion. We could even represent
	//this interconnectedness as a graph: whether or not a child goroutine should continue
	//executing might be predicated on knowledge of the state of many other goroutines.
	//The parent goroutine (often the main goroutine) with this full contextual knowledge
	//should be able to tell its child goroutines to terminate. We’ll continue looking at
	//large-scale goroutine interdependence in the next chapter, but for now let’s consider
	//how to ensure a single child goroutine is guaranteed to be cleaned up. Let’s start with
	//a simple example of a goroutine leak:
	//doWork := func(strings <-chan string) <-chan interface{} {
	//	completed := make(chan interface{})
	//	go func() {
	//		defer fmt.Println("doWork exited.")
	//		defer close(completed)
	//		for s := range strings {
	//			// Do something interesting
	//			fmt.Println(s)
	//		}
	//	}()
	//	return completed
	//}
	//
	//doWork(nil)
	//// Perhaps more work is done here
	//fmt.Println("Done.")
	//Here we see that the main goroutine passes a nil channel into doWork. Therefore, the strings channel
	//will never actually gets any strings written onto it, and the goroutine containing doWork will
	//remain in memory for the lifetime of this process (we would even deadlock if we joined the
	//goroutine within doWork and the main goroutine).
	//The way to successfully mitigate this is to establish a signal between the parent goroutine
	//and its children that allows the parent to signal cancellation to its children. By convention,
	//this signal is usually a read-only channel named done. The parent goroutine passes this channel
	//to the child goroutine and then closes the channel when it wants to cancel the child goroutine.
	//Here’s an example:
	//doWork := func(
	//	done <-chan interface{},
	//	strings <-chan string,
	//) <-chan interface{} { //Here we pass the done channel to the doWork function. As a convention, this channel is the first parameter.
	//	terminated := make(chan interface{})
	//	go func() {
	//		defer fmt.Println("doWork exited.")
	//		defer close(terminated)
	//		for {
	//			select {
	//			case s := <-strings:
	//				// Do something interesting
	//				fmt.Println(s)
	//			case <-done: //On this line we see the ubiquitous for-select pattern in use. One of our case statements is checking whether our done channel has been signaled. If it has, we return from the goroutine.
	//				return
	//			}
	//		}
	//	}()
	//	return terminated
	//}
	//
	//done := make(chan interface{})
	//terminated := doWork(done, nil)
	//
	//go func() { //Here we create another goroutine that will cancel the goroutine spawned in doWork if more than one second passes.
	//	// Cancel the operation after 1 second.
	//	time.Sleep(1 * time.Second)
	//	fmt.Println("Canceling doWork goroutine...")
	//	close(done)
	//}()
	//
	//<-terminated //This is where we join the goroutine spawned from doWork with the main goroutine.
	//fmt.Println("Done.")
	//You can see that despite passing in nil for our strings channel, our goroutine still exits
	//successfully. Unlike the example before it, in this example we do join the two goroutines,
	//and yet do not receive a deadlock. This is because before we join the two goroutines, we create
	//a third goroutine to cancel the goroutine within doWork after a second. We have successfully
	//eliminated our goroutine leak!
	//The previous example handles the case for goroutines receiving on a channel nicely, but what if
	//we’re dealing with the reverse situation: a goroutine blocked on attempting to write a value
	//to a channel? Here’s a quick example to demonstrate the issue:
	//newRandStream := func() <-chan int {
	//	randStream := make(chan int)
	//	go func() {
	//		defer fmt.Println("newRandStream closure exited.") //Here we print out a message when the goroutine successfully terminates.
	//		defer close(randStream)
	//		for {
	//			randStream <- rand.Int()
	//		}
	//	}()
	//
	//	return randStream
	//}
	//
	//randStream := newRandStream()
	//fmt.Println("3 random ints:")
	//for i := 1; i <= 3; i++ {
	//	fmt.Printf("%d: %d\n", i, <-randStream)
	//}
	//You can see from the output that the deferred fmt.Println statement never gets run. After the third
	//iteration of our loop, our goroutine blocks trying to send the next random integer to a channel that
	//is no longer being read from. We have no way of telling the producer it can stop. The solution, just
	//like for the receiving case, is to provide the producer goroutine with a channel informing it to exit:
	//newRandStream := func(done <-chan interface{}) <-chan int {
	//	randStream := make(chan int)
	//	go func() {
	//		defer fmt.Println("newRandStream closure exited.")
	//		defer close(randStream)
	//		for {
	//			select {
	//			case randStream <- rand.Int():
	//			case <-done:
	//				return
	//			}
	//		}
	//	}()
	//
	//	return randStream
	//}
	//
	//done := make(chan interface{})
	//randStream := newRandStream(done)
	//fmt.Println("3 random ints:")
	//for i := 1; i <= 3; i++ {
	//	fmt.Printf("%d: %d\n", i, <-randStream)
	//}
	//close(done)
	//
	//// Simulate ongoing work
	//time.Sleep(1 * time.Second)
	//We see now that the goroutine is being properly cleaned up.
	//Now that we know how to ensure goroutines don’t leak, we can stipulate a convention: If a goroutine is
	//responsible for creating a goroutine, it is also responsible for ensuring it can stop the goroutine.

	//The or-channel
	//At times you may find yourself wanting to combine one or more done channels into a single done channel
	//that closes if any of its component channels close. It is perfectly acceptable, albeit verbose, to write
	//a select statement that performs this coupling; however, sometimes you can’t know the number of done
	//channels you’re working with at runtime. In this case, or if you just prefer a one-liner, you can combine
	//these channels together using the or-channel pattern.
	//var or func(channels ...<-chan interface{}) <-chan interface{}
	//or = func(channels ...<-chan interface{}) <-chan interface{} { //Here we have our function, or, which takes in a variadic slice of channels and returns a single channel.
	//	switch len(channels) {
	//	case 0: //Since this is a recursive function, we must set up termination criteria. The first is that if the variadic slice is empty, we simply return a nil channel. This is consistant with the idea of passing in no channels; we wouldn’t expect a composite channel to do anything.
	//		return nil
	//	case 1: //Our second termination criteria states that if our variadic slice only contains one element, we just return that element.
	//		return channels[0]
	//	}
	//
	//	orDone := make(chan interface{})
	//	go func() { //Here is the main body of the function, and where the recursion happens. We create a goroutine so that we can wait for messages on our channels without blocking.
	//		defer close(orDone)
	//
	//		switch len(channels) {
	//		case 2: //Because of how we’re recursing, every recursive call to or will at least have two channels. As an optimization to keep the number of goroutines constrained, we place a special case here for calls to or with only two channels.
	//			select {
	//			case <-channels[0]:
	//			case <-channels[1]:
	//			}
	//		default: //Here we recursively create an or-channel from all the channels in our slice after the third index, and then select from this. This recurrence relation will destructure the rest of the slice into or-channels to form a tree from which the first signal will return. We also pass in the orDone channel so that when goroutines up the tree exit, goroutines down the tree also exit.
	//			select {
	//			case <-channels[0]:
	//			case <-channels[1]:
	//			case <-channels[2]:
	//			case <-or(append(channels[3:], orDone)...):
	//			}
	//		}
	//	}()
	//	return orDone
	//}
	//This is a fairly concise function that enables you to combine any number of channels together into a
	//single channel that will close as soon as any of its component channels are closed, or written to.
	//Let’s take a look at how we can use this function. Here’s a brief example that takes channels that
	//close after a set duration, and uses the or function to combine these into a single channel that closes:
	//sig := func(after time.Duration) <-chan interface{}{ //This function simply creates a channel that will close when the time specified in the after elapses.
	//	c := make(chan interface{})
	//	go func() {
	//		defer close(c)
	//		time.Sleep(after)
	//	}()
	//	return c
	//}
	//
	//start := time.Now() //Here we keep track of roughly when the channel from the or function begins to block.
	//<-or(
	//	sig(2*time.Hour),
	//	sig(5*time.Minute),
	//	sig(1*time.Second),
	//	sig(1*time.Hour),
	//	sig(1*time.Minute),
	//)
	//fmt.Printf("done after %v", time.Since(start))//And here we print the time it took for the read to occur.
	//If you run this program you will get:
	//done after 1.000216772s

	//Error Handling
	//With concurrent processes, this question becomes a little more complex. Because a concurrent process is
	//operating independently of its parent or siblings, it can be difficult for it to reason about what the
	//right thing to do with the error is.
	//Take a look at the following code for an example of this issue:
	//checkStatus := func(
	//	done <-chan interface{},
	//	urls ...string,
	//) <-chan *http.Response {
	//	responses := make(chan *http.Response)
	//	go func() {
	//		defer close(responses)
	//		for _, url := range urls {
	//			resp, err := http.Get(url)
	//			if err != nil {
	//				fmt.Println(err) //Here we see the goroutine doing its best to signal that there’s an error. What else can it do? It can’t pass it back! How many errors is too many? Does it continue making requests?
	//				continue
	//			}
	//			select {
	//			case <-done:
	//				return
	//			case responses <- resp:
	//			}
	//		}
	//	}()
	//	return responses
	//}
	//
	//done := make(chan interface{})
	//defer close(done)
	//
	//urls := []string{"https://www.google.com", "https://badhost"}
	//for response := range checkStatus(done, urls...) {
	//	fmt.Printf("Response: %v\n", response.Status)
	//}
	//Running this code produces:
	//Response: 200 OK
	//Get https://badhost: dial tcp: lookup badhost on 127.0.1.1:53: no such host
	//Here we see that the goroutine has been given no choice in the matter. It can’t simply swallow the error,
	//and so it does the only sensible thing: it prints the error and hopes something is paying attention.
	//Don’t put your goroutines in this awkward position. I suggest you separate your concerns: in general,
	//your concurrent processes should send their errors to another part of your program that has complete
	//information about the state of your program, and can make a more informed decision about what to do.
	//The following example demonstrates a correct solution to this problem:
	//type Result struct { //Here we create a type that encompasses both the *http.Response and the error possible from an iteration of the loop within our goroutine.
	//	Error error
	//	Response *http.Response
	//}
	//checkStatus := func(done <-chan interface{}, urls ...string) <-chan Result { //This line returns a channel that can be read from to retrieve results of an iteration of our loop.
	//	results := make(chan Result)
	//	go func() {
	//		defer close(results)
	//
	//		for _, url := range urls {
	//			var result Result
	//			resp, err := http.Get(url)
	//			result = Result{Error: err, Response: resp}
	//			select {
	//			case <-done:
	//				return
	//			case results <- result:
	//			}
	//		}
	//	}()
	//	return results
	//}
	//done := make(chan interface{})
	//defer close(done)
	//
	//urls := []string{"https://www.google.com", "https://badhost"}
	//for result := range checkStatus(done, urls...) {
	//	if result.Error != nil {
	//		fmt.Printf("error: %v", result.Error)
	//		continue
	//	}
	//	fmt.Printf("Response: %v\n", result.Response.Status)
	//}
	//In the previous example, we simply wrote errors out to stdio, but we could do something else. Let’s
	//alter our program slightly so that it stops trying to check for status if three or more errors occur:
	//done := make(chan interface{})
	//defer close(done)
	//
	//errCount := 0
	//urls := []string{"a", "https://www.google.com", "b", "c", "d"}
	//for result := range checkStatus(done, urls...) {
	//	if result.Error != nil {
	//		fmt.Printf("error: %v\n", result.Error)
	//		errCount++
	//		if errCount >= 3 {
	//			fmt.Println("Too many errors, breaking!")
	//			break
	//		}
	//		continue
	//	}
	//	fmt.Printf("Response: %v\n", result.Response.Status)
	//}

	//Pipelines
	//A pipeline is just another tool you can use to form an abstraction in your system. In particular, it is a
	//very powerful tool to use when your program needs to process streams, or batches of data. The word pipeline
	//is believed to have first been used in 1856, and likely referred to a line of pipes that transported liquid
	//from one place to another. We borrow this term in computer science because we’re also transporting something
	//from one place to another: data. A pipeline is nothing more than a series of things that take data in, perform
	//an operation on it, and pass the data back out. We call each of these operations a stage of the pipeline.
	//By using a pipeline, you separate the concerns of each stage, which provides numerous benefits. You can
	//modify stages independent of one another, you can mix and match how stages are combined independent of
	//modifying the stages, you can process each stage concurrent to upstream or downstream stages, and you
	//can fan-out, or rate-limit portions of your pipeline. We’ll cover fan-out in the section “Fan-Out, Fan-In”.
	//As mentioned previously, a stage is just something that takes data in, performs a transformation on it,
	//and sends the data back out. Here is a function that could be considered a pipeline stage:
	//multiply := func(values []int, multiplier int) []int {
	//	multipliedValues := make([]int, len(values))
	//	for i, v := range values {
	//		multipliedValues[i] = v * multiplier
	//	}
	//	return multipliedValues
	//}
	//This function takes a slice of integers in with a multiplier, loops through them multiplying as it goes,
	//and returns a new transformed slice out. Looks like a boring function, right? Let’s create another stage:
	//add := func(values []int, additive int) []int {
	//	addedValues := make([]int, len(values))
	//	for i, v := range values {
	//		addedValues[i] = v + additive
	//	}
	//	return addedValues
	//}
	//Another boring function! This one just creates a new slice and adds a value to each element.
	//At this point, you might be wondering what makes these two functions pipeline stages and not
	//just functions. Let’s try combining them:
	//ints := []int{1, 2, 3, 4}
	//for _, v := range add(multiply(ints, 2), 1) {
	//	fmt.Println(v)
	//}
	//you work with every day, but because we constructed them to have the properties of a pipeline stage,
	//we’re able to combine them to form a pipeline. That’s interesting; what are the properties of a pipeline stage?
	//A stage consumes and returns the same type.
	//A stage must be reified by the language so that it may be passed around. Functions in Go are reified and fit
	//this purpose nicely.
	//Here, our add and multiply stages satisfy all the properties of a pipeline stage: they both consume
	//a slice of int and return a slice of int, and because Go has reified functions, we can pass add
	//and multiple around. These properties give rise to the interesting properties of pipeline stages
	//we mentioned earlier: namely it becomes very easy to combine our stages at a higher level without
	//modifying the stages themselves.
	//For example, if we wanted to now add an additional stage to our pipeline to multiply by two, we’d simply
	//wrap our previous pipeline in a new multiply stage, like so:
	//ints := []int{1, 2, 3, 4}
	//for _, v := range multiply(add(multiply(ints, 2), 1), 2) {
	//	fmt.Println(v)
	//}
	//Notice how each stage is taking a slice of data and returning a slice of data? These stages are performing
	//what we call batch processing. This just means that they operate on chunks of data all at once instead
	//of one discrete value at a time. There is another type of pipeline stage that performs stream processing.
	//This means that the stage receives and emits one element at a time.
	//multiply := func(value, multiplier int) int {
	//	return value * multiplier
	//}
	//
	//add := func(value, additive int) int {
	//	return value + additive
	//}
	//
	//ints := []int{1, 2, 3, 4}
	//for _, v := range ints {
	//	fmt.Println(multiply(add(multiply(v, 2), 1), 2))
	//}
	//Each stage is receiving and emitting a discrete value, and the memory footprint of our program is back
	//down to only the size of the pipeline’s input. But we had to pull the pipeline down into the body of the
	//for loop and let the range do the heavy lifting of feeding our pipeline. Not only does this limit the
	//reuse of how we feed the pipeline, but as we’ll see later in this section, it also limits our ability
	//to scale. We have other problems too. Effectively, we’re instantiating our pipeline for every iteration
	//of the loop. Though it’s cheap to make function calls, we’re making three function calls for each iteration
	//of the loop. And what about concurrency? I stated earlier that one of the benefits of utilizing pipelines was
	//the ability to process individual stages concurrently, and I mentioned something about fan-out. Where does
	//all that come in?
	//Best Practices for Constructing Pipelines
	//Channels are uniquely suited to constructing pipelines in Go because they fulfill all of our basic
	//requirements. They can receive and emit values, they can safely be used concurrently, they can be
	//ranged over, and they are reified by the language. Let’s take a moment and convert the previous
	//example to utilize channels instead:
	//generator := func(done <-chan interface{}, integers ...int) <-chan int {
	//	intStream := make(chan int)
	//	go func() {
	//		defer close(intStream)
	//		for _, i := range integers {
	//			select {
	//			case <-done:
	//				return
	//			case intStream <- i:
	//			}
	//		}
	//	}()
	//	return intStream
	//}
	//
	//multiply := func(done <-chan interface{}, intStream <-chan int, multiplier int, ) <-chan int {
	//	multipliedStream := make(chan int)
	//	go func() {
	//		defer close(multipliedStream)
	//		for i := range intStream {
	//			select {
	//			case <-done:
	//				return
	//			case multipliedStream <- i * multiplier:
	//			}
	//		}
	//	}()
	//	return multipliedStream
	//}
	//
	//add := func(done <-chan interface{}, intStream <-chan int, additive int, ) <-chan int {
	//	addedStream := make(chan int)
	//	go func() {
	//		defer close(addedStream)
	//		for i := range intStream {
	//			select {
	//			case <-done:
	//				return
	//			case addedStream <- i + additive:
	//			}
	//		}
	//	}()
	//	return addedStream
	//}
	//
	//done := make(chan interface{})
	//defer close(done)
	//
	//intStream := generator(done, 1, 2, 3, 4)
	//pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)
	//
	//for v := range pipeline {
	//	fmt.Println(v)
	//}
	//It looks like we’ve replicated the desired output, but at the cost of having a lot more code. What exactly
	//have we gained? First, let’s examine what we’ve written. We now have three functions instead of two.
	//They all look like they start one goroutine inside their bodies, and use the pattern we established in
	//“Preventing Goroutine Leaks” of taking in a channel to signal that the goroutine should exit.
	//They all look like they return channels, and some of them look like they take in an additional channel
	//as well. Interesting! Let’s start breaking this down further:
	//done := make(chan interface{})
	//defer close(done)
	//The first thing our program does is create a done channel and call close on it in a defer statement.
	//As discussed previously, this ensures our program exits cleanly and never leaks goroutines. Nothing
	//new there.
	//Some Handy Generators
	//I promised earlier I would talk about some fun generators that might be widely useful. As a reminder,
	//a generator for a pipeline is any function that converts a set of discrete values into a stream of
	//values on a channel. Let’s take a look at a generator called repeat:
	//repeat := func(//This function will repeat the values you pass to it infinitely until you tell it to stop.
	//done <-chan interface{},
	//	values ...interface{},
	//) <-chan interface{} {
	//	valueStream := make(chan interface{})
	//	go func() {
	//		defer close(valueStream)
	//		for {
	//			for _, v := range values {
	//				select {
	//				case <-done:
	//					return
	//				case valueStream <- v:
	//				}
	//			}
	//		}
	//	}()
	//	return valueStream
	//}
	//take := func(//This pipeline stage will only take the first num items off of its incoming valueStream and then exit. Together, the two can be very powerful:
	//	done <-chan interface{},
	//	valueStream <-chan interface{},
	//	num int,
	//) <-chan interface{} {
	//	takeStream := make(chan interface{})
	//	go func() {
	//		defer close(takeStream)
	//		for i := 0; i < num; i++ {
	//			select {
	//			case <-done:
	//				return
	//			case takeStream <- <- valueStream:
	//			}
	//		}
	//	}()
	//	return takeStream
	//}
	//done := make(chan interface{})
	//defer close(done)
	//
	//for num := range take(done, repeat(done, 1), 10) {
	//	fmt.Printf("%v ", num)
	//}
	//In this basic example, we create a repeat generator to generate an infinite number of ones, but then only
	//take the first 10. Because the repeat generator’s send blocks on the take stage’s receive, the repeat
	//generator is very efficient. Although we have the capability of generating an infinite stream of ones,
	//we only generate N+1 instances where N is the number we pass into the take stage.
	//We can expand on this. Let’s create another repeating generator, but this time, let’s create one that
	//repeatedly calls a function. Let’s call it repeatFn:
	//repeatFn := func(
	//	done <-chan interface{},
	//	fn func() interface{},
	//) <-chan interface{} {
	//	valueStream := make(chan interface{})
	//	go func() {
	//		defer close(valueStream)
	//		for {
	//			select {
	//			case <-done:
	//				return
	//			case valueStream <- fn():
	//			}
	//		}
	//	}()
	//	return valueStream
	//}
	//take := func(
	//	done <-chan interface{},
	//	valueStream <-chan interface{},
	//	num int,
	//) <-chan interface{} {
	//	takeStream := make(chan interface{})
	//	go func() {
	//		defer close(takeStream)
	//		for i := 0; i < num; i++ {
	//			select {
	//			case <-done:
	//				return
	//			case takeStream <- <-valueStream:
	//			}
	//		}
	//	}()
	//	return takeStream
	//}
	//done := make(chan interface{})
	//defer close(done)
	//
	//rand := func() interface{} { return rand.Int()}
	//
	//for num := range take(done, repeatFn(done, rand), 10) {
	//	fmt.Println(num)
	//}
	//That’s pretty cool—an infinite channel of random integers generated on an as-needed basis!
	//You may be wondering why all of these generators and stages are receiving and sending on channels of
	//interface{}. We could have just as easily written these functions to be specific to a type, or maybe
	//written a Go generator.
	//Empty interfaces are a bit taboo in Go, but for pipeline stages it is my opinion that it’s OK to deal
	//in channels of interface{} so that you can use a standard library of pipeline patterns. As we discussed
	//earlier, a lot of a pipeline’s utility comes from reusable stages. This is best achieved when the stages
	//operate at the level of specificity appropriate to itself. In the repeat and repeatFn generators, the
	//concern is generating a stream of data by looping over a list or operator. With the take stage, the concern
	//is limiting our pipeline. None of these operations require information about the types they’re working on,
	//but instead only require knowledge of the arity of their parameters.
	//When you need to deal in specific types, you can place a stage that performs the type assertion for you.
	//The performance overhead of having an extra pipeline stage (and thus goroutine) and the type assertion
	//are negligible, as we’ll see in just a bit. Here’s a small example that introduces a toString pipeline stage:
	//toString := func(
	//	done <-chan interface{},
	//	valueStream <-chan interface{},
	//) <-chan string {
	//	stringStream := make(chan string)
	//	go func() {
	//		defer close(stringStream)
	//		for v := range valueStream {
	//			select {
	//			case <-done:
	//				return
	//			case stringStream <- v.(string):
	//			}
	//		}
	//	}()
	//	return stringStream
	//}
	//done := make(chan interface{})
	//defer close(done)
	//
	//var message string
	//for token := range toString(done, take(done, repeat(done, "I", "am."), 5)) {
	//	message += token
	//}
	//
	//fmt.Printf("message: %s...", message)
	//So let’s prove to ourselves that the performance cost of generalizing portions of our pipeline is
	//negligible. We’ll write two benchmarking functions: one to test the generic stages, and one to
	//test the type-specific stages:
	// see benchmark_4_test.go
	//You can see that the type-specific stages are twice as fast, but only marginally faster
	//in magnitude. Generally, the limiting factor on your pipeline will either be your generator,
	//or one of the stages that is computationally intensive. If the generator isn’t creating a
	//stream from memory as with the repeat and repeatFn generators, you’ll probably be I/O bound.
	//Reading from disk or the network will likely eclipse the meager performance overhead shown here.
	//If one of your stages is computationally expensive, this will certainly eclipse this performance
	//overhead. If this technique still leaves a bad taste in your mouth, you can always write a Go
	//generator for creating your generator stages. Speaking of one stage being computationally expensive,
	//how can we help mitigate this? Won’t it rate-limit the entire pipeline?
	//For ways to help mitigate this, let’s discuss the fan-out, fan-in technique.

	//Fan-Out, Fan-In
	//So you’ve got a pipeline set up. Data is flowing through your system beautifully, transforming
	//as it makes its way through the stages you’ve chained together. It’s like a beautiful stream;
	//a beautiful, slow stream, and oh my god why is this taking so long?
	//Sometimes, stages in your pipeline can be particularly computationally expensive. When this happens,
	//upstream stages in your pipeline can become blocked while waiting for your expensive stages to complete.
	//Not only that, but the pipeline itself can take a long time to execute as a whole. How can we address this?
	//One of the interesting properties of pipelines is the ability they give you to operate on the stream of
	//data using a combination of separate, often reorderable stages. You can even reuse stages of the pipeline
	//multiple times. Wouldn’t it be interesting to reuse a single stage of our pipeline on multiple goroutines
	//in an attempt to parallelize pulls from an upstream stage? Maybe that would help improve the performance
	//of the pipeline.
	//In fact, it turns out it can, and this pattern has a name: fan-out, fan-in.
	//Fan-out is a term to describe the process of starting multiple goroutines to handle input from
	//the pipeline, and fan-in is a term to describe the process of combining multiple results into one channel.
	//So what makes a stage of a pipeline suited for utilizing this pattern? You might consider
	//fanning out one of your stages if both of the following apply:
	//It doesn’t rely on values that the stage had calculated before.
	//It takes a long time to run.
	//Let’s take a look at an example.
	//toInt := func(
	//	done <-chan interface{},
	//	valueStream <-chan interface{},
	//) <-chan int {
	//	Stream := make(chan int)
	//	go func() {
	//		defer close(Stream)
	//		for {
	//			select {
	//			case <-done:
	//				return
	//			case Stream <- (<-valueStream).(int):
	//			}
	//		}
	//	}()
	//	return Stream
	//}
	//
	//primeFinder := func(
	//	done <-chan interface{},
	//	valueStream <-chan int,
	//) <-chan interface{} {
	//	Stream := make(chan interface{})
	//	go func() {
	//		defer close(Stream)
	//		for {
	//			select {
	//			case <-done:
	//				return
	//			default:
	//				x := <-valueStream
	//				flag := 0
	//				for i := 2; i < x; i++ {
	//					if x%i == 0 {
	//						flag = 1
	//						break
	//					}
	//				}
	//				if flag == 0 {
	//					Stream <- x
	//				}
	//			}
	//		}
	//	}()
	//	return Stream
	//}
	//rand := func() interface{} { return rand.Intn(50000000) }
	//
	//done := make(chan interface{})
	//defer close(done)
	//
	//start := time.Now()
	//
	//randIntStream := toInt(done, repeatFn(done, rand))
	//fmt.Println("Primes:")
	//for prime := range take(done, primeFinder(done, randIntStream), 10) {
	//	fmt.Printf("\t%d\n", prime)
	//}
	//
	//fmt.Printf("Search took: %v", time.Since(start))
	//We’re generating a stream of random numbers, capped at 50,000,000, converting the stream into an integer
	//stream, and then passing that into our primeFinder stage. primeFinder naively begins to attempt to divide
	//the number provided on the input stream by every number below it. If it’s unsuccessful, it passes the value
	//on to the next stage. Certainly, this is a horrible way to try and find prime numbers, but it fulfills our
	//requirement of taking a long time.
	//In our for loop, we range over the found primes, print them out as they come in, and—thanks to our
	//take stage—close the pipeline after 10 primes are found. We then print out how long the search took, and
	//the done channel is closed by a defer statement and the pipeline is torn down.
	//You can see it took roughly 3 seconds to find 10 primes. Not great.
	//we’ll look at how we can fan-out one or more of the stages to chew through slow operations more quickly.
	//This is a relatively simple example, so we only have two stages: random number generation and prime
	//sieving. In a larger program, your pipeline might be composed of many more stages; how do we know which
	//one to fan out? Remember our criteria from earlier: order-independence and duration. Our random integer
	//generator is certainly order-independent, but it doesn’t take a particularly long time to run. The primeFinder
	//stage is also order-independent—numbers are either prime or not—and because of our naive algorithm, it
	//certainly takes a long time to run. It looks like a good candidate for fanning out.
	//fanIn := func(
	//	done <-chan interface{},
	//	channels ...<-chan interface{},
	//) <-chan interface{} {
	//	var wg sync.WaitGroup
	//	multiplexedStream := make(chan interface{})
	//
	//	multiplex := func(c <-chan interface{}) {
	//		defer wg.Done()
	//		for i := range c {
	//			select {
	//			case <-done:
	//				return
	//			case multiplexedStream <- i:
	//			}
	//		}
	//	}
	//
	//	// Select from all the channels
	//	wg.Add(len(channels))
	//	for _, c := range channels {
	//		go multiplex(c)
	//	}
	//
	//	// Wait for all the reads to complete
	//	go func() {
	//		wg.Wait()
	//		close(multiplexedStream)
	//	}()
	//
	//	return multiplexedStream
	//}
	//
	//done := make(chan interface{})
	//defer close(done)
	//
	//start := time.Now()
	//
	//rand := func() interface{} { return rand.Intn(50000000) }
	//
	//randIntStream := toInt(done, repeatFn(done, rand))
	//
	//numFinders := runtime.NumCPU()// Here we’re starting up as many copies of this stage as we have CPUs.
	//fmt.Printf("Spinning up %d prime finders.\n", numFinders)
	//finders := make([]<-chan interface{}, numFinders)
	//fmt.Println("Primes:")
	//for i := 0; i < numFinders; i++ {
	//	finders[i] = primeFinder(done, randIntStream)
	//}
	//
	//for prime := range take(done, fanIn(done, finders...), 10) {
	//	fmt.Printf("\t%d\n", prime)
	//}
	//
	//fmt.Printf("Search took: %v", time.Since(start))
	//And that’s it! We now have eight goroutines pulling from the random number generator and attempting to
	//determine whether the number is prime. Generating random numbers shouldn’t take much time, and so each
	//goroutine for the findPrimes stage should be able to determine whether its number is prime and then have
	//another random number available to it immediately.
	//We still have a problem though: now that we have four goroutines, we also have four channels,
	//but our range over primes is only expecting one channel. This brings us to the fan-in portion of the pattern.
	//As we discussed earlier, fanning in means multiplexing or joining together multiple streams of data
	//into a single stream. The algorithm to do so is relatively simple.
	//In a nutshell, fanning in involves creating the multiplexed channel consumers will read from, and then
	//spinning up one goroutine for each incoming channel, and one goroutine to close the multiplexed channel
	//when the incoming channels have all been closed. Since we’re going to be creating a goroutine that is
	//waiting on N other goroutines to complete, it makes sense to create a sync.WaitGroup to coordinate things.
	//The multiplex function also notifies the WaitGroup that it’s done.
	//A naive implementation of the fan-in, fan-out algorithm only works if the order in which results arrive is
	//unimportant. We have done nothing to guarantee that the order in which items are read from the randIntStream
	//is preserved as it makes its way through the sieve. Later, we’ll look at an example of a way to maintain order.
	//So down from 3 seconds to 1 seconds, not bad!

	//The or-done-channel
	//At times you will be working with channels from disparate parts of your system. Unlike with pipelines,
	//you can’t make any assertions about how a channel will behave when code you’re working with is canceled
	//via its done channel. That is to say, you don’t know if the fact that your goroutine was canceled means
	//the channel you’re reading from will have been canceled. For this reason, as we laid out in “Preventing
	//Goroutine Leaks”, we need to wrap our read from the channel with a select statement that also selects
	//from a done channel. This is perfectly fine, but doing so takes code that’s easily read like this:
	//for val := range myChan {
	//	// Do something with val
	//}
	//And explodes it out into this:
	//loop:
	//for {
	//	select {
	//	case <-done:
	//		break loop
	//	case maybeVal, ok := <-myChan:
	//		if ok == false {
	//			return // or maybe break from for
	//		}
	//		// Do something with val
	//	}
	//}
	//This can get busy quite quickly—especially if you have nested loops. Continuing with the theme of utilizing
	//goroutines to write clearer concurrent code, and not prematurely optimizing, we can fix this with a single
	//goroutine. We encapsulate the verbosity so that others don’t have to:
	//orDone := func(done, c <-chan interface{}) <-chan interface{} {
	//	valStream := make(chan interface{})
	//	go func() {
	//		defer close(valStream)
	//		for {
	//			select {
	//			case <-done:
	//				return
	//			case v, ok := <-c:
	//				if ok == false {
	//					return
	//				}
	//				select {
	//				case valStream <- v:
	//				case <-done:
	//				}
	//			}
	//		}
	//	}()
	//	return valStream
	//}
	//Doing this allows us to get back to simple for loops, like so:
	//for val := range orDone(done, myChan) {
	//    // Do something with val
	//}

	//The tee-channel
	//Sometimes you may want to split values coming in from a channel so that you can send them off into
	//two separate areas of your codebase. Imagine a channel of user commands: you might want to take in
	//a stream of user commands on a channel, send them to something that executes them, and also send
	//them to something that logs the commands for later auditing.
	//Taking its name from the tee command in Unix-like systems, the tee-channel does just this. You can
	//pass it a channel to read from, and it will return two separate channels that will get the same value:
	//tee := func(done <-chan interface{}, in <-chan interface{}, ) (<-chan interface{}, <-chan interface{}) {
	//	out1 := make(chan interface{})
	//	out2 := make(chan interface{})
	//	go func() {
	//		defer close(out1)
	//		defer close(out2)
	//		for val := range orDone(done, in) {
	//			var out1, out2 = out1, out2// We will want to use local versions of out1 and out2, so we shadow these variables.
	//			for i := 0; i < 2; i++ {
	//				select {
	//				case <-done:
	//				case out1 <- val:
	//					out1 = nil//We’re going to use one select statement so that writes to out1 and out2 don’t block each other. To ensure both are written to, we’ll perform two iterations of the select statement: one for each outbound channel.
	//				case out2 <- val:
	//					out2 = nil// Once we’ve written to a channel, we set its shadowed copy to nil so that further writes will block and the other channel may continue.
	//				}
	//			}
	//		}
	//	}()
	//	return out1, out2
	//}
	//Notice that writes to out1 and out2 are tightly coupled. The iteration over in cannot continue until
	//both out1 and out2 have been written to. Usually this is not a problem as handling the throughput of
	//the process reading from each channel should be a concern of something other than the tee command
	//anyway, but it’s worth noting. Here’s a quick example to demonstrate:
	//done := make(chan interface{})
	//defer close(done)
	//
	//out1, out2 := tee(done, take(done, repeat(done, 1, 2), 4))
	//
	//for val1 := range out1 {
	//	fmt.Printf("out1: %v, out2: %v\n", val1, <-out2)
	//}

	//The bridge-channel
	//In some circumstances, you may find yourself wanting to consume values from a sequence of channels:
	//<-chan <-chan interface{}
	//As a consumer, the code may not care about the fact that its values come from a sequence of channels.
	//In that case, dealing with a channel of channels can be cumbersome. If we instead define a function
	//that can destructure the channel of channels into a simple channel—a technique called bridging the
	//channels—this will make it much easier for the consumer to focus on the problem at hand. Here’s how
	//we can achieve that:
	//bridge := func(done <-chan interface{}, chanStream <-chan <-chan interface{}) <-chan interface{} {
	//	valStream := make(chan interface{}) //This is the channel that will return all values from bridge.
	//	go func() {
	//		defer close(valStream)
	//		for { //This loop is responsible for pulling channels off of chanStream and providing them to a nested loop for use.
	//			var stream <-chan interface{}
	//			select {
	//			case maybeStream, ok := <-chanStream:
	//				if ok == false {
	//					return
	//				}
	//				stream = maybeStream
	//			case <-done:
	//				return
	//			}
	//			for val := range orDone(done, stream) { //This loop is responsible for reading values off the channel it has been given and repeating those values onto valStream. When the stream we’re currently looping over is closed, we break out of the loop performing the reads from this channel, and continue with the next iteration of the loop, selecting channels to read from. This provides us with an unbroken stream of values.
	//				select {
	//				case valStream <- val:
	//				case <-done:
	//				}
	//			}
	//		}
	//	}()
	//	return valStream
	//}
	//This is pretty straightforward code. Now we can use bridge to help present a single-channel facade over
	//a channel of channels. Here’s an example that creates a series of 10 channels, each with one element
	//written to them, and passes these channels into the bridge function:
	//genVals := func() <-chan <-chan interface{} {
	//	chanStream := make(chan (<-chan interface{}))
	//	go func() {
	//		defer close(chanStream)
	//		for i := 0; i < 10; i++ {
	//			stream := make(chan interface{}, 1)
	//			stream <- i
	//			close(stream)
	//			chanStream <- stream
	//		}
	//	}()
	//	return chanStream
	//}
	//
	//for v := range bridge(nil, genVals()) {
	//	fmt.Printf("%v ", v)
	//}
	//Thanks to bridge, we can use the channel of channels from within a single range statement and focus
	//on our loop’s logic. Destructuring the channel of channels is left to code that is specific to this concern.

	//The context Package
	//As we’ve seen, in concurrent programs it’s often necessary to preempt operations because of timeouts,
	//cancellation, or failure of another portion of the system. We’ve looked at the idiom of creating a
	//done channel, which flows through your program and cancels all blocking concurrent operations. This
	//works well, but it’s also somewhat limited.
	//It would be useful if we could communicate extra information alongside the simple notification
	//to cancel: why the cancellation was occuring, or whether or not our function has a deadline
	//by which it needs to complete.
	//it turns out that the need to wrap a done channel with this information is very common in systems of any size.
	//If we take a peek into the context package, we see that it’s very simple:
	//var Canceled = errors.New("context canceled")
	//var DeadlineExceeded error = deadlineExceededError{}
	//
	//type CancelFunc
	//type Context
	//
	//func Background() Context
	//func TODO() Context
	//func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
	//func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)
	//func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
	//func WithValue(parent Context, key, val interface{}) Context
	//We’ll revisit these types and functions in a bit, but for now let’s focus on the Context type.
	//This is the type that will flow through your system much like a done channel does. If you use the
	//context package, each function that is downstream from your top-level concurrent call would take
	//in a Context as its first argument. The type looks like this:
	//type Context interface {
	//
	//	// Deadline returns the time when work done on behalf of this
	//	// context should be canceled. Deadline returns ok==false when no
	//	// deadline is set. Successive calls to Deadline return the same
	//	// results.
	//	Deadline() (deadline time.Time, ok bool)
	//
	//	// Done returns a channel that's closed when work done on behalf
	//	// of this context should be canceled. Done may return nil if this
	//	// context can never be canceled. Successive calls to Done return
	//	// the same value.
	//	Done() <-chan struct{}
	//
	//	// Err returns a non-nil error value after Done is closed. Err
	//	// returns Canceled if the context was canceled or
	//	// DeadlineExceeded if the context's deadline passed. No other
	//	// values for Err are defined.  After Done is closed, successive
	//	// calls to Err return the same value.
	//	Err() error
	//
	//	// Value returns the value associated with this context for key,
	//	// or nil if no value is associated with key. Successive calls to
	//	// Value with the same key returns the same result.
	//	Value(key interface{}) interface{}
	//}
	//This also looks pretty simple. There’s a Done method which returns a channel that’s closed when our
	//function is to be preempted. There’s also some new, but easy to understand methods: a Deadline
	//function to indicate if a goroutine will be canceled after a certain time, and an Err method that
	//will return non-nil if the goroutine was canceled. But the Value method looks a little out of place.
	//What’s it for?
	//The Go authors noticed that one of the primary uses of goroutines was programs that serviced requests.
	//Usually in these programs, request-specific information needs to be passed along in addition to
	//information about preemption. This is the purpose of the Value function. We’ll talk about this more
	//in a bit, but for now we just need to know that the context package serves two primary purposes:
	// 1. To provide an API for canceling branches of your call-graph.
	// 2. To provide a data-bag for transporting request-scoped data through your call-graph.
	// Let’s focus on the first aspect: cancellation.
	//As we learned in “Preventing Goroutine Leaks”, cancellation in a function has three aspects:
	//A goroutine’s parent may want to cancel it.
	//A goroutine may want to cancel its children.
	//Any blocking operations within a goroutine need to be preemptable so that it may be canceled.
	//The context package helps manage all three of these.
	//As we mentioned, the Context type will be the first argument to your function. If you look at the methods
	//on the Context interface, you’ll see that there’s nothing present that can mutate the state of the
	//underlying structure. Further, there’s nothing that allows the function accepting the Context to cancel
	//it. This protects functions up the call stack from children canceling the context. Combined with the
	//Done method, which provides a done channel, this allows the Context type to safely manage cancellation
	//from its antecedents.
	//This raises a question: if a Context is immutable, how do we affect the behavior of
	//cancellations in functions below a current function in the call stack?
	//This is where the functions in the context package become important. Let’s take a look
	//at a few of them one more time to refresh our memory:
	//func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
	//func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)
	//func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
	//Notice that all these functions take in a Context and return one as well. Some of these also
	//take in other arguments like deadline and timeout. The functions all generate new instances
	//of a Context with the options relative to these functions.
	//WithCancel returns a new Context that closes its done channel when the returned cancel function is called.
	//WithDeadline returns a new Context that closes its done channel when the machine’s clock advances past
	//the given deadline. WithTimeout returns a new Context that closes its done channel after the given timeout duration.

	//If your function needs to cancel functions below it in the call-graph in some manner, it will call
	//one of these functions and pass in the Context it was given, and then pass the Context returned into
	//its children. If your function doesn’t need to modify the cancellation behavior, the function simply
	//passes on the Context it was given.
	//
	//In this way, successive layers of the call-graph can create a Context that adheres to their needs
	//without affecting their parents. This provides a very composable, elegant solution for how to manage
	//branches of your call-graph.
	//
	//In this spirit, instances of a Context are meant to flow through your program’s call-graph. In an
	//object-oriented paradigm, it’s common to store references to often-used data as member variables,
	//but it’s important to not do this with instances of context.Context. Instances of context.Context
	//may look equivalent from the outside, but internally they may change at every stack-frame. For this
	//reason, it’s important to always pass instances of Context into your functions. This way functions
	//have the Context intended for it, and not the Context intended for a stack-frame N levels up the stack.
	//
	//At the top of your asynchronous call-graph, your code probably won’t have been passed a Context.
	//To start the chain, the context package provides you with two functions to create empty instances of Context:
	//func Background() Context
	//func TODO() Context
	//Background simply returns an empty Context. TODO is not meant for use in production,
	//but also returns an empty Context; TODO’s intended purpose is to serve as a placeholder
	//for when you don’t know which Context to utilize, or if you expect your code to be provided with a Context,
	//but the upstream code hasn’t yet furnished one.
	//
	//So let’s put all this to use. Let’s look at an example that uses the done channel pattern,
	//and see what benefits we might gain from switching to use of the context package. Here is a
	//program that concurrently prints a greeting and a farewell:
	//var wg sync.WaitGroup
	//done := make(chan interface{})
	//defer close(done)
	//
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	if err := printGreeting(done); err != nil {
	//		fmt.Printf("%v", err)
	//		return
	//	}
	//}()
	//
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	if err := printFarewell(done); err != nil {
	//		fmt.Printf("%v", err)
	//		return
	//	}
	//}()
	//
	//wg.Wait()
	//Running this code produces:
	//goodbye world!
	//hello world!
	//Ignoring the race condition (we could receive our farewell before we’re greeted!), we can see that
	//we have two branches of our program running concurrently. We’ve set up the standard preemption
	//method by creating a done channel and passing it down through our call-graph. If we close the done
	//channel at any point in main, both branches will be canceled.
	//By introducing goroutines in main, we’ve opened up the possibility of controlling this program
	//in a few different and interesting ways. Maybe we want genGreeting to time out if it takes too
	//long. Maybe we don’t want genFarewell to invoke locale if we know its parent is going to be
	//canceled soon. At each stack-frame, a function can affect the entirety of the call stack below it.
	//Using the done channel pattern, we could accomplish this by wrapping the incoming
	//done channel in other done channels and then returning if any of them fire, but we wouldn’t
	//have the extra information about deadlines and errors a Context gives us.
	//To make comparing the done channel pattern to the use of the context package easier,
	//let’s represent this program as a tree. Each node in the tree represents an invocation of a function.
	// see image1.png
	//Let’s modify our program to use the context package instead of a done channel. Because we now have
	//the flexibility of a context.Context, we can introduce a fun scenario.
	//Let’s say that genGreeting only wants to wait one second before abandoning the call to locale—a timeout
	//of one second. We also want to build some smart logic into main. If printGreeting is unsuccessful,
	//we also want to cancel our call to printFarewell. After all, it wouldn’t make sense to say goodbye
	//if we don’t say hello!
	//var wg sync.WaitGroup
	//ctx, cancel := context.WithCancel(context.Background()) //Here main creates a new Context with context.Background() and wraps it with context.WithCancel to allow for cancellations.
	//defer cancel()
	//
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//
	//	if err := printGreeting2(ctx); err != nil {
	//		fmt.Printf("cannot print greeting: %v\n", err)
	//		cancel() //On this line, it will cancel the Context if there is an error returned from printGreeting2.
	//	}
	//}()
	//
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	if err := printFarewell2(ctx); err != nil {
	//		fmt.Printf("cannot print farewell: %v\n", err)
	//	}
	//}()
	//
	//wg.Wait()
	//Here are the results of running this code:
	//cannot print greeting: context deadline exceeded
	//cannot print farewell: context canceled
	//Let’s use our call-graph to understand what’s going on. The numbers here correspond to the
	//code callouts in the preceding example.
	//Notice how genGreeting was able to build up a custom context.Context to meet its needs without
	//having to affect its parent’s Context. If genGreeting were to return successfully, and
	//printGreeting needed to make another call, it could do so without leaking information
	//about how genGreeting operated. This composability enables you to write large systems without
	//mixing concerns throughout your call-graph.
	//We can make another improvement on this program: since we know locale takes roughly one minute to
	//run, in locale we can check to see whether we were given a deadline, and if so, whether we’ll
	//meet it. This example demonstrates using the context.Context’s Deadline method to do so:
	//var wg sync.WaitGroup
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//
	//	if err := printGreeting3(ctx); err != nil {
	//		fmt.Printf("cannot print greeting: %v\n", err)
	//		cancel()
	//	}
	//}()
	//
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	if err := printFarewell3(ctx); err != nil {
	//		fmt.Printf("cannot print farewell: %v\n", err)
	//	}
	//}()
	//
	//wg.Wait()
	//Although the difference in this iteration of the program is small, it allows the locale function to
	//fail fast. In programs that may have a high cost for calling the next bit of functionality, this
	//may save a significant amount of time, but at the very least it also allows the function to fail
	//immediately instead of having to wait for the actual timeout to occur. The only catch is that you
	//have to have some idea of how long your subordinate call-graph will take—an exercise that can be very difficult.

	//This brings us to the other half of what the context package provides: a data-bag for a Context to store
	//and retrieve request-scoped data. Remember that oftentimes when a function creates a goroutine and Context,
	//it’s starting a process that will service requests, and functions further down the stack may need information
	//about the request. Here’s an example of how to store data within the Context, and how to retrieve it:
	//ProcessRequest("jane", "abc123")
	//This produces:
	//
	//handling response for jane (abc123)
	//Pretty simple stuff. The only qualifications are that:
	//The key you use must satisfy Go’s notion of comparability; that is, the equality operators
	//== and != need to return correct results when used.
	//Values returned must be safe to access from multiple goroutines.
	//Since both the Context’s key and value are defined as interface{}, we lose Go’s type-safety when
	//attempting to retrieve values. The key could be a different type, or slightly different than the
	//key we provide. The value could be a different type than we’re expecting. For these reasons, the
	//Go authors recommend you follow a few rules when storing and retrieving value from a Context.
	//First, they recommend you define a custom key-type in your package. As long as other packages do the same,
	//this prevents collisions within the Context. As a reminder as to why, let’s take a look at a short program
	//that attempts to store keys in a map that have different types, but the same underlying value:
	//type foo int
	//type bar int
	//
	//m := make(map[interface{}]int)
	//m[foo(1)] = 1
	//m[bar(1)] = 2
	//
	//fmt.Printf("%v", m)
	//This produces:
	//
	//map[1:1 1:2]
	//You can see that though the underlying values are the same, the different type information differentiates
	//them within a map. Since the type you define for your package’s keys is unexported, other packages cannot
	// conflict with keys you generate within your package.
	//Since we don’t export the keys we use to store the data, we must therefore export functions that
	//retrieve the data for us. This works out nicely since it allows consumers of this data to use static,
	//type-safe functions.
	//When you put all of this together, you get something like the following example:
	//ProcessRequest1("jane", "abc123")
	//We now have a type-safe way to retrieve values from the Context, and—if the consumers were in a
	//different package—they wouldn’t know or care what keys were used to store the information. However,
	//this technique does pose a problem.
	//In the previous example, let’s say HandleResponse did live in another package named response,
	//and let’s say the package ProcessRequest lived in a package named process. The process package
	//would have to import the response package to make the call to HandleResponse, but HandleResponse
	//would have no way to access the accessor functions defined in the process package because importing
	//process would form a circular dependency. Because the types used to store the keys in Context are
	//private to the process package, the response package has no way to retrieve this data!
	//This coerces the architecture into creating packages centered around data types that are imported
	//from multiple locations. This certainly isn’t a bad thing, but it’s something to be aware of.

}
//func printGreeting(done <-chan interface{}) error {
//	greeting, err := genGreeting(done)
//	if err != nil {
//		return err
//	}
//	fmt.Printf("%s world!\n", greeting)
//	return nil
//}
//
//func printFarewell(done <-chan interface{}) error {
//	farewell, err := genFarewell(done)
//	if err != nil {
//		return err
//	}
//	fmt.Printf("%s world!\n", farewell)
//	return nil
//}
//
//func genGreeting(done <-chan interface{}) (string, error) {
//	switch locale, err := locale(done); {
//	case err != nil:
//		return "", err
//	case locale == "EN/US":
//		return "hello", nil
//	}
//	return "", fmt.Errorf("unsupported locale")
//}
//
//func genFarewell(done <-chan interface{}) (string, error) {
//	switch locale, err := locale(done); {
//	case err != nil:
//		return "", err
//	case locale == "EN/US":
//		return "goodbye", nil
//	}
//	return "", fmt.Errorf("unsupported locale")
//}
//
//func locale(done <-chan interface{}) (string, error) {
//	select {
//	case <-done:
//		return "", fmt.Errorf("canceled")
//	case <-time.After(1 * time.Minute):
//	}
//	return "EN/US", nil
//}
//func printGreeting2(ctx context.Context) error {
//	greeting, err := genGreeting2(ctx)
//	if err != nil {
//		return err
//	}
//	fmt.Printf("%s world!\n", greeting)
//	return nil
//}
//
//func printFarewell2(ctx context.Context) error {
//	farewell, err := genFarewell2(ctx)
//	if err != nil {
//		return err
//	}
//	fmt.Printf("%s world!\n", farewell)
//	return nil
//}
//
//func genGreeting2(ctx context.Context) (string, error) {
//	ctx, cancel := context.WithTimeout(ctx, 1*time.Second) //Here genGreeting wraps its Context with context.WithTimeout. This will automatically cancel the returned Context after 1 second, thereby canceling any children it passes the Context into, namely locale.
//	defer cancel()
//
//	switch locale, err := locale2(ctx); {
//	case err != nil:
//		return "", err
//	case locale == "EN/US":
//		return "hello", nil
//	}
//	return "", fmt.Errorf("unsupported locale")
//}
//
//func genFarewell2(ctx context.Context) (string, error) {
//	switch locale, err := locale2(ctx); {
//	case err != nil:
//		return "", err
//	case locale == "EN/US":
//		return "goodbye", nil
//	}
//	return "", fmt.Errorf("unsupported locale")
//}
//
//func locale2(ctx context.Context) (string, error) {
//	select {
//	case <-ctx.Done():
//		return "", ctx.Err() //This line returns the reason why the Context was canceled. This error will bubble all the way up to main, which will cause the cancellation at 2.
//	case <-time.After(1 * time.Minute):
//	}
//	return "EN/US", nil
//}

func printGreeting3(ctx context.Context) error {
	greeting, err := genGreeting3(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", greeting)
	return nil
}

func printFarewell3(ctx context.Context) error {
	farewell, err := genFarewell3(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", farewell)
	return nil
}

func genGreeting3(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	switch locale, err := locale3(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func genFarewell3(ctx context.Context) (string, error) {
	switch locale, err := locale3(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func locale3(ctx context.Context) (string, error) {
	if deadline, ok := ctx.Deadline(); ok { //Here we check to see whether our Context has provided a deadline. If it did, and our system’s clock has advanced past the deadline, we simply return with a special error defined in the context package, DeadlineExceeded.
		if deadline.Sub(time.Now().Add(1*time.Minute)) <= 0 {
			return "", context.DeadlineExceeded
		}
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(1 * time.Minute):
	}
	return "EN/US", nil
}
func ProcessRequest(userID, authToken string) {
	ctx := context.WithValue(context.Background(), "userID", userID)
	ctx = context.WithValue(ctx, "authToken", authToken)
	HandleResponse(ctx)
}

func HandleResponse(ctx context.Context) {
	fmt.Printf(
		"handling response for %v (%v)",
		ctx.Value("userID"),
		ctx.Value("authToken"),
	)
}

type ctxKey int

const (
	ctxUserID ctxKey = iota
	ctxAuthToken
)

func UserID(c context.Context) string {
	return c.Value(ctxUserID).(string)
}

func AuthToken(c context.Context) string {
	return c.Value(ctxAuthToken).(string)
}

func ProcessRequest1(userID, authToken string) {
	ctx := context.WithValue(context.Background(), ctxUserID, userID)
	ctx = context.WithValue(ctx, ctxAuthToken, authToken)
	HandleResponse1(ctx)
}

func HandleResponse1(ctx context.Context) {
	fmt.Printf(
		"handling response for %v (auth: %v)",
		UserID(ctx),
		AuthToken(ctx),
	)
}

