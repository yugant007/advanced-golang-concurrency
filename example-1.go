package main

func main() {
	//var wg sync.WaitGroup
	//x:= "ff"
	//sayHello := func() {
	//	defer wg.Done()
	//	x = "dd"
	//	fmt.Println("hello")
	//}
	//wg.Add(1)
	//go sayHello()
	//wg.Wait()
	//fmt.Println(x)


	//We’ve been using a lot of anonymous functions in our examples to create quick goroutine examples.
	//Let’s shift our attention to closures. Closures close around the lexical scope they are created in,
	//thereby capturing variables. If you run a closure in a goroutine, does the closure operate on a copy
	//of these variables, or the original references? Let’s give it a try and see:
	//var wg sync.WaitGroup
	//salutation := "hello"
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	salutation = "welcome"
	//}()
	//wg.Wait()
	//fmt.Println(salutation)
	//What do you think the value of salutation will be: “hello” or “welcome”? Let’s run it and find out:
	// it will be welcome
	//Interesting! It turns out that goroutines execute within the same address space they were created
	//in, and so our program prints out the word “welcome.”


	//Let’s try another example. What do you think this program will output?
	//var wg sync.WaitGroup
	//for _, salutation := range []string{"hello", "greetings", "good day"} {
	//	wg.Add(1)
	//	go func() {
	//		defer wg.Done()
	//		fmt.Println(salutation) // Here we reference the loop variable salutation created by ranging over a string slice.
	//	}()
	//}
	//wg.Wait()
	// the output will be
	//good day
	//good day
	//good day
	//That’s kind of surprising! Let’s figure out what’s going on here. In this example, the goroutine is
	//running a closure that has closed over the iteration variable salutation, which has a type of string.
	//As our loop iterates, salutation is being assigned to the next string value in the slice literal.
	//Because the goroutines being scheduled may run at any point in time in the future, it is undetermined
	//what values will be printed from within the goroutine. On my machine, there is a high probability the
	//loop will exit before the goroutines are begun. This means the salutation variable falls out of scope.
	//What happens then? Can the goroutines still reference something that has fallen out of scope? Won’t
	//the goroutines be accessing memory that has potentially been garbage collected?


	//Usually on my machine, the loop exits before any goroutines begin running, so salutation is transferred
	//to the heap holding a reference to the last value in my string slice, “good day.” And so I usually see
	//“good day” printed three times. The proper way to write this loop is to pass a copy of salutation into
	//the closure so that by the time the goroutine is run, it will be operating on the data from its iteration
	//of the loop:
	//var wg sync.WaitGroup
	//for _, salutation := range []string{"hello", "greetings", "good day"} {
	//	wg.Add(1)
	//	go func(salutation string) {
	//		defer wg.Done()
	//		fmt.Println(salutation)
	//	}(salutation)
	//}
	//wg.Wait()
	// the output is
	//good day
	//hello
	//greetings
	//Since multiple goroutines can operate against the same address space, we still have to worry about
	//synchronization. As we’ve discussed, we can choose either to synchronize access to the shared memory
	//the goroutines access, or we can use CSP primitives to share memory by communication.


	//Yet another benefit of goroutines is that they’re extraordinarily lightweight. Here’s an excerpt from the Go FAQ:
	//A newly minted goroutine is given a few kilobytes, which is almost always enough. When it isn’t, the
	//run-time grows (and shrinks) the memory for storing the stack automatically, allowing many goroutines
	//to live in a modest amount of memory. The CPU overhead averages about three cheap instructions per function
	//call. It is practical to create hundreds of thousands of goroutines in the same address space. If goroutines
	//were just threads, system resources would run out at a much smaller number.\


	//A few kilobytes per goroutine; that isn’t bad at all! Let’s try and verify that for ourselves. But
	//before we do, we have to cover one interesting thing about goroutines: the garbage collector does
	//nothing to collect goroutines that have been abandoned somehow. If I write the following:
	//go func() {
	//    // <operation that will block forever>
	//}()
	// Do work
	//The goroutine here will hang around until the process exits.


	//In the following example, we combine the fact that goroutines are not garbage collected with
	//the runtime’s ability to introspect upon itself and measure the amount of memory allocated before
	//and after goroutine creation:
	//memConsumed := func() uint64 {
	//	runtime.GC()
	//	var s runtime.MemStats
	//	runtime.ReadMemStats(&s)
	//	return s.Sys
	//}
	//
	//var c <-chan interface{}
	//var wg sync.WaitGroup
	//noop := func() {
	//	wg.Done()
	//	<-c // We require a goroutine that will never exit so that we can keep a number of them in memory for measurement. Don’t worry about how we’re achieving this at this time; just know that this goroutine won’t exit until the process is finished.
	//}
	//
	//const numGoroutines = 1e4 //Here we define the number of goroutines to create. We will use the law of large numbers to asymptotically approach the size of a goroutine.
	//wg.Add(numGoroutines)
	//before := memConsumed() //Here we measure the amount of memory consumed before creating our goroutines.
	//fmt.Println(before)
	//for i := numGoroutines; i > 0; i-- {
	//	go noop()
	//}
	//wg.Wait()
	//after := memConsumed() //And here we measure the amount of memory consumed after creating our goroutines.
	//fmt.Println(after)
	//fmt.Printf("%.3fkb", float64(after-before)/numGoroutines/1000)
	// And here’s the result:
	// 2.817kb
	//It looks like the documentation is correct! These are just empty goroutines that don’t do anything,
	//but it still gives us an idea of the number of goroutines we can likely create. Table
	//gives some rough estimates of how many goroutines you could likely create with a 64-bit CPU without
	//using swap space.
	//Memory (GB)	Goroutines (#/100,000)	Order of magnitude
	// 2^0             3.718                   3
	// 2^1			   7.436                   3
	// 2^2 			   14.873				   6
	// 2^3			   29.746                  6
	// 2^4			   59.492				   6
	// 2^5			   118.983				   6


	//Something that might dampen our spirits is context switching, which is when something hosting a
	//concurrent process must save its state to switch to running a different concurrent process. If we
	//have too many concurrent processes, we can spend all of our CPU time context switching between
	//them and never get any real work done. At the OS level, with threads, this can be quite costly.
	//The OS thread must save things like register values, lookup tables, and memory maps to successfully
	//be able to switch back to the current thread when it is time. Then it has to load the same information
	//for the incoming thread.
	// Context switching in software is comparatively much, much cheaper. Under a software-defined scheduler,
	//the runtime can be more selective in what is persisted for retrieval, how it is persisted, and when
	//the persisting need occur. Let’s take a look at the relative performance of context switching on my
	//laptop between OS threads and goroutines. First, we’ll utilize Linux’s built-in benchmarking suite to
	//measure how long it takes to send a message between two threads on the same core:
	//That gives us 1.467 μs per context switch. That does’t seem too bad, but let’s reserve judgment
	//until we examine context switches between goroutines.
	// run benchmark_1_test.go file.
	// We run the benchmark specifying that we only want to utilize one CPU so that it’s a similar test
	//to the Linux benchmark. Let’s take a look at the results:
	// go test -bench=. -cpu=1 ./benchmark_1_test.go
	// 225 ns per context switch, wow! That’s 0.225 μs, or 92% faster than an OS context switch
	//on my machine, which if you recall took 1.467 μs. It’s difficult to make any claims about
	//how many goroutines will cause too much context switching, but we can comfortably say that
	//the upper limit is likely not to be any kind of barrier to using goroutines.


}
