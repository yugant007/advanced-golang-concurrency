package main

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
	"sort"
	"time"
)

func main() {
	//Error Propagation
	//Errors indicate that your system has entered a state in which it cannot fulfill an operation that a user
	//either explicitly or implicitly requested. Because of this,it needs to relay a few pieces of critical information:
	//What happened. ->
	//This is the part of the error that contains information about what happened, e.g., “disk full,”
	//“socket closed,” or “credentials expired.” This information is likely to be generated implicitly
	//by whatever it was that generated the errors, although you can probably decorate this with some
	//context that will help the user.
	//When and where it occurred. ->
	//Errors should always contain a complete stack trace starting with how the call was initiated and ending with
	//where the error was instantiated. The stack trace should not be contained in the error message
	//but should be easily accessible when handling the error up the stack.
	//Further, the error should contain information regarding the context it’s running within. For example,
	//in a distributed system, it should have some way of identifying what machine the error occurred on.
	//Later, when trying to understand what happened in your system, this information will be invaluable.
	//In addition, the error should contain the time on the machine the error was instantiated on, in UTC.
	//A friendly user-facing message. ->
	//The message that gets displayed to the user should be customized to suit your system and its users.
	//It should only contain abbreviated and relevant information from the previous two points. A friendly
	//message is human-centric, gives some indication of whether the issue is transitory, and should be
	//about one line of text.
	//How the user can get more information.->
	//At some point, someone will likely want to know, in detail, what happened when the error occurred.
	//Errors that are presented to users should provide an ID that can be cross-referenced to a corresponding
	//log that displays the full information of the error: time the error occurred (not the time the error was logged),
	//the stack trace—everything you stuffed into the error when it was created. It can also be helpful to include
	//a hash of the stack trace to aid in aggregating like issues in bug trackers.
	// It’s possible to place all errors into one of two categories:
	//1. Bugs
	//2. Known edge cases (e.g., broken network connections, failed disk writes, etc.)
	//Imagine a large system with multiple modules:
	// cli component -> intermediary component -> low level component
	//Let’s say an error occurs in the “Low Level Component” and we’ve crafted a well-formed error there to be passed
	//up the stack. Within the context of the “Low Level Component,” this error might be considered well-formed,
	//but within the context of our system, it may not be. Let’s take the stance that at the boundaries of each
	//component, all incoming errors must be wrapped in a well-formed error for the component our code is within.
	//For example, if we were in “Intermediary Component,” and we were calling code from “Low Level Component,”
	//which might error, we could have this:
	//func PostReport(id string) error {
	//	result, err := lowlevel.DoWork()
	//	if err != nil {
	//	if _, ok := err.(lowlevel.Error); ok { // Here we check to ensure we’re receiving a well-formed error. If we aren’t, we’ll simply ferry the malformed error up the stack to indicate a bug.
	//	err = WrapErr(err, "cannot post report with id %q", id) // Here we use a hypothetical function call to wrap the incoming error with pertinent information for our module, and to give it a new type. Note that wrapping the error might involve hiding some low-level details that may not be important for the user within this context.
	//}
	//	return err
	//}
	//	// ...
	//}
	//The low-level details of where the root of the error occurred (e.g., what goroutine, machine, stack trace,
	//etc.) are still filled in when the error is initially instantiated, but our architecture dictates that at
	//module boundaries we convert the error to our module’s error type—potentially filling in pertinent
	//information. Now, any error that escapes our module without our module’s error type can be considered
	//malformed, and a bug. Note that it is only necessary to wrap errors in this fashion at your own module
	//boundaries—public functions/methods—or when your code can add valuable context. Usually this prevents
	//the need for wrapping errors in most of the code.

	//Let’s take a look at a complete example. This example won’t be extremely robust (e.g., the error type is
	//perhaps simplistic), and the call stack is linear, which obfuscates the fact that it’s only necessary
	//to wrap errors at module boundaries.
	// look at MyError{}, wrapError(), LowLevelErr{}, isGloballyExec()

	//Then, let’s create another module, intermediate, which calls functions from the lowlevel package:
	// look at IntermediateErr{}, runJob()

	//Finally, let’s create a top-level main function that calls functions from the intermediate package.
	//This is the user-facing portion of our program:
	//log.SetOutput(os.Stdout)
	//log.SetFlags(log.Ltime | log.LUTC)
	//
	//err := runJob("1")
	//if err != nil {
	//	msg := "There was an unexpected issue; please report this as a bug."
	//	if _, ok := err.(IntermediateErr); ok { // Here we check to see if the error is of the expected type. If it is, we know it’s a well-crafted error, and we can simply pass its message on to the user.
	//		msg = err.Error()
	//	}
	//	handleError(1, err, msg) // On this line we bind the log and error message together with an ID of 1. We could easily make this increase monotonically, or use a GUID to ensure a unique ID.
	//}

	//When we run this, we get a log message that contains:
	//[logID: 1]: 22:11:04 main.IntermediateErr{error:main.MyError
	//  {Inner:main.LowLevelErr{error:main.MyError{Inner:(*os.PathError)
	//  (0xc4200123f0), Message:"stat /bad/job/binary: no such file or directory",
	//  StackTrace:"goroutine 1 [running]:
	//  runtime/debug.Stack(0xc420012420, 0x2f, 0x0)
	//      /home/kate/.guix-profile/src/runtime/debug/stack.go:24 +0x79
	//  main.wrapError(0x530200, 0xc4200123f0, 0xc420012420, 0x2f, 0x0, 0x0,
	//  0x0, 0x0, 0x0, 0x0, ...)
	//      /tmp/babel-79540aE/go-src-7954DTN.go:22 +0xbb
	//  main.isGloballyExec(0x4d1313, 0xf, 0x4daecc, 0x30, 0x4c5800)
	//      /tmp/babel-79540aE/go-src-7954DTN.go:39 +0xc5
	//  main.runJob(0x4cfada, 0x1, 0x4d4c19, 0x22)
	//      /tmp/babel-79540aE/go-src-7954DTN.go:51 +0x4b
	//  main.main()
	//      /tmp/babel-79540aE/go-src-7954DTN.go:71 +0x63
	//  ", Misc:map[string]interface {}{}}}, Message:"cannot run job \"1\":
	//  requisite binaries not available", StackTrace:"goroutine 1 [running]:
	//  runtime/debug.Stack(0x4d63f0, 0x33, 0xc420045e40)
	//      /home/kate/.guix-profile/src/runtime/debug/stack.go:24 +0x79
	//  main.wrapError(0x530380, 0xc42000a370, 0x4d63f0, 0x33,
	//  0xc420045e40, 0x1, 0x1, 0x0, 0x0, 0x0, ...)
	//      /tmp/babel-79540aE/go-src-7954DTN.go:22 +0xbb
	//  main.runJob(0x4cfada, 0x1, 0x4d4c19, 0x22)
	//      /tmp/babel-79540aE/go-src-7954DTN.go:53 +0x356
	//  main.main()
	//      /tmp/babel-79540aE/go-src-7954DTN.go:71 +0x63
	//  ", Misc:map[string]interface {}{}}}
	//But our error message is now exactly what we want users to see:
	//[1] cannot run job "1": requisite binaries not available

	//Timeouts and Cancellation
	//Take the following code, and assume it’s running in its own goroutine:
	//var value interface{}
	//select {
	//case <-done:
	//	return
	//case value = <-valueStream:
	//}
	//
	//result := reallyLongCalculation(value)
	//
	//select {
	//case <-done:
	//	return
	//case resultStream<-result:
	//}
	//We’ve dutifully coupled the read from valueStream and the write to resultStream with a check against the
	//done channel to see if the goroutine has been canceled, but we still have a problem. reallyLongCalculation
	//doesn’t look to be preemptable, and, according to the name, it looks like it might take a really long time!
	//This means that if something attempts to cancel this goroutine while reallyLongCalculation is executing, it
	//could be a very long time before we acknowledge the cancellation and halt. Let’s try and make
	//reallyLongCaluclation preemptable and see what happens:
	//reallyLongCalculation := func(done <-chan interface{}, value interface{}) interface{} {
	//	intermediateResult := longCalculation(value)
	//	select {
	//	case <-done:
	//		return nil
	//	default:
	//	}
	//
	//	return longCaluclation(intermediateResult)
	//}
	//We’ve made some progress: reallyLongCaluclation is now preemptable, but we can see that we’ve only halved
	//the problem: we can only preempt reallyLongCalculation in between calls to other, seemingly long-running,
	//function calls. To solve this, we need to make longCalculation preemptable as well:
	//reallyLongCalculation := func(done <-chan interface{}, value interface{}) interface{} {
	//	intermediateResult := longCalculation(done, value)
	//	return longCaluclation(done, intermediateResult)
	//}
	//there’s another problem lurking here as well: if our goroutine happens to modify shared state—e.g.,
	//a database, a file, an in-memory data structure—what happens when the goroutine is canceled? Does
	//your goroutine try and roll back the intermediary work it’s done? How long does it have to do this
	//work? Something has told the goroutine that it should halt, so the goroutine shouldn’t take too long
	//to roll back its work, right?
	//It’s difficult to give general advice on how to handle this problem because the nature of your algorithm
	//will dictate so much of how you handle this situation; however, if you keep your modifications to any
	//shared state within a tight scope, and/or ensure those modifications are easily rolled back, you can
	//usually handle cancellations pretty well. If possible, build up intermediate results in-memory and
	//then modify state as quickly as possible. As an example, here is the wrong way to do it:
	//result := add(1, 2, 3)
	//writeTallyToState(result)
	//result = add(result, 4, 5, 6)
	//writeTallyToState(result)
	//result = add(result, 7, 8, 9)
	//writeTallyToState(result)
	//Here we write to state three times. If a goroutine running this code were canceled before the final write,
	//we’d need to somehow roll back the previous two calls to writeTallyToState. Contrast that approach with this:
	//result := add(1, 2, 3, 4, 5, 6, 7, 8, 9)
	//writeTallyToState(result)
	//Here the surface area we have to worry about rolling back is much smaller. If the cancellation comes
	//in after our call to writeToState, we still need a way to back out our changes, but the probability
	//that this will happen is much smaller since we only modify state once.

	//Heartbeats
	//Heartbeats are a way for concurrent processes to signal life to outside parties. They get their name from
	//human anatomy wherein a heartbeat signifies life to an observer. Heartbeats have been around since before
	//Go, and remain useful within it.
	//There are a few different reasons heartbeats are interesting for concurrent code. They allow us insights
	//into our system, and they can make testing the system deterministic when it might otherwise not be.
	//There are two different types of heartbeats we’ll discuss in this section:
	//Heartbeats that occur on a time interval.
	//Heartbeats that occur at the beginning of a unit of work.
	//Heartbeats that occur on a time interval are useful for concurrent code that might be waiting for something
	//else to happen for it to process a unit of work. Because you don’t know when that work might come in,
	//your goroutine might be sitting around for a while waiting for something to happen. A heartbeat is a
	//way to signal to its listeners that everything is well, and that the silence is expected.
	//The following code demonstrates a goroutine that exposes a heartbeat:
	//doWork := func(done <-chan interface{}, pulseInterval time.Duration, ) (<-chan interface{}, <-chan time.Time) {
	//	heartbeat := make(chan interface{}) //Here we set up a channel to send heartbeats on. We return this out of doWork.
	//	results := make(chan time.Time)
	//	go func() {
	//		defer close(heartbeat)
	//		defer close(results)
	//
	//		pulse := time.Tick(pulseInterval) //Here we set the heartbeat to pulse at the pulseInterval we were given. Every pulseInterval there will be something to read on this channel.
	//		workGen := time.Tick(2*pulseInterval) //This is just another ticker used to simulate work coming in. We choose a duration greater than the pulseInterval so that we can see some heartbeats coming out of the goroutine.
	//
	//		sendPulse := func() {
	//			select {
	//			case heartbeat <-struct{}{}:
	//			default: //Note that we include a default clause. We must always guard against the fact that no one may be listening to our heartbeat. The results emitted from the goroutine are critical, but the pulses are not.
	//			}
	//		}
	//		sendResult := func(r time.Time) {
	//			for {
	//				select {
	//				case <-done:
	//					return
	//				case <-pulse: //Just like with done channels, anytime you perform a send or receive, you also need to include a case for the heartbeat’s pulse.
	//					sendPulse()
	//				case results <- r:
	//					return
	//				}
	//			}
	//		}
	//
	//		for {
	//			select {
	//			case <-done:
	//				return
	//			case <-pulse: //Just like with done channels, anytime you perform a send or receive, you also need to include a case for the heartbeat’s pulse.
	//				sendPulse()
	//			case r := <-workGen:
	//				sendResult(r)
	//			}
	//		}
	//	}()
	//	return heartbeat, results
	//}
	//Notice that because we might be sending out multiple pulses while we wait for input, or multiple pulses
	//while waiting to send results, all the select statements need to be within for loops. Looking good so
	//far; how do we utilize this function and consume the events it emits? Let’s take a look:
	//done := make(chan interface{})
	//time.AfterFunc(10*time.Second, func() { close(done) }) //We set up the standard done channel and close it after 10 seconds. This gives our goroutine time to do some work.
	//
	//const timeout = 2*time.Second //Here we set our timeout period. We’ll use this to couple our heartbeat interval to our timeout.
	//heartbeat, results := doWork(done, timeout/2) //We pass in timeout/2 here. This gives our heartbeat an extra tick to respond so that our timeout isn’t too sensitive.
	//for {
	//	select {
	//	case _, ok := <-heartbeat: //Here we select on the heartbeat. When there are no results, we are at least guaranteed a message from the heartbeat channel every timeout/2. If we don’t receive it, we know there’s something wrong with the goroutine itself.
	//		if ok == false {
	//			return
	//		}
	//		fmt.Println("pulse")
	//	case r, ok := <-results: //Here we select from the results channel; nothing fancy going on here.
	//		if ok == false {
	//			return
	//		}
	//		fmt.Printf("results %v\n", r.Second())
	//	case <-time.After(timeout): //Here we time out if we haven’t received either a heartbeat or a new result.
	//		return
	//	}
	//}
	//Running this code produces:
	//pulse
	//pulse
	//results 52
	//pulse
	//pulse
	//results 54
	//pulse
	//pulse
	//results 56
	//pulse
	//pulse
	//results 58
	//pulse
	//You can see that we receive about two pulses per result as we intended.
	//Now in a properly functioning system, heartbeats aren’t that interesting. We might use them to gather
	//statistics regarding idle time, but the utility for interval-based heartbeats really shines when
	//your goroutine isn’t behaving as expected.
	//Consider the next example. We’ll simulate an incorrectly written goroutine with a panic by stopping
	//the goroutine after only two iterations, and then not closing either of our channels. Let’s have a look:

	//doWork := func(done <-chan interface{}, pulseInterval time.Duration, ) (<-chan interface{}, <-chan time.Time) {
	//	heartbeat := make(chan interface{})
	//	results := make(chan time.Time)
	//	go func() {
	//		pulse := time.Tick(pulseInterval)
	//		workGen := time.Tick(2 * pulseInterval)
	//
	//		sendPulse := func() {
	//			select {
	//			case heartbeat <- struct{}{}:
	//			default:
	//			}
	//		}
	//		sendResult := func(r time.Time) {
	//			for {
	//				select {
	//				case <-pulse:
	//					sendPulse()
	//				case results <- r:
	//					return
	//				}
	//			}
	//		}
	//
	//		for i := 0; i < 2; i++ { //Here is our simulated panic. Instead of infinitely looping until we’re asked to stop, as in the previous example, we’ll only loop twice.
	//			select {
	//			case <-done:
	//				return
	//			case <-pulse:
	//				sendPulse()
	//			case r := <-workGen:
	//				sendResult(r)
	//			}
	//		}
	//	}()
	//	return heartbeat, results
	//}
	//
	//done := make(chan interface{})
	//time.AfterFunc(10*time.Second, func() { close(done) })
	//
	//const timeout = 2 * time.Second
	//heartbeat, results := doWork(done, timeout/2)
	//for {
	//	select {
	//	case _, ok := <-heartbeat:
	//		if ok == false {
	//			return
	//		}
	//		fmt.Println("pulse")
	//	case r, ok := <-results:
	//		if ok == false {
	//			return
	//		}
	//		fmt.Printf("results %v\n", r)
	//	case <-time.After(timeout):
	//		fmt.Println("worker goroutine is not healthy!")
	//		return
	//	}
	//}
	//Beautiful! Within two seconds our system realizes something is amiss with our goroutine and breaks the
	//for-select loop. By using a heartbeat, we have successfully avoided a deadlock, and we remain deterministic
	//by not having to rely on a longer timeout.
	//Also note that heartbeats help with the opposite case: they let us know that long-running goroutines
	//remain up, but are just taking a while to produce a value to send on the values channel.
	//Now let’s shift over to looking at heartbeats that occur at the beginning of a unit of work. These
	//are extremely useful for tests. Here’s an example that sends a pulse before every unit of work:
	//doWork := func(done <-chan interface{}) (<-chan interface{}, <-chan int) {
	//	heartbeatStream := make(chan interface{}, 1) //Here we create the heartbeat channel with a buffer of one. This ensures that there’s always at least one pulse sent out even if no one is listening in time for the send to occur.
	//	workStream := make(chan int)
	//	go func () {
	//		defer close(heartbeatStream)
	//		defer close(workStream)
	//
	//		for i := 0; i < 10; i++ {
	//			select { //Here we set up a separate select block for the heartbeat. We don’t want to include this in the same select block as the send on results because if the receiver isn’t ready for the result, they’ll receive a pulse instead, and the current value of the result will be lost. We also don’t include a case statement for the done channel since we have a default case that will just fall through.
	//			case heartbeatStream <- struct{}{}:
	//			default: //Once again we guard against the fact that no one may be listening to our heartbeats. Because our heartbeat channel was created with a buffer of one, if someone is listening, but not in time for the first pulse, they’ll still be notified of a pulse.
	//			}
	//
	//			select {
	//			case <-done:
	//				return
	//			case workStream <- rand.Intn(10):
	//			}
	//		}
	//	}()
	//
	//	return heartbeatStream, workStream
	//}
	//
	//done := make(chan interface{})
	//defer close(done)
	//
	//heartbeat, results := doWork(done)
	//for {
	//	select {
	//	case _, ok := <-heartbeat:
	//		if ok {
	//			fmt.Println("pulse")
	//		} else {
	//			return
	//		}
	//	case r, ok := <-results:
	//		if ok {
	//			fmt.Printf("results %v\n", r)
	//		} else {
	//			return
	//		}
	//	}
	//}
	//Running this code produces:
	//You can see in this example that we receive one pulse for every result, as intended.

	//Where this technique really shines is in writing tests. Interval-based heartbeats can be used in
	//the same fashion, but if you only care that the goroutine has started doing its work, this style
	//of heartbeat is simple. Consider the following snippet of code:
	// see DoWork()
	//The DoWork function is a pretty simple generator that converts the numbers we pass in to a stream
	//on the channel it returns. Let’s try testing this function. Here’s an example of a bad test:
	// see TestDoWork_GeneratesAllNumbers()
	// run ->  go test case_test.go example-4.go
	//Running this test produces:
	//--- FAIL: TestDoWork_GeneratesAllNumbers (1.00s)
	//    case_test.go:30: test timed out
	//FAIL
	//FAIL    command-line-arguments  12.845s
	//This test is bad because it’s nondeterministic. In our example function, I’ve ensured this test will
	//always fail, but if I were to remove the time.Sleep, the situation actually gets worse: this test
	//will pass at times, and fail at others.

	//We mentioned earlier how factors external to the process can cause the goroutine to take longer
	//to get to its first iteration. Even whether or not the goroutine is scheduled in the first place
	//is a concern. The point is that we can’t be guaranteed that the first iteration of the goroutine
	//will occur before our timeout is reached, and so we begin thinking in terms of probabilities: how
	//likely is it that this timeout will be significant? We could increase the timeout, but that means
	//failures will take a long time, thereby slowing down our test suite.
	//Fortunately with a heartbeat this is easily solved. Here is a test that is deterministic:
	// see TestDoWork_GeneratesAllNumbers1()
	//Because of the heartbeat, we can safely write our test without timeouts. The only risk we run is of
	//one of our iterations taking an inordinate amount of time. If that’s important to us, we can utilize
	//the safer interval-based heartbeats and achieve perfect safety.
	//Here is an example of a test utilizing interval-based heartbeats:
	// see DoWork1()

	//Replicated Requests
	//For some applications, receiving a response as quickly as possible is the top priority.
	//For example, maybe the application is servicing a user’s HTTP request, or retrieving a
	//replicated blob of data. In these instances you can make a trade-off: you can replicate
	//the request to multiple handlers (whether those be goroutines, processes, or servers),
	//and one of them will return faster than the other ones; you can then immediately return
	//the result. The downside is that you’ll have to utilize resources to keep multiple copies
	//of the handlers running.
	//If this replication is done in-memory, it might not be that costly, but if replicating the handlers
	//requires replicating processes, servers, or even data centers, this can become quite costly.
	//The decision you’ll have to make is whether or not the cost is worth the benefit.
	//Let’s look at how you can replicate requests within a single process. We’ll use multiple goroutines
	//to serve as request handlers, and the goroutines will sleep for a random amount of time between one
	//and six nanoseconds to simulate load. This will give us handlers that return a result at various
	//times and will allow us to see how this can lead to faster results.
	//Here’s an example that replicates a simulated request over 10 handlers:
	//doWork := func(done <-chan interface{}, id int, wg *sync.WaitGroup, result chan<- int) {
	//	started := time.Now()
	//	defer wg.Done()
	//
	//	// Simulate random load
	//	simulatedLoadTime := time.Duration(1+rand.Intn(5)) * time.Second
	//	select {
	//	case <-done:
	//	case <-time.After(simulatedLoadTime):
	//	}
	//
	//	select {
	//	case <-done:
	//	case result <- id:
	//	}
	//
	//	took := time.Since(started)
	//	// Display how long handlers would have taken
	//	if took < simulatedLoadTime {
	//		took = simulatedLoadTime
	//	}
	//	fmt.Printf("%v took %v\n", id, took)
	//}
	//
	//done := make(chan interface{})
	//result := make(chan int)
	//
	//var wg sync.WaitGroup
	//wg.Add(10)
	//
	//for i := 0; i < 10; i++ { //Here we start 10 handlers to handle our requests.
	//	go doWork(done, i, &wg, result)
	//}
	//
	//firstReturned := <-result //This line grabs the first returned value from the group of handlers.
	//close(done)               //Here we cancel all the remaining handlers. This ensures they don’t continue to do unnecessary work.
	//wg.Wait()
	//
	//fmt.Printf("Received an answer from #%v\n", firstReturned)
	//In this run, it looks like handler #8 returned fastest. Note that in the output we’re displaying how long
	//each handler would have taken so that you can get a sense of how much time this technique can save.
	//Imagine if you only spun up one handler and it happened to be handler #5. Instead of waiting just over
	//a second for the request to be handled, you would have had to wait for five seconds.
	//The only caveat to this approach is that all of your handlers need to have equal opportunity to service
	//the request. In other words, you’re not going to have a chance at receiving the fastest time from a
	//handler that can’t service the request. As I mentioned, whatever resources the handlers are using to
	//do their job need to be replicated as well.
	//A different symptom of the same problem is uniformity. If your handlers are too much alike, the chances
	//that any one will be an outlier is smaller. You should only replicate out requests like this to handlers
	//that have different runtime conditions: different processes, machines, paths to a data store, or access
	//to different data stores altogether.
	//Although this is can be expensive to set up and maintain, if speed is your goal, this is a valuable
	//technique. In addition, this naturally provides fault tolerance and scalability.

	//Rate Limiting
	//If you’ve ever worked with an API for a service, you’ve likely had to contend with rate limiting,
	//which constrains the number of times some kind of resource is accessed to some finite number per
	//unit of time. The resource can be anything: API connections, disk reads/writes, network packets, errors.
	//Have you ever wondered why services put rate limits in place? Why not allow unfettered access to a
	//system? The most obvious answer is that by rate limiting a system, you prevent entire classes of
	//attack vectors against your system. If malicious users can access your system as quickly as their
	//resources allow it, they can do all kinds of things.
	//For example, they could fill up your service’s disk either with log messages or valid requests.
	//If you’ve misconfigured your log rotation, they could even perform something malicious and then
	//make enough requests that any record of the activity would be rotated out of the log and into /dev/null.
	//They could attempt to brute-force access to a resource, or maybe they would just perform a distributed
	//denial of service attack. The point is: if you don’t rate limit requests to your system, you cannot
	//easily secure it.
	//Malicious use isn’t the only reason. In distributed systems, a legitimate user could degrade the
	//performance of the system for other users if they’re performing operations at a high enough volume,
	//or if the code they’re exercising is buggy. This can even cause the death-spirals we discussed earlier.
	//From a product standpoint, this is awful! Usually you want to make some kind of guarantees to your users
	//about what kind of performance they can expect on a consistent basis. If one user can affect that
	//agreement, you’re in for a bad time. A user’s mental model is usually that their access to the system
	//is sandboxed and can neither affect nor be affected by other users’ activities. If you break that mental
	//model, your system can feel like it’s not well engineered, and even cause users to become angry or leave.
	//So how do we go about implementing rate limits in Go?
	//Let’s assume that to utilize a resource, you have to have an access token for the resource. Without
	//the token, your request is denied. Now imagine these tokens are stored in a bucket waiting to be
	//retrieved for usage. The bucket has a depth of d, which indicates it can hold d access tokens at
	//a time. For example, if the bucket has a depth of five, it can hold five tokens.
	//Now, every time you need to access a resource, you reach into the bucket and remove a token. If
	//your bucket contains five tokens, and you access the resource five times, you’d be able to do so;
	//but on the sixth try, no access token would be available. You either have to queue your request
	//until a token becomes available, or deny the request.
	//Let’s put this algorithm to use and see how a Go program might behave when written against an
	//implementation of the token bucket algorithm.
	//Let’s pretend we have access to an API, and a Go client has been provided to utilize it. This
	//API has two endpoints: one for reading a file, and one for resolving a domain name to an IP
	//address. For simplicity’s sake, I’m going to leave off any arguments and return values that
	//would be needed to actually access a service. So here’s our client:
	// see ReadFile(), ResolveAddress(), APIConnection{}, Open()
	//Since in theory this request is going over the wire, we take a context.Context in as the first argument
	//in case we need to cancel the request or pass values over to the server. Pretty standard stuff.
	//We’ll now create a simple driver to access this API. The driver needs to read 10 files and
	//resolve 10 addresses, but the files and addresses have no relation to each other and so the
	//driver can make these API calls concurrent to one another. Later this will help stress our
	//APIClient and exercise our rate limiter.
	//defer log.Printf("Done.")
	//log.SetOutput(os.Stdout)
	//log.SetFlags(log.Ltime | log.LUTC)
	//
	//apiConnection := Open()
	//var wg sync.WaitGroup
	//wg.Add(20)
	//
	//for i := 0; i < 10; i++ {
	//	go func() {
	//		defer wg.Done()
	//		err := apiConnection.ReadFile(context.Background())
	//		if err != nil {
	//			log.Printf("cannot ReadFile: %v", err)
	//		}
	//		log.Printf("ReadFile")
	//	}()
	//}
	//
	//for i := 0; i < 10; i++ {
	//	go func() {
	//		defer wg.Done()
	//		err := apiConnection.ResolveAddress(context.Background())
	//		if err != nil {
	//			log.Printf("cannot ResolveAddress: %v", err)
	//		}
	//		log.Printf("ResolveAddress")
	//	}()
	//}
	//
	//wg.Wait()
	//We can see that all API requests are fielded almost simultaneously. We have no rate limiting set up
	//and so our clients are free to access the system as frequently as they like. Now is a good time to
	//remind you that a bug could exist in our driver that could result in an infinite loop. Without rate
	//limiting, I could be staring down a nasty bill.
	//OK, so let’s introduce a rate limiter! I’m going to do so within the APIConnection, but
	//normally a rate limiter would be running on a server so the users couldn’t trivially bypass it.
	//Production systems might also include a client-side rate limiter to help prevent the client from
	//making unnecessary calls only to be denied, but that is an optimization. For our purposes, a
	//client-side rate limiter keeps things simple.
	//We’re going to be looking at examples that use an implementation of a token bucket rate limiter from
	//the golang.org/x/time/rate package. I chose this package because this is as close to the standard
	//library as I could get. There are certainly other packages out there that do the same thing with more
	//bells and whistles, and those may serve you better for use in production systems. The golang.org/x/time/rate
	//package is pretty simple, so it should work well for our purposes.
	//The first two ways we’ll interact with this package are the Limit type and the NewLimiter function, defined here:
	// Limit defines the maximum frequency of some events.  Limit is
	// represented as number of events per second.  A zero Limit allows no
	// events.
	//type Limit float64

	// NewLimiter returns a new Limiter that allows events up to rate r
	// and permits bursts of at most b tokens.
	//func NewLimiter(r Limit, b int) *Limiter
	//In NewLimiter, we see two familiar parameters: r and b. r is the rate we discussed previously,
	//and b is the bucket depth we discussed.
	//The rates package also defines a helper method, Every, to assist in converting a time.Duration into a Limit:
	// Every converts a minimum time interval between events to a Limit.
	//func Every(interval time.Duration) Limit
	//The Every function makes sense, but I want to discuss rate limits in terms of the number of operations
	//per time measurement, not the interval between requests. We can express this as the following:
	//rate.Limit(events/timePeriod.Seconds())
	//But I don’t want to type that every time, and the Every function has some special logic that will return
	//rate.Inf—an indication that there is no limit—if the interval provided is zero. Because of this, we’ll
	//express our helper function in terms of the Every function:
	// check par()
	//After we create a rate.Limiter, we’ll want to use it to block our requests until we’re given an access token.
	//We can do that with the Wait method, which simply calls WaitN with an argument of 1:
	//After we create a rate.Limiter, we’ll want to use it to block our requests until we’re given an access
	//token. We can do that with the Wait method, which simply calls WaitN with an argument of 1:
	// Wait is shorthand for WaitN(ctx, 1).
	//func (lim *Limiter) Wait(ctx context.Context)
	// WaitN blocks until lim permits n events to happen.
	// It returns an error if n exceeds the Limiter's burst size, the Context is
	// canceled, or the expected wait time exceeds the Context's Deadline.
	//func (lim *Limiter) WaitN(ctx context.Context, n int) (err error)
	//We should now have all the ingredients we’ll need to begin rate limiting our API requests.
	//Let’s modify our APIConnection type and give it a try!
	//defer log.Printf("Done.")
	//log.SetOutput(os.Stdout)
	//log.SetFlags(log.Ltime | log.LUTC)
	//
	//apiConnection := Open1()
	//var wg sync.WaitGroup
	//wg.Add(20)
	//
	//for i := 0; i < 10; i++ {
	//	go func() {
	//		defer wg.Done()
	//		err := apiConnection.ReadFile1(context.Background())
	//		if err != nil {
	//			log.Printf("cannot ReadFile: %v", err)
	//		}
	//		log.Printf("ReadFile")
	//	}()
	//}
	//
	//for i := 0; i < 10; i++ {
	//	go func() {
	//		defer wg.Done()
	//		err := apiConnection.ResolveAddress1(context.Background())
	//		if err != nil {
	//			log.Printf("cannot ResolveAddress: %v", err)
	//		}
	//		log.Printf("ResolveAddress")
	//	}()
	//}
	//
	//wg.Wait()
	//You can see that whereas before we were fielding all of our API requests simultaneously, we’re now
	//completing a request once a second. It looks like our rate limiter is working!
	//This gets us very basic rate limiting, but in production we’re likely going to want something a little
	//more complex. We will probably want to establish multiple tiers of limits: fine-grained controls to
	//limit requests per second, and coarse-grained controls to limit requests per minute, hour, or day.

	//In certain instances, it’s possible to do this with a single rate limiter; however, it’s not possible in
	//all cases, and by attempting to roll the semantics of limits per unit of time into a single layer, you
	//lose a lot of information around the intent of the rate limiter. For these reasons, I find it easier
	//to keep the limiters separate and then combine them into one rate limiter that manages the interaction
	//for you. To this end I’ve created a simple aggregate rate limiter called multiLimiter. Here is the definition:

	//The Wait method loops through all the child rate limiters and calls Wait on each of them.
	//These calls may or may not block, but we need to notify each rate limiter of the request so we can
	//decrement our token bucket. By waiting for each limiter, we are guaranteed to wait for exactly the
	//time of the longest wait. This is because if we perform smaller waits that only wait for segments of
	//the longest wait and then hit the longest wait, the longest wait will be recalculated to only be the
	//remaining time. This is because while the earlier waits were blocking, the latter waits were refilling
	//their buckets; any waits after will be returned instantaneously.
	//Now that we have the means to express rate limits from multiple rate limits, let’s take the opportunity
	//to do so. Let’s redefine our APIConnection to have limits both per second and per minute:

	//defer log.Printf("Done.")
	//log.SetOutput(os.Stdout)
	//log.SetFlags(log.Ltime | log.LUTC)
	//
	//apiConnection := Open2()
	//var wg sync.WaitGroup
	//wg.Add(20)
	//
	//for i := 0; i < 10; i++ {
	//	go func() {
	//		defer wg.Done()
	//		err := apiConnection.ReadFile2(context.Background())
	//		if err != nil {
	//			log.Printf("cannot ReadFile: %v", err)
	//		}
	//		log.Printf("ReadFile")
	//	}()
	//}
	//
	//for i := 0; i < 10; i++ {
	//	go func() {
	//		defer wg.Done()
	//		err := apiConnection.ResolveAddress2(context.Background())
	//		if err != nil {
	//			log.Printf("cannot ResolveAddress: %v", err)
	//		}
	//		log.Printf("ResolveAddress")
	//	}()
	//}
	//
	//wg.Wait()
	//As you can see we make two requests per second up until request #11, at which point we begin making requests
	//every six seconds. This is because we drained our available pool of per-minute request tokens, and become
	//limited by this cap.
	//It might be slightly counterintuitive why request #11 occurs after only two seconds rather than six like
	//the rest of the requests. Remember that although we limit our API requests to 10 a minute, that minute
	//is a sliding window of time. By the time we reach the eleventh request, our per-minute rate limiter has
	//accrued another token.
	//Defining limits like this allows us to express our coarse-grained limits plainly while still limiting the
	//number of requests at a finer level of detail.
	//This technique also allows us to begin thinking across dimensions other than time. When you rate limit a
	//system, you’re probably going to limit more than one thing. You’ll likely have some kind of limit on the
	// number of API requests, but in addition, you’ll probably also have limits on other resources like disk
	//access, network access, etc. Let’s flesh out our example a bit and set up rate limits for disk and network:

	//Healing Unhealthy Goroutines
	//In long-lived processes such as daemons, it’s very common to have a set of long-lived goroutines.
	//These goroutines are usually blocked, waiting on data to come to them through some means, so that
	//they can wake up, do their work, and then pass the data on. Sometimes the goroutines are dependent
	//on a resource that you don’t have very good control of. Maybe a goroutine receives a request to
	//pull data from a web service, or maybe it’s monitoring an ephemeral file. The point is that it can
	//be very easy for a goroutine to become stuck in a bad state from which it cannot recover without
	//external help. If you separate your concerns, you might even say that it shouldn’t be the concern
	//of a goroutine doing work to know how to heal itself from a bad state. In a long-running process,
	//it can be useful to create a mechanism that ensures your goroutines remain healthy and restarts
	//them if they become unhealthy. We’ll refer to this process of restarting goroutines as “healing.”
	//To heal goroutines, we’ll use our heartbeat pattern to check up on the liveliness of the goroutine we’re
	//monitoring. The type of heartbeat will be determined by what you’re trying to monitor, but if your
	//goroutine can become livelocked, make sure that the heartbeat contains some kind of information indicating
	//that the goroutine is not only up, but doing useful work. In this section, for simplicity, we’ll only
	//consider whether goroutines are live or dead.

	//We’ll call the logic that monitors a goroutine’s health a steward, and the goroutine that it monitors a
	//ward. Stewards will also be responsible for restarting a ward’s goroutine should it become unhealthy. To
	//do so, it will need a reference to a function that can start the goroutine. Let’s see what a steward
	//might look like:
	//type startGoroutineFn func(done <-chan interface{}, pulseInterval time.Duration) (heartbeat <-chan interface{}) //Here we define the signature of a goroutine that can be monitored and restarted. We see the familiar done channel, and pulseInterval and heartbeat from the heartbeat pattern.
	//
	//newSteward := func(timeout time.Duration, startGoroutine startGoroutineFn) startGoroutineFn { //On this line we see that a steward takes in a timeout for the goroutine it will be monitoring, and a function, startGoroutine, to start the goroutine it’s monitoring. Interestingly, the steward itself returns a startGoroutineFn indicating that the steward itself is also monitorable.
	//	return func(done <-chan interface{}, pulseInterval time.Duration) <-chan interface{} {
	//		heartbeat := make(chan interface{})
	//		go func() {
	//			defer close(heartbeat)
	//			var wardDone chan interface{}
	//			var wardHeartbeat <-chan interface{}
	//			startWard := func() { //Here we define a closure that encodes a consistent way to start the goroutine we’re monitoring.
	//				wardDone = make(chan interface{})                             //This is where we create a new channel that we’ll pass into the ward goroutine in case we need to signal that it should halt.
	//				wardHeartbeat = startGoroutine(or(wardDone, done), timeout/2) //Here we start the goroutine we’ll be monitoring. We want the ward goroutine to halt if either the steward is halted, or the steward wants to halt the ward goroutine, so we wrap both done channels in a logical-or. The pulseInterval we pass in is half of the timeout period, although as we discussed in “Heartbeats”, this can be tweaked.
	//			}
	//			startWard()
	//			pulse := time.Tick(pulseInterval)
	//
	//		monitorLoop:
	//			for {
	//				timeoutSignal := time.After(timeout)
	//
	//				for { //This is our inner loop, which ensures that the steward can send out pulses of its own.
	//					select {
	//					case <-pulse:
	//						select {
	//						case heartbeat <- struct{}{}:
	//						default:
	//						}
	//					case <-wardHeartbeat: //Here we see that if we receive the ward’s pulse, we continue our monitoring loop.
	//						continue monitorLoop
	//					case <-timeoutSignal: //This line indicates that if we don’t receive a pulse from the ward within our timeout period, we request that the ward halt and we begin a new ward goroutine. We then continue monitoring.
	//						log.Println("steward: ward unhealthy; restarting")
	//						close(wardDone)
	//						startWard()
	//						continue monitorLoop
	//					case <-done:
	//						return
	//					}
	//				}
	//			}
	//		}()
	//
	//		return heartbeat
	//	}
	//}
	//Our for loop is a little busy, but as long as you’re familiar with the patterns involved, it’s relatively
	//straightforward to read through. Let’s give our steward a test run. What happens if we monitor a
	//goroutine that is misbehaving? Let’s take a look:
	//log.SetOutput(os.Stdout)
	//log.SetFlags(log.Ltime | log.LUTC)
	//
	//doWork := func(done <-chan interface{}, _ time.Duration) <-chan interface{} {
	//	log.Println("ward: Hello, I'm irresponsible!")
	//	go func() {
	//		<-done //Here we see that this goroutine isn’t doing anything but waiting to be canceled. It’s also not sending out any pulses.
	//		log.Println("ward: I am halting.")
	//	}()
	//	return nil
	//}
	//doWorkWithSteward := newSteward(4*time.Second, doWork) //This line creates a function that will create a steward for the goroutine doWork starts. We set the timeout for doWork at four seconds.
	//
	//done := make(chan interface{})
	//time.AfterFunc(9*time.Second, func() { //Here we halt the steward and its ward after nine seconds so that our example will end.
	//	log.Println("main: halting steward and ward.")
	//	close(done)
	//})
	//
	//for range doWorkWithSteward(done, 4*time.Second) {
	//} //Finally, we start the steward and range over its pulses to prevent our example from halting.
	//log.Println("Done")
}

func or(done chan interface{}, done2 <-chan interface{}) <-chan interface{} {
	d := make(chan interface{})
	go func() {
		select {
		case <-done:
			close(d)
		case <-done2:
			close(d)
		}
	}()
	return d
}

type MyError struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]interface{}
}

func wrapError(err error, messagef string, msgArgs ...interface{}) MyError {
	return MyError{
		Inner:      err, // Here we store the error we’re wrapping. We always want to be able to get back to the lowest-level error in case we need to investigate what happened.
		Message:    fmt.Sprintf(messagef, msgArgs...),
		StackTrace: string(debug.Stack()),        //This line of code takes note of the stack trace when the error was created. A more sophisticated error type might elide the stack-frame from wrapError.
		Misc:       make(map[string]interface{}), //Here we create a catch-all for storing miscellaneous information. This is where we might store the concurrent ID, a hash of the stack trace, or other contextual information that might help in diagnosing the error.
	}
}

func (err MyError) Error() string {
	return err.Message
}

type LowLevelErr struct {
	error
}

func isGloballyExec(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, LowLevelErr{wrapError(err, err.Error())} //Here we wrap the raw error from calling os.Stat with a customized error. In this case we are OK with the message coming out of this error, and so we won’t mask it.
	}
	return info.Mode().Perm()&0100 == 0100, nil
}

type IntermediateErr struct {
	error
}

func runJob(id string) error {
	const jobBinPath = "/bad/job/binary"
	isExecutable, err := isGloballyExec(jobBinPath)
	if err != nil {
		return IntermediateErr{wrapError( // Here we are now customizing the error with a crafted message. In this case, we want to obfuscate the low-level details of why the job isn’t running because we feel it’s not important information to consumers of our module.
			err,
			"cannot run job %q: requisite binaries not available",
			id,
		)}
	} else if isExecutable == false {
		return wrapError(
			nil,
			"cannot run job %q: requisite binaries are not executable",
			id,
		)
	}

	return exec.Command(jobBinPath, "--id="+id).Run()
}

func handleError(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logID: %v]: ", key))
	log.Printf("%#v", err) //Here we log out the full error in case someone needs to dig into what happened.
	fmt.Printf("[%v] %v", key, message)
}

func DoWork(done <-chan interface{}, nums ...int) (<-chan interface{}, <-chan int) {
	heartbeat := make(chan interface{}, 1)
	intStream := make(chan int)
	go func() {
		defer close(heartbeat)
		defer close(intStream)

		time.Sleep(2 * time.Second) //Here we simulate some kind of delay before the goroutine can begin working. In practice this can be all kinds of things and is nondeterministic. I’ve seen delays caused by CPU load, disk contention, network latency, and goblins.

		for _, n := range nums {
			select {
			case heartbeat <- struct{}{}:
			default:
			}

			select {
			case <-done:
				return
			case intStream <- n:
			}
		}
	}()

	return heartbeat, intStream
}

func DoWork1(done <-chan interface{}, pulseInterval time.Duration, nums ...int, ) (<-chan interface{}, <-chan int) {
	heartbeat := make(chan interface{}, 1)
	intStream := make(chan int)
	go func() {
		defer close(heartbeat)
		defer close(intStream)

		time.Sleep(2 * time.Second)

		pulse := time.Tick(pulseInterval)
	numLoop: // We’re using a label here to make continuing from the inner loop a little simpler.
		for _, n := range nums {
			for { //We require two loops: one to range over our list of numbers, and this inner loop to run until the number is successfully sent on the intStream.
				select {
				case <-done:
					return
				case <-pulse:
					select {
					case heartbeat <- struct{}{}:
					default:
					}
				case intStream <- n:
					continue numLoop //Here we continue executing the outer loop.
				}
			}
		}
	}()

	return heartbeat, intStream
}

func Open() *APIConnection {
	return &APIConnection{}
}

type APIConnection struct{}

func (a *APIConnection) ReadFile(ctx context.Context) error {
	// Pretend we do work here
	return nil
}

func (a *APIConnection) ResolveAddress(ctx context.Context) error {
	// Pretend we do work here
	return nil
}

func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}

func Open1() *APIConnection1 {
	return &APIConnection1{
		rateLimiter: rate.NewLimiter(rate.Limit(1), 1), //Here we set the rate limit for all API connections to one event per second.
	}
}

type APIConnection1 struct {
	rateLimiter *rate.Limiter
}

func (a *APIConnection1) ReadFile1(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil { //Here we wait on the rate limiter to have enough access tokens for us to complete our request.
		return err
	}
	// Pretend we do work here
	return nil
}

func (a *APIConnection1) ResolveAddress1(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil { //Here we wait on the rate limiter to have enough access tokens for us to complete our request.
		return err
	}
	// Pretend we do work here
	return nil
}

type RateLimiter interface { //Here we define a RateLimiter interface so that a MultiLimiter can recursively define other MultiLimiter instances.
	Wait(context.Context) error
	Limit() rate.Limit
}

func MultiLimiter(limiters ...RateLimiter) *multiLimiter {
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}
	sort.Slice(limiters, byLimit) // Here we implement an optimization and sort by the Limit() of each RateLimiter.
	return &multiLimiter{limiters: limiters}
}

type multiLimiter struct {
	limiters []RateLimiter
}

func (l *multiLimiter) Wait(ctx context.Context) error {
	for _, l := range l.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (l *multiLimiter) Limit() rate.Limit {
	return l.limiters[0].Limit() // Because we sort the child RateLimiter instances when multiLimiter is instantiated, we can simply return the most restrictive limit, which will be the first element in the slice.
}

func Open2() *APIConnection2 {
	secondLimit := rate.NewLimiter(Per(2, time.Second), 1)   //Here we define our limit per second with no burstiness.
	minuteLimit := rate.NewLimiter(Per(10, time.Minute), 10) // Here we define our limit per minute with a burstiness of 10 to give the users their initial pool. The limit per second will ensure we don’t overload our system with requests.
	return &APIConnection2{
		rateLimiter: MultiLimiter(secondLimit, minuteLimit), //We then combine the two limits and set this as the master rate limiter for our APIConnection.
	}
}

type APIConnection2 struct {
	rateLimiter RateLimiter
}

func (a *APIConnection2) ReadFile2(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	// Pretend we do work here
	return nil
}

func (a *APIConnection2) ResolveAddress2(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	// Pretend we do work here
	return nil
}

func Open3() *APIConnection3 {
	return &APIConnection3{
		apiLimit: MultiLimiter( //Here we set up a rate limiter for API calls. There are limits for both requests per second and requests per minute.
			rate.NewLimiter(Per(2, time.Second), 2),
			rate.NewLimiter(Per(10, time.Minute), 10),
		),
		diskLimit: MultiLimiter( //Here we set up a rate limiter for disk reads. We’ll only limit this to one read per second.
			rate.NewLimiter(rate.Limit(1), 1),
		),
		networkLimit: MultiLimiter( //For networking, we’ll set up a limit of three requests per second.
			rate.NewLimiter(Per(3, time.Second), 3),
		),
	}
}

type APIConnection3 struct {
	networkLimit,
	diskLimit,
	apiLimit RateLimiter
}

func (a *APIConnection3) ReadFile3(ctx context.Context) error {
	err := MultiLimiter(a.apiLimit, a.diskLimit).Wait(ctx) //When we go to read a file, we’ll combine the limits from the API limiter and the disk limiter.
	if err != nil {
		return err
	}
	// Pretend we do work here
	return nil
}

func (a *APIConnection3) ResolveAddress3(ctx context.Context) error {
	err := MultiLimiter(a.apiLimit, a.networkLimit).Wait(ctx) //When we require network access, we’ll combine the limits from the API limiter and the network limiter.
	if err != nil {
		return err
	}
	// Pretend we do work here
	return nil
}
