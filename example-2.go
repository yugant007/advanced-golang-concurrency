package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func main() {
	// The sync package contains the concurrency primitives that are most useful for low-level memory access
	//synchronization.
	//These operations have their use—mostly in small scopes such as a struct. It will be up to you to decide
	//when memory access synchronization is appropriate.
	//Follow example-2.go file

	//WaitGroup ->
	//WaitGroup is a great way to wait for a set of concurrent operations to complete when you either don’t care
	//about the result of the concurrent operation, or you have other means of collecting their results.
	//If neither of those conditions are true, I suggest you use channels and a select statement instead.
	//var wg sync.WaitGroup
	//
	//wg.Add(1) // Here we call Add with an argument of 1 to indicate that one goroutine is beginning.
	//go func() {
	//	defer wg.Done() // Here we call Done using the defer keyword to ensure that before we exit the goroutine’s closure, we indicate to the WaitGroup that we’ve exited.
	//	fmt.Println("1st goroutine sleeping...")
	//	time.Sleep(1)
	//}()
	//
	//wg.Add(1) // Here we call Add with an argument of 1 to indicate that one goroutine is beginning.
	//go func() {
	//	defer wg.Done() // Here we call Done using the defer keyword to ensure that before we exit the goroutine’s closure, we indicate to the WaitGroup that we’ve exited.
	//	fmt.Println("2nd goroutine sleeping...")
	//	time.Sleep(2)
	//}()
	//
	//wg.Wait() // Here we call Wait, which will block the main goroutine until all goroutines have indicated they have exited.
	//fmt.Println("All goroutines complete.")
	//You can think of a WaitGroup like a concurrent-safe counter: calls to Add increment the
	//counter by the integer passed in, and calls to Done decrement the counter by one. Calls
	//to Wait block until the counter is zero.
	//Notice that the calls to Add are done outside the goroutines they’re helping to track.
	//If we didn’t do this, we would have introduced a race condition, because remember from “Goroutines”
	//that we have no guarantees about when the goroutines will be scheduled; we could reach the call to
	//Wait before either of the goroutines begin. Had the calls to Add been placed inside the goroutines’
	//closures, the call to Wait could have returned without blocking at all because the calls to Add would
	//not have taken place.
	//It’s customary to couple calls to Add as closely as possible to the goroutines they’re helping to
	//track, but sometimes you’ll find Add called to track a group of goroutines all at once. I usually do
	//this before for loops like this:
	//hello := func(wg *sync.WaitGroup, id int) {
	//	defer wg.Done()
	//	fmt.Printf("Hello from %v!\n", id)
	//}
	//
	//const numGreeters = 5
	//var wg sync.WaitGroup
	//wg.Add(numGreeters)
	//for i := 0; i < numGreeters; i++ {
	//	go hello(&wg, i+1) // we don't need to pass wg here, but since hello anonymous function call defined first before wg, therefore we are passing wg.
	//}
	//wg.Wait()

	//Mutex and RWMutex
	//Mutex stands for “mutual exclusion” and is a way to guard critical sections of your program.
	//a critical section is an area of your program that requires exclusive access to a shared resource.
	// A Mutex provides a concurrent-safe way to express exclusive access to these shared resources.
	//whereas channels share memory by communicating, a Mutex shares memory by
	//creating a convention developers must follow to synchronize access to the memory. You are
	//responsible for coordinating access to this memory by guarding access to it with a mutex.
	//var count int
	//var lock sync.Mutex
	//
	//increment := func() {
	//	lock.Lock() // Here we request exclusive use of the critical section—in this case the count variable—guarded by a Mutex, lock.
	//	defer lock.Unlock() // Here we indicate that we’re done with the critical section lock is guarding.
	//	count++
	//	fmt.Printf("Incrementing: %d\n", count)
	//}
	//
	//decrement := func() {
	//	lock.Lock() // Here we request exclusive use of the critical section—in this case the count variable—guarded by a Mutex, lock.
	//	defer lock.Unlock() // Here we indicate that we’re done with the critical section lock is guarding.
	//	// if we don't have unlock then this will create a deadlock bcs since we are using wait group that means
	//	// main program will not terminate until all go-routine finishes, other go-routine's can't acquire this lock,
	//	// bcs lock was already acquired by some other go-routine, and hence all go-routine are sleep.
	//	count--
	//	fmt.Printf("Decrementing: %d\n", count)
	//}
	//
	//// Increment
	//var arithmetic sync.WaitGroup
	//for i := 0; i <= 5; i++ {
	//	arithmetic.Add(1)
	//	go func() {
	//		defer arithmetic.Done()
	//		increment()
	//	}()
	//}
	//
	//// Decrement
	//for i := 0; i <= 5; i++ {
	//	arithmetic.Add(1)
	//	go func() {
	//		defer arithmetic.Done()
	//		decrement()
	//	}()
	//}
	//
	//arithmetic.Wait()
	//fmt.Println("Arithmetic complete.")
	// You’ll notice that we always call Unlock within a defer statement. This is a very common idiom when
	//utilizing a Mutex to ensure the call always happens, even when panicing. Failing to do so will probably
	//cause your program to deadlock.
	//Critical sections are so named because they reflect a bottleneck in your program. It is somewhat
	//expensive to enter and exit a critical section, and so generally people attempt to minimize the
	//time spent in critical sections.
	// One strategy for doing so is to reduce the cross-section of the critical section. There may be
	//memory that needs to be shared between multiple concurrent processes, but perhaps not all of
	//these processes will read and write to this memory. If this is the case, you can take advantage
	//of a different type of mutex: sync.RWMutex.
	// The sync.RWMutex is conceptually the same thing as a Mutex: it guards access to memory; however,
	//RWMutex gives you a little bit more control over the memory. You can request a lock for reading,
	//in which case you will be granted access unless the lock is being held for writing. This means
	//that an arbitrary number of readers can hold a reader lock so long as nothing else is holding a
	//writer lock. Here’s an example that demonstrates a producer that is less active than the numerous
	//consumers the code creates:
	//producer := func(wg *sync.WaitGroup, l sync.Locker) { // The producer function’s second parameter is of the type sync.Locker. This interface has two methods, Lock and Unlock, which the Mutex and RWMutex types satisfy.
	//	defer wg.Done()
	//	for i := 5; i > 0; i-- {
	//		l.Lock()
	//		l.Unlock()
	//		time.Sleep(1) //Here we make the producer sleep for one second to make it less active than the observer goroutines.
	//	}
	//}
	//
	//observer := func(wg *sync.WaitGroup, l sync.Locker) {
	//	defer wg.Done()
	//	l.Lock()
	//	defer l.Unlock()
	//}
	//
	//test := func(count int, mutex, rwMutex sync.Locker) time.Duration {
	//	var wg sync.WaitGroup
	//	wg.Add(count+1)
	//	beginTestTime := time.Now()
	//	go producer(&wg, mutex)
	//	for i := count; i > 0; i-- {
	//		go observer(&wg, rwMutex)
	//	}
	//
	//	wg.Wait()
	//	return time.Since(beginTestTime)
	//}
	//
	//tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	//defer tw.Flush()
	//
	//var m sync.RWMutex
	//_, _ = fmt.Fprintf(tw, "Readers\tRWMutext\tMutex\n")
	//for i := 0; i < 20; i++ {
	//	count := int(math.Pow(2, float64(i)))
	//	_, _ = fmt.Fprintf(tw, "%d\t%v\t%v\n", count, test(count, &m, m.RLocker()), test(count, &m, &m))
	//}
	// but it’s usually advisable to use RWMutex instead of Mutex when it logically makes sense.

	//Cond
	// a rendezvous point for goroutines waiting for or announcing the occurrence
	//of an event.
	// In that definition, an “event” is any arbitrary signal between two or more goroutines that carries no
	//information other than the fact that it has occurred. Very often you’ll want to wait for one of these
	//signals before continuing execution on a goroutine. If we were to look at how to accomplish this without
	//the Cond type, one naive approach to doing this is to use an infinite loop:
	//for conditionTrue() == false {
	//}
	// However this would consume all cycles of one core. To fix that, we could introduce a time.Sleep:
	//for conditionTrue() == false {
	//    time.Sleep(1*time.Millisecond)
	//}
	//This is better, but it’s still inefficient, and you have to figure out how long to sleep for: too long,
	//and you’re artificially degrading performance; too short, and you’re unnecessarily consuming too much
	//CPU time. It would be better if there were some kind of way for a goroutine to efficiently sleep until
	//it was signaled to wake and check its condition. This is exactly what the Cond type does for us. Using
	//a Cond, we could write the previous examples like this:
	//c := sync.NewCond(&sync.Mutex{}) // Here we instantiate a new Cond. The NewCond function takes in a type that satisfies the sync.Locker interface. This is what allows the Cond type to facilitate coordination with other goroutines in a concurrent-safe way.
	//c.L.Lock()// Here we lock the Locker for this condition. This is necessary because the call to Wait automatically calls Unlock on the Locker when entered.
	//for conditionTrue() == false {
	//	c.Wait() // Here we wait to be notified that the condition has occurred. This is a blocking call and the goroutine will be suspended.
	//}
	//c.L.Unlock() // Here we unlock the Locker for this condition. This is necessary because when the call to Wait exits, it calls Lock on the Locker for the condition.
	// This approach is much more efficient. Note that the call to Wait doesn’t just block,
	//it suspends the current goroutine, allowing other goroutines to run on the OS thread.
	//A few other things happen when you call Wait: upon entering Wait, Unlock is called on the
	//Cond variable’s Locker, and upon exiting Wait, Lock is called on the Cond variable’s Locker.
	//In my opinion, this takes a little getting used to; it’s effectively a hidden side effect of
	//the method. It looks like we’re holding this lock the entire time while we wait for the
	//condition to occur, but that’s not actually the case. When you’re scanning code, you’ll just
	//have to keep an eye out for this pattern.
	// Let’s expand on this example and show both sides of the equation: a goroutine that is waiting
	//for a signal, and a goroutine that is sending signals. Say we have a queue of fixed length 2,
	//and 10 items we want to push onto the queue. We want to enqueue items as soon as there is room,
	//so we want to be notified as soon as there’s room in the queue. Let’s try using a Cond to manage
	//this coordination:
	//c := sync.NewCond(&sync.Mutex{}) // First, we create our condition using a standard sync.Mutex as the Locker.
	//queue := make([]interface{}, 0, 10) // Next, we create a slice with a length of zero. Since we know we’ll eventually add 10 items, we instantiate it with a capacity of 10.
	//removeFromQueue := func(delay time.Duration) {
	//	time.Sleep(delay)
	//	c.L.Lock() // We once again enter the critical section for the condition so we can modify data pertinent to the condition.
	//	queue = queue[1:] // Here we simulate dequeuing an item by reassigning the head of the slice to the second item.
	//	fmt.Println("Removed from queue")
	//	c.L.Unlock() // Here we exit the condition’s critical section since we’ve successfully dequeued an item.
	//	c.Signal() // Here we let a goroutine waiting on the condition know that something has occurred.
	//}
	//for i := 0; i < 10; i++{
	//	c.L.Lock() // We enter the critical section for the condition by calling Lock on the condition’s Locker.
	//	for len(queue) == 2 { // Here we check the length of the queue in a loop. This is important because a signal on the condition doesn’t necessarily mean what you’ve been waiting for has occurred—only that something has occurred.
	//		c.Wait() // We call Wait, which will suspend the main goroutine until a signal on the condition has been sent.
	//	}
	//	fmt.Println("Adding to queue")
	//	queue = append(queue, struct{}{})
	//	go removeFromQueue(1*time.Second) // Here we create a new goroutine that will dequeue an element after one second.
	//	c.L.Unlock() // Here we exit the condition’s critical section since we’ve successfully enqueued an item.
	//}
	//fmt.Println(queue)
	//As you can see, the program successfully adds all 10 items to the queue (and exits before it has
	//a chance to dequeue the last two items). It also always waits until at least one item is dequeued
	//before enqueing another.
	// We also have a new method in this example, Signal. This is one of two methods that the Cond
	//type provides for notifying goroutines blocked on a Wait call that the condition has been
	//triggered. The other is a method called Broadcast. Internally, the runtime maintains a FIFO
	//list of goroutines waiting to be signaled; Signal finds the goroutine that’s been waiting the
	//longest and notifies that, whereas Broadcast sends a signal to all goroutines that are waiting.
	//Broadcast is arguably the more interesting of the two methods as it provides a way to communicate
	//with multiple goroutines at once. We can trivially reproduce Signal with channels (as we’ll see
	//in the section “Channels”), but reproducing the behavior of repeated calls to Broadcast would be
	//more difficult. In addition, the Cond type is much more performant than utilizing channels.
	// To get a feel for what it’s like to use Broadcast, let’s imagine we’re creating a GUI
	//application with a button on it. We want to register an arbitrary number of functions
	//that will run when that button is clicked. A Cond is perfect for this because we can use
	//its Broadcast method to notify all registered handlers. Let’s see how that might look:
	//type Button struct { // We define a type Button that contains a condition, Clicked.
	//	Clicked *sync.Cond
	//}
	//button := Button{ Clicked: sync.NewCond(&sync.Mutex{}) }
	//
	//subscribe := func(c *sync.Cond, fn func()) { //Here we define a convenience function that will allow us to register functions to handle signals from a condition. Each handler is run on its own goroutine, and subscribe will not exit until that goroutine is confirmed to be running.
	//	var goroutineRunning sync.WaitGroup
	//	goroutineRunning.Add(1)
	//	go func() {
	//		goroutineRunning.Done()
	//		c.L.Lock()
	//		defer c.L.Unlock()
	//		c.Wait()
	//		fn()
	//	}()
	//	goroutineRunning.Wait()
	//}
	//
	//var clickRegistered sync.WaitGroup
	//clickRegistered.Add(3) //Here we set a handler for when the mouse button is raised. It in turn calls Broadcast on the Clicked Cond to let all handlers know that the mouse button has been clicked (a more robust implementation would first check that it had been depressed).
	//subscribe(button.Clicked, func() {
	//	fmt.Println("Maximizing window.")
	//	clickRegistered.Done()
	//})
	//subscribe(button.Clicked, func() {
	//	fmt.Println("Displaying annoying dialog box!")
	//	clickRegistered.Done()
	//})
	//subscribe(button.Clicked, func() {
	//	fmt.Println("Mouse clicked.")
	//	clickRegistered.Done()
	//})
	//
	//button.Clicked.Broadcast() //Next, we simulate a user raising the mouse button from having clicked the application’s button.
	//
	//clickRegistered.Wait()
	//You can see that with one call to Broadcast on the Clicked Cond, all three handlers are run.
	//Were it not for the clickRegistered WaitGroup, we could call button.Clicked.Broadcast()
	//multiple times, and each time all three handlers would be invoked. This is something channels
	//can’t do easily and thus is one of the main reasons to utilize the Cond type.
	// Like most other things in the sync package, usage of Cond works best when constrained to a
	//tight scope, or exposed to a broader scope through a type that encapsulates it.

	//Once
	//What do you think this code will print out?
	//var count int
	//
	//increment := func() {
	//	count++
	//}
	//
	//var once sync.Once
	//
	//var increments sync.WaitGroup
	//increments.Add(100)
	//for i := 0; i < 100; i++ {
	//	go func() {
	//		defer increments.Done()
	//		once.Do(increment)
	//	}()
	//}
	//
	//increments.Wait()
	//fmt.Printf("Count is %d\n", count)
	//It’s tempting to say the result will be Count is 100, but I’m sure you’ve noticed the sync.Once variable,
	//and that we’re somehow wrapping the call to increment within the Do method of once. In fact, this code
	//will print out the following:
	//Count is 1
	//As the name implies, sync.Once is a type that utilizes some sync primitives internally to ensure that only
	//one call to Do ever calls the function passed in—even on different goroutines. This is indeed because we
	//wrap the call to increment in a sync.Once Do method.
	//It may seem like the ability to call a function exactly once is a strange thing to encapsulate and put into
	//the standard package, but it turns out that the need for this pattern comes up rather frequently. Just for
	//fun, let’s check Go’s standard library and see how often Go itself uses this primitive. Here’s a grep
	//command that will perform the search:
	//grep -ir sync.Once $(go env GOROOT)/src |wc -l
	// this produces 112
	//There are a few things to note about utilizing sync.Once. Let’s take a look at another example;
	//what do you think it will print?
	//var count int
	//increment := func() { count++ }
	//decrement := func() { count-- }
	//
	//var once sync.Once
	//once.Do(increment)
	//once.Do(decrement)
	//
	//fmt.Printf("Count: %d\n", count)
	//Is it surprising that the output displays 1 and not 0? This is because sync.Once only counts the number
	//of times Do is called, not how many times unique functions passed into Do are called. In this way,
	//copies of sync.Once are tightly coupled to the functions they are intended to be called with; once
	//again we see how usage of the types within the sync package work best within a tight scope. I
	//recommend that you formalize this coupling by wrapping any usage of sync.Once in a small lexical
	//block: either a small function, or by wrapping both in a type. What about this example? What do
	//you think will happen?
	//var onceA, onceB sync.Once
	//var initB func()
	//initA := func() { onceB.Do(initB) }
	//initB = func() { onceA.Do(initA) }
	//onceA.Do(initA)
	//This program will deadlock because the call to Do at 1 won’t proceed until the call to Do at 2 exits—a
	//classic example of a deadlock. For some, this may be slightly counterintuitive since it appears
	//as though we’re using sync.Once as intended to guard against multiple initialization, but the
	//only thing sync.Once guarantees is that your functions are only called once.

	//Pool
	//At a high level, a the pool pattern is a way to create and make available a fixed number, or pool, of
	//things for use. It’s commonly used to constrain the creation of things that are expensive (e.g.,
	//database connections) so that only a fixed number of them are ever created, but an indeterminate
	//number of operations can still request access to these things. In the case of Go’s sync.Pool,
	//this data type can be safely used by multiple goroutines.
	//Pool’s primary interface is its Get method. When called, Get will first check whether there are any available
	//instances within the pool to return to the caller, and if not, call its New member variable to create a
	//new one. When finished, callers call Put to place the instance they were working with back in the pool
	//for use by other processes. Here’s a simple example to demonstrate:
	//myPool := &sync.Pool{
	//	New: func() interface{} {
	//		fmt.Println("Creating new instance.")
	//		return struct{}{}
	//	},
	//}
	//
	//myPool.Get() // Here we call Get on the pool. These calls will invoke the New function defined on the pool since instances haven’t yet been instantiated.
	//instance := myPool.Get() // Here we call Get on the pool. These calls will invoke the New function defined on the pool since instances haven’t yet been instantiated.
	//myPool.Put(instance) // Here we put an instance previously retrieved back in the pool. This increases the available number of instances to one.
	//myPool.Get() // When this call is executed, we will reuse the instance previously allocated and put it back in the pool. The New function will not be invoked.
	// So why use a pool and not just instantiate objects as you go? Go has a garbage collector,
	//so the instantiated objects will be automatically cleaned up. What’s the point? Consider this example:
	//var numCalcsCreated int
	//calcPool := &sync.Pool {
	//	New: func() interface{} {
	//		numCalcsCreated += 1
	//		mem := make([]byte, 1024)
	//		return &mem // Notice that we are storing the address of the slice of bytes. bcs passing and storing a pointer variable is faster.
	//	},
	//}
	//
	//// Seed the pool with 4KB
	//calcPool.Put(calcPool.New())
	//calcPool.Put(calcPool.New())
	//calcPool.Put(calcPool.New())
	//calcPool.Put(calcPool.New())
	//
	//const numWorkers = 1024*1024
	//var wg sync.WaitGroup
	//wg.Add(numWorkers)
	//for i := numWorkers; i > 0; i-- {
	//	go func() {
	//		defer wg.Done()
	//
	//		mem := calcPool.Get().(*[]byte) // And here we are asserting the type is a pointer to a slice of bytes.
	//		defer calcPool.Put(mem)
	//
	//		// Assume something interesting, but quick is being done with
	//		// this memory.
	//	}()
	//}
	//wg.Wait()
	//fmt.Printf("%d calculators were created.", numCalcsCreated)
	//This produces:
	//4 calculators were created.
	// Had I run this example without a sync.Pool, though the results are non-deterministic, in the worst case
	//I could have been attempting to allocate a gigabyte of memory.
	//Another common situation where a Pool is useful is for warming a cache of pre-allocated objects for
	//operations that must run as quickly as possible. In this case, instead of trying to guard the host
	//machine’s memory by constraining the number of objects created, we’re trying to guard consumers’
	//time by front-loading the time it takes to get a reference to another object. This is very common
	//when writing high-throughput network servers that attempt to respond to requests as quickly as
	//possible. Let’s take a look at such a scenario.
	//First, let’s create a function that simulates creating a connection to a service.
	//We’ll make this connection take a long time:
	// connectToService()
	// Next, let’s see how performant a network service would be if for every request we started a new
	//connection to the service. We’ll write a network handler that opens a connection to another
	//service for every connection the network handler accepts. To make the benchmarking simple,
	//we’ll only allow one connection at a time:
	// startNetworkDaemon1()
	// Now let’s benchmark this:
	// see benchmark_2_test.go
	// cmd  -> go test -benchtime=10s -bench=. ./benchmark_2_test.go example-2.go
	//This produces:
	//BenchmarkNetworkRequest1-8 	10 		1000385643ns/op
	//Looks like like roughly 1E9 ns/op. This seems reasonable as far as performance goes,
	//but let’s see if we can improve it by using a sync.Pool to host connections to our fictitious service:
	// warmServiceConnCache(), startNetworkDaemon2().
	// Now let’s benchmark this:
	// see benchmark_3_test.go
	// cmd  -> go test -benchtime=10s -bench=. ./benchmark_3_test.go example-2.go
	//This produces:
	//BenchmarkNetworkRequest2-4          9494           2374108 ns/op
	//2.3E6 ns/op: three orders of magnitude faster! You can see how utilizing this pattern when
	//working with things that are expensive to create can drastically improve response time.
	// As we’ve seen, the object pool design pattern is best used either when you have concurrent
	//processes that require objects, but dispose of them very rapidly after instantiation,
	//or when construction of these objects could negatively impact memory.
	//However, there is one thing to be wary of when determining whether or not you should utilize
	//a Pool: if the code that utilizes the Pool requires things that are not roughly homogenous,
	//you may spend more time converting what you’ve retrieved from the Pool than it would have
	//taken to just instantiate it in the first place. For instance, if your program requires
	//slices of random and variable length, a Pool isn’t going to help you much. The probability
	//that you’ll receive a slice the length you require is low.
	//So when working with a Pool, just remember the following points:
	//1. When instantiating sync.Pool, give it a New member variable that is thread-safe when called.
	//2. When you receive an instance from Get, make no assumptions regarding the state of the object you receive back.
	//Make sure to call Put when you’re finished with the object you pulled out of the pool. Otherwise, the
	//3. Pool is useless. Usually this is done with defer.
	//4. Objects in the pool must be roughly uniform in makeup.

	//Channels
	//Like a river, a channel serves as a conduit for a stream of information; values may be passed along the channel,
	//and then read out downstream. For this reason I usually end my chan variable names with the word “Stream.”
	//When using channels, you’ll pass a value into a chan variable, and then somewhere else in your program
	//read it off the channel. The disparate parts of your program don’t require knowledge of each other,
	//only a reference to the same place in memory where the channel resides. This can be done by passing
	//references of channels around your program.
	// creating a channel
	//var dataStream chan interface{} // Here we declare a channel. We say it is “of type” interface{} since the type we’ve declared is the empty interface.
	//dataStream = make(chan interface{}) // Here we instantiate the channel using the built-in make function.
	//To declare a unidirectional channel, you’ll simply include the <- operator. To both declare and instantiate
	//a channel that can only read, place the <- operator on the lefthand side, like so:
	//var dataStream <-chan interface{}
	//dataStream := make(<-chan interface{})
	//And to declare and create a channel that can only send, you place the <- operator on the righthand side, like so:
	//var dataStream chan<- interface{}
	//dataStream := make(chan<- interface{})
	//This is possible because Go will implicitly convert bidirectional channels to unidirectional
	//channels when needed. Here’s an example:
	//var receiveChan <-chan interface{}
	//var sendChan chan<- interface{}
	//dataStream := make(chan interface{})
	//
	//// Valid statements:
	//receiveChan = dataStream
	//sendChan = dataStream
	// Keep in mind channels are typed. In this example, we created a chan interface{} variable, which
	//means that we can place any kind of data onto it, but we can also give it a stricter type to
	//constrain the type of data it could pass along. Here’s an example of a channel for integers.
	//intStream := make(chan int)
	//To use channels, we’ll once again make use of the <- operator. Sending is done by placing the
	//<- operator to the right of a channel, and receiving is done by placing the <- operator to the left of the channel
	//stringStream := make(chan string)
	//go func() {
	//	stringStream <- "Hello channels!" //Here we pass a string literal onto the channel stringStream.
	//}()
	//fmt.Println(<-stringStream) // Here we read the string literal off of the channel and print it out to stdout.
	//Pretty simple, right? All you need is a channel variable and you can pass data onto it and read
	//data off of it; however, it is an error to try and write a value onto a read-only channel, and
	//an error to read a value from a write-only channel. If we try and compile the following example,
	//Go’s compiler will let us know that we’re doing something illegal:
	//writeStream := make(chan<- interface{})
	//readStream := make(<-chan interface{})
	//
	//<-writeStream
	//readStream <- struct{}{}
	// channels in Go are said to be blocking.
	// This example works because channels in Go are said to be blocking. This means that any goroutine
	//that attempts to write to a channel that is full will wait until the channel has been emptied, and any goroutine
	//that attempts to read from a channel that is empty will wait until at least one item is placed on it.
	// This can cause deadlocks if you don’t structure your program correctly. Take a look at the following example,
	//which introduces a nonsensical conditional to prevent the anonymous goroutine from placing a value on the channel:
	//stringStream := make(chan string)
	//go func() {
	//	if 0 != 1 {
	//		return
	//	}
	//	stringStream <- "Hello channels!"
	//}()
	//fmt.Println(<-stringStream)
	// The receiving form of the <- operator can also optionally return two values, like so:
	//stringStream := make(chan string)
	//go func() {
	//	stringStream <- "Hello channels!"
	//}()
	//salutation, ok := <-stringStream
	//fmt.Printf("(%v): %v", ok, salutation)
	// Very curious! What does the boolean signify? The second return value is a way for a read operation to
	//indicate whether the read off the channel was a value generated by a write elsewhere in the process, or
	//a default value generated from a closed channel. Wait a second; a closed channel, what’s that?
	// In programs, it’s very useful to be able to indicate that no more values will be sent over a channel. This helps
	//downstream processes know when to move on, exit, re-open communications on a new or different channel, etc.
	// To close a channel, we use the close keyword, like so:
	//valueStream := make(chan interface{})
	//close(valueStream)
	//Interestingly, we can read from a closed channel as well. Take this example:
	//intStream := make(chan int)
	//close(intStream)
	//integer, ok := <- intStream
	//fmt.Printf("(%v): %v", ok, integer)
	// Notice that we never placed anything on this channel; we closed it immediately. We were still able to
	//perform a read operation, and in fact, we could continue performing reads on this channel indefinitely
	//despite the channel remaining closed. This is to allow support for multiple downstream reads from a
	//single upstream writer on the channel. The second value returned—here stored in the ok variable—is false,
	//indicating that the value we received is the zero value for int, or 0, and not a value placed on the stream.
	// This opens up a few new patterns for us. The first is ranging over a channel. The range keyword—used in
	//conjunction with the for statement—supports channels as arguments, and will automatically break the loop
	//when a channel is closed. This allows for concise iteration over the values on a channel. Let’s take a
	//look at an example:
	//intStream := make(chan int)
	//go func() {
	//	defer close(intStream)
	//	for i := 1; i <= 5; i++ {
	//		intStream <- i
	//	}
	//}()
	//
	//for integer := range intStream {
	//	fmt.Printf("%v ", integer)
	//}
	// Notice how the loop doesn’t need an exit criteria, and the range does not return the second boolean value.
	//The specifics of handling a closed channel are managed for you to keep the loop concise.
	// Closing a channel is also one of the ways you can signal multiple goroutines simultaneously. If you
	//have n goroutines waiting on a single channel, instead of writing n times to the channel to unblock each
	//goroutine, you can simply close the channel. Since a closed channel can be read from an infinite number
	//of times, it doesn’t matter how many goroutines are waiting on it, and closing the channel is both cheaper
	//and faster than performing n writes. Here’s an example of unblocking multiple goroutines at once:
	//begin := make(chan interface{})
	//var wg sync.WaitGroup
	//for i := 0; i < 5; i++ {
	//	wg.Add(1)
	//	go func(i int) {
	//		defer wg.Done()
	//		<-begin
	//		fmt.Printf("%v has begun\n", i)
	//	}(i)
	//}
	//
	//fmt.Println("Unblocking goroutines...")
	//close(begin)
	//wg.Wait()
	// Remember in “The sync Package” we discussed using the sync.Cond type to perform the same behavior.
	//You can certainly use that, but as we’ve discussed, channels are composable, so this is my favorite
	//way to unblock multiple goroutines at the same time.
	// We can also create buffered channels, which are channels that are given a capacity when they’re instantiated.
	//This means that even if no reads are performed on the channel, a goroutine can still perform n writes,
	//where n is the capacity of the buffered channel. Here’s how to declare and instantiate one:
	//var dataStream chan interface{}
	//dataStream = make(chan interface{}, 4) // Here we create a buffered channel with a capacity of four. This means that we can place four things onto the channel regardless of whether it’s being read from.
	// Once again, I’ve exploded out the instantiation into two lines so you can see that the declaration of a
	//buffered channel is no different than an unbuffered one. This is somewhat interesting because it means
	//that the goroutine that instantiates a channel controls whether it’s buffered. This suggests that the
	//creation of a channel should probably be tightly coupled to goroutines that will be performing writes
	//on it so that we can reason about its behavior and performance more easily.
	// Unbuffered channels are also defined in terms of buffered channels: an unbuffered channel is simply a buffered
	//channel created with a capacity of 0. Here’s an example of two channels that have equivalent functionality:
	//a := make(chan int)
	//b := make(chan int, 0)
	// Both channels are int channels with a capacity of zero. Remember that when we discussed blocking, we said
	//that writes to a channel block if a channel is full, and reads from a channel block if the channel is empty?
	//“Full” and “empty” are functions of the capacity, or buffer size. An unbuffered channel has a capacity of
	//zero and so it’s already full before any writes. A buffered channel with no receivers and a capacity of
	//four would be full after four writes, and block on the fifth write since it has nowhere else to place
	//the fifth element. Like unbuffered channels, buffered channels are still blocking; the preconditions
	//that the channel be empty or full are just different. In this way, buffered channels are an in-memory
	//FIFO queue for concurrent processes to communicate over.
	// An unbuffered channel is a channel that needs a receiver as soon as a message is emitted to the channel.
	// If the capacity is zero or absent, the channel is unbuffered and communication succeeds only when both a
	//sender and receiver are ready.
	// If the channel is unbuffered, the sender blocks until the receiver has received the value
	// Buffered channels can be useful in certain situations, but you should create them with care.
	// buffered channels can easily become a premature optimization and also hide deadlocks by making them more
	//unlikely to happen. This sounds like a good thing, but I’m guessing you’d much rather find a deadlock while
	//writing code the first time, and not in the middle of the night when your production system goes down.
	// Let’s examine another, more complete code example that uses buffered channels
	//var stdoutBuff bytes.Buffer //Here we create an in-memory buffer to help mitigate the nondeterministic nature of the output. It doesn’t give us any guarantees, but it’s a little faster than writing to stdout directly.
	//defer stdoutBuff.WriteTo(os.Stdout) // Here we ensure that the buffer is written out to stdout before the process exits.
	//
	//intStream := make(chan int, 4) // Here we create a buffered channel with a capacity of four.
	//go func() {
	//	defer close(intStream)
	//	defer fmt.Fprintln(&stdoutBuff, "Producer Done.")
	//	for i := 0; i < 5; i++ {
	//		fmt.Fprintf(&stdoutBuff, "Sending: %d\n", i)
	//		intStream <- i
	//	}
	//}()
	//
	//for integer := range intStream {
	//	fmt.Fprintf(&stdoutBuff, "Received %v.\n", integer)
	//}
	// In this example, the order in which output to stdout is written is nondeterministic, but you can
	//still get a rough idea of how the anonymous goroutine is working. If you look at the output, you can see
	//how our anonymous goroutine is able to place all five of its results on the intStream and exit before the
	//main goroutine pulls even one result off:
	//Sending: 0
	//Sending: 1
	//Sending: 2
	//Sending: 3
	//Sending: 4
	//Producer Done.
	//Received 0.
	//Received 1.
	//Received 2.
	//Received 3.
	//Received 4.
	//This is an example of an optimization that can be useful under the right conditions: if a goroutine making
	//writes to a channel has knowledge of how many writes it will make, it can be useful to create a buffered
	//channel whose capacity is the number of writes to be made, and then make those writes as quickly as possible.
	//There are, of course, caveats.
	//  The only aspect of channels we haven’t covered is the default value for channels: nil. How do programs
	// interact with a nil channel? First, let’s try reading from a nil channel:
	//var dataStream chan interface{}
	//<-dataStream
	// A deadlock! This indicates that reading from a nil channel will block (although not necessarily deadlock)
	//a program. What about writes?
	//var dataStream chan interface{}
	//dataStream <- struct{}{}
	// It looks like writes to a nil channel will also block. That just leaves one operation, close. What
	//happens if we attempt to close a nil channel?
	//var dataStream chan interface{}
	//close(dataStream)
	//panic: close of nil channel
	//goroutine 1 [running]:
	// Yipes! This is probably the worst outcome of all the operations performed on a nil channel: a panic.
	//Be sure to ensure the channels you’re working with are always initialized first.
	// Result of channel operations given a channel’s state
	//Operation		 Channel-state				Result
	//Read 			 nil 	            		Block
	//				 Open and Empty 			Block
	//				 Open and Not Empty 		Value
	//				 Closed 					<default value>, false
	//				 Write Only 				Compilation Error
	//write			 nil 						Block
	// 				 Open and Full 				Block
	// 				 Open and Not Full 			Write Value
	//				 Closed						panic
	// 				 Receive Only 				Compilation Error
	//close			 nil						panic
	//				 Open and Not Empty 		Closes Channel; reads succeed until channel is drained,
	//               Open and Empty 			Closes Channel; reads produces default value
	// 				 Closed						panic
	// 				 Receive Only 				Compilation Error
	// For now, let’s look at an example to help clarify these concepts. Let’s create a goroutine that
	//clearly owns a channel, and a consumer that clearly handles blocking and closing of a channel:
	//chanOwner := func() <-chan int {
	//	resultStream := make(chan int, 5) // Here we instantiate a buffered channel. Since we know we’ll produce six results, we create a buffered channel of five so that the goroutine can complete as quickly as possible.
	//	go func() { // Here we start an anonymous goroutine that performs writes on resultStream. Notice that we’ve inverted how we create goroutines. It is now encapsulated within the surrounding function.
	//		defer close(resultStream) // Here we ensure resultStream is closed once we’re finished with it. As the channel owner, this is our responsibility.
	//		for i := 0; i <= 5; i++ {
	//			resultStream <- i
	//		}
	//	}()
	//	return resultStream // Here we return the channel. Since the return value is declared as a read-only channel, resultStream will implicitly be converted to read-only for consumers.
	//}
	//
	//resultStream := chanOwner()
	//for result := range resultStream { //Here we range over resultStream. As a consumer, we are only concerned with blocking and closed channels.
	//	fmt.Printf("Received: %d\n", result)
	//}
	//fmt.Println("Done receiving!")
	// Notice how the lifecycle of the resultStream channel is encapsulated within the chanOwner function.
	//It’s very clear that the writes will not happen on a nil or closed channel, and that the close will always
	//happen once. This removes a large swath of risk from our program. I highly encourage you to do what you can
	//in your programs to keep the scope of channel ownership small so that these things remain obvious. If you
	//have a channel as a member variable of a struct with numerous methods on it, it’s going to quickly become
	//unclear how the channel will behave.
	// The consumer function only has access to a read channel, and therefore only needs to know how it should
	//handle blocking reads and channel closes. In this small example, we’ve taken the stance that it’s perfectly
	//OK to block the life of the program until the channel is closed. As range will not read default value

	//select
	//The select statement is the glue that binds channels together; it’s how we’re able to compose channels
	//together in a program to form larger abstractions. If channels are the glue that binds goroutines together,
	// It is not an overstatement to say that select statements are one of the most crucial things in a
	//Go program with concurrency. You can find select statements binding together channels locally, within
	//a single function or type, and also globally, at the intersection of two or more components in a system.
	//In addition to joining components, at these critical junctures in your program, select statements can
	//help safely bring channels together with concepts like cancellations, timeouts, waiting, and default values.
	// So what are these powerful select statements? How do we use them, and how do they work? Let’s start
	//by just laying one out. Here’s a very simple example:
	//var c1, c2 <-chan interface{}
	//var c3 chan<- interface{}
	//select {
	//case <- c1:
	//	// Do something
	//case <- c2:
	//	// Do something
	//case c3<- struct{}{}:
	//	// Do something
	//}
	//It looks a bit like a switch block, doesn’t it? Just like a switch block, a select block encompasses a series
	//of case statements that guard a series of statements; however, that’s where the similarities end. Unlike
	//switch blocks, case statements in a select block aren’t tested sequentially, and execution won’t automatically
	//fall through if none of the criteria are met.
	//
	//Instead, all channel reads and writes are considered simultaneously to see if any of them are ready:
	//populated or closed channels in the case of reads, and channels that are not at capacity in the case of
	//writes. If none of the channels are ready, the entire select statement blocks. Then when one the channels
	//is ready, that operation will proceed, and its corresponding statements will execute. Let’s take a look
	//at a quick example:
	//start := time.Now()
	//c := make(chan interface{})
	//go func() {
	//	time.Sleep(5*time.Second)
	//	close(c)
	//}()
	//
	//fmt.Println("Blocking on read...")
	//select {
	//case <-c:
	//	fmt.Printf("Unblocked %v later.\n", time.Since(start))
	//}
	//This produces:
	//Blocking on read...
	//Unblocked 5.000170047s later.
	// The first question of multiple channels being ready simultaneously seems interesting.
	//Let’s just try it and see what happens!
	//c1 := make(chan interface{}); close(c1)
	//c2 := make(chan interface{}); close(c2)
	//
	//var c1Count, c2Count int
	//for i := 1000; i >= 0; i-- {
	//	select {
	//	case <-c1:
	//		c1Count++
	//	case <-c2:
	//		c2Count++
	//	}
	//}
	//
	//fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)
	// This produces:
	//c1Count: 505
	//c2Count: 496
	// what happens if there are never any channels that become ready? If there’s nothing useful you can do
	//when all the channels are blocked, but you also can’t block forever, you may want to time out. Go’s time
	//package provides an elegant way to do this with channels that fits nicely within the paradigm of select
	//statements. Here’s an example using one:
	//var c <-chan int
	//select {
	//case <-c:
	//case <-time.After(1 * time.Second):
	//	fmt.Println("Timed out.")
	//}
	//what happens when no channel is ready, and we need to do something in the meantime? Like
	//case statements, the select statement also allows for a default clause in case you’d like to
	//do something if all the channels you’re selecting against are blocking. Here’s an example:
	//start := time.Now()
	//var c1, c2 <-chan int
	//select {
	//case <-c1:
	//case <-c2:
	//default:
	//	fmt.Printf("In default after %v\n\n", time.Since(start))
	//}
	//You can see that it ran the default statement almost instantaneously. This allows you to exit a select
	//block without blocking. Usually you’ll see a default clause used in conjunction with a for-select loop.
	//This allows a goroutine to make progress on work while waiting for another goroutine to report a result.
	//Here’s an example of that:
	//done := make(chan interface{})
	//go func() {
	//	time.Sleep(5*time.Second)
	//	close(done)
	//}()
	//
	//workCounter := 0
	//loop:
	//for {
	//	select {
	//	case <-done:
	//		break loop
	//	default:
	//	}
	//
	//	// Simulate work
	//	workCounter++
	//	time.Sleep(1*time.Second)
	//}
	//
	//fmt.Printf("Achieved %v cycles of work before signalled to stop.\n", workCounter)
	// Finally, there is a special case for empty select statements: select statements with no case clauses.
	//These look like this:
	// select{}
	// This statement will simply block forever.

	//The GOMAXPROCS Lever
	// In the runtime package, there is a function called GOMAXPROCS. In my opinion, the name is misleading:
	//people often think this function relates to the number of logical processors on the host machine—and in
	//a roundabout way it does—but really this function controls the number of OS threads that will host so-called
	//“work queues.
	// Prior to Go 1.5, GOMAXPROCS was always set to one, and usually you’d find this snippet in most Go programs:
	// runtime.GOMAXPROCS(runtime.NumCPU())
	// Almost universally, developers want to take advantage of all the cores on the machine their process is
	//running in. Because of this, in subsequent Go versions, it is now automatically set to the number of
	//logical CPUs on the host machine.
	//So why would you want to tweak this value? Most of the time you won’t want to. Go’s scheduling algorithm
	//is good enough in most situations that increasing or decreasing the number of worker queues and threads
	//will likely do more harm than good, but there are still some situations where changing this value might be useful.
	// For instance, I worked on one project that had a test suite plagued by race conditions.
	//However it came to be, the team had a handful of packages that had tests that sometimes failed.
	//The infrastructure on which we ran our tests only had four logical CPUs, and so at any one point
	//in time we had four goroutines executing simultaneously. By increasing GOMAXPROCS beyond the number
	//of logical CPUs we had, we were able to trigger the race conditions much more often, and thus get
	//them corrected faster.
	//  Tweaking this value pushes your program closer to the metal it’s running on, but at the cost
	// of abstraction and long-term performance stability.
}
func connectToService() interface{} {
	time.Sleep(1 * time.Second)
	return struct{}{}
}

func startNetworkDaemon1() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		server, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer server.Close()

		wg.Done()

		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("cannot accept connection: %v", err)
				continue
			}
			connectToService()
			_, _ = fmt.Fprintln(conn, "")
			_ = conn.Close()
		}
	}()
	return &wg
}

func warmServiceConnCache() *sync.Pool {
	p := &sync.Pool {
		New: connectToService,
	}
	for i := 0; i < 10; i++ {
		p.Put(p.New())
	}
	return p
}


func startNetworkDaemon2() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		connPool := warmServiceConnCache()

		server, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer server.Close()

		wg.Done()

		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("cannot accept connection: %v", err)
				continue
			}
			svcConn := connPool.Get()
			_, _ = fmt.Fprintln(conn, "")
			connPool.Put(svcConn)
			_ = conn.Close()
		}
	}()
	return &wg
}
