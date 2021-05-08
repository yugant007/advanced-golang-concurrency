package main

import (
	"runtime/pprof"
)

func main() {
	//Anatomy of a Goroutine Error
	//For example, when this simple program is executed:
	//waitForever := make(chan interface{})
	//go func() {
	//	panic("test panic")
	//}()
	//<-waitForever
	//panic: test panic
	//
	//  goroutine 4 [running]:
	//  main.main.func1() //3
	//      /tmp/babel-3271QbD/go-src-32713Rn.go:6 +0x65 //1
	//  created by main.main
	//      /tmp/babel-3271QbD/go-src-32713Rn.go:7 +0x4e //2
	//  exit status 2

	//1. Refers to where the panic occurred.
	//2. Refers to where the goroutine was started.
	//3. Indicates the name of the function running as a goroutine. If it’s an anonymous function as in this example,
	//an automatic and unique identifier is assigned.
	//If you’d like to see the stack traces of all the goroutines that were executing
	//when the program panicked, you can enable the old behavior by setting the GOTRACEBACK
	//environmental variable to all.

	//Race Detection
	//a -race flag was added as a flag for most go commands:
	//	$ go test -race mypkg    # test the package
	//  $ go run -race mysrc.go  # compile and run the program
	//  $ go build -race mycmd   # build the command
	//  $ go install -race mypkg # install the package
	//If you’re a developer and all you need is a more reliable way to detect race conditions, this is
	//really all you need to know. One caveat of using the race detector is that the algorithm will
	//only find races that are contained in code that is exercised. For this reason, the Go team
	//recommends running a build of your application built with the race flag under real-world load.
	//This increases the probability of finding races by virtue of increasing the probability that
	//more code is exercised.
	//There are also some options you can specify via environmental variables to tweak the
	//behavior of the race detector, although generally the defaults are sufficient:
	//LOG_PATH
	//This tells the race detector to write reports to the LOG_PATH.pid file. You can also
	//pass it special values: stdout and stderr. The default value is stderr.
	//STRIP_PATH_PREFIX
	//This tells the race detector to strip the beginnings of file paths in reports to make them more concise.
	//HISTORY_SIZE
	//This sets the per-goroutine history size, which controls how many previous memory accesses are remembered
	//per goroutine. The valid range of values is [0, 7]. The memory allocated for goroutine history begins at
	//32 KB when HISTORY_SIZE is 0, and doubles with each subsequent value for a maximum of 4 MB at a HISTORY_SIZE
	//of 7. When you see “failed to restore the stack” in reports, that’s an indicator to increase this value; however
	//, it can significantly increase memory consumption.
	//Given this simple program we first looked at
	//var data int
	//go func() { //1
	//	data++
	//}()
	//if data == 0 {
	//	fmt.Printf("the value is %v.\n", data)
	//}
	// ==================
	//  WARNING: DATA RACE
	//  Write by goroutine 6:
	//    main.main.func1()
	//        /tmp/babel-10285ejY/go-src-10285GUP.go:6 +0x44 //1
	//
	//  Previous read by main goroutine:
	//    main.main()
	//        /tmp/babel-10285ejY/go-src-10285GUP.go:7 +0x8e //2
	//
	//  Goroutine 6 (running) created at:
	//    main.main()
	//        /tmp/babel-10285ejY/go-src-10285GUP.go:6 +0x80
	//  ==================
	//  Found 1 data race(s)
	//  exit status 66
	// 1. Signifies a goroutine that is attempting to write unsynchronized memory access.
	// 2. Signifies a goroutine (in this case the main goroutine) trying to read this same memory.

	//pprof
	//In large codebases, it can sometimes be difficult to ascertain how your program is performing at runtime.
	//many goroutines are running? Are your CPUs being fully utilized? How’s memory usage doing? Profiling is a
	//great way to answer these questions, and Go has a package in the standard library to support a profiler
	//named “pprof.”
	//pprof is a tool that was created at Google and can display profile data either while a program is running,
	//or by consuming saved runtime statistics. The usage of the program is pretty well described by its help flag,
	//so instead we’ll stick to discussing the runtime/pprof package here—specifically as it pertains to concurrency.
	//The runtime/pprof package is pretty simple, and has predefined profiles to hook into and display:
	// goroutine    - stack traces of all current goroutines
	//  heap         - a sampling of all heap allocations
	//  threadcreate - stack traces that led to the creation of new OS threads
	//  block        - stack traces that led to blocking on synchronization primitives
	//  mutex        - stack traces of holders of contended mutexes
	//From the context of concurrency, most of these are useful for understanding what’s happening
	//within your running program. For example, here’s a goroutine that can help you detect goroutine leaks:
	//log.SetFlags(log.Ltime | log.LUTC)
	//log.SetOutput(os.Stdout)
	//
	//// Every second, log how many goroutines are currently running.
	//go func() {
	//	goroutines := pprof.Lookup("goroutine")
	//	for range time.Tick(1*time.Second) {
	//		log.Printf("goroutine count: %d\n", goroutines.Count())
	//	}
	//}()
	//
	//// Create some goroutines which will never exit.
	//var blockForever chan struct{}
	//for i := 0; i < 10; i++ {
	//	go func() { <-blockForever }()
	//	time.Sleep(500*time.Millisecond)
	//}
	//These built-in profiles can really help you profile and diagnose issues with your program, but of course
	//you can write custom profiles tailored to help you monitor your programs:
	//prof := newProfIfNotDef("my_package_namespace")
}

func newProfIfNotDef(name string) *pprof.Profile {
	prof := pprof.Lookup(name)
	if prof == nil {
		prof = pprof.NewProfile(name)
	}
	return prof
}

