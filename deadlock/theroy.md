A race condition occurs when two or more operations must execute in the correct order,
but the program has not been written so that this order is guaranteed to be maintained.

Most of the time, this shows up in what’s called a data race, where one concurrent
operation attempts to read a variable while at some undetermined time another concurrent
operation is attempting to write to the same variable.

Here’s a basic example:
var data int
go func() {
    data++ // line 3
}()
if data == 0 { // line 5
     fmt.Printf("the value is %v.\n", data)
}
Here, lines 3 and 5 are both trying to access the variable data, but there is no guarantee
what order this might happen. There are three possible outcomes to running this code:
1. Nothing is printed. In this case, line 3 was executed before line 5.
2. “the value is 0” is printed. In this case, lines 5 and 6 were executed before line 3.
3. “the value is 1” is printed. In this case, line 5 was executed before line 3, but line 3 was executed before line 6.


                                        ----------------------
                                              ATOMICITY
                                        ----------------------

When something is considered atomic, or to have the property of atomicity, this means that
within the context that it is operating, it is indivisible, or uninterruptible.

The first thing that’s very important is the word “context.” Something may be atomic
in one context, but not another. Operations that are atomic within the context of your
process may not be atomic in the context of the operating system; operations that are
atomic within the context of the operating system may not be atomic within the context
of your machine; and operations that are atomic within the context of your machine may
not be atomic within the context of your application. In other words, the atomicity of
an operation can change depending on the currently defined scope. When thinking about
an atomicity, very often the first thing you need to do is to define the context, or scope,
the operation will be considered to be atomic in.

Now let’s look at the terms “indivisible” and “uninterruptible.” These terms mean that
within the context you’ve defined, something that is atomic will happen in its entirety
without anything happening in that context simultaneously.
i++
It may look atomic, but a brief analysis reveals several operations:
1. Retrieve the value of i.
2. Increment the value of i.
3. Store the value of i.
While each of these operations alone is atomic, the combination of the three may not be,
depending on your context. This reveals an interesting property of atomic operations:
combining them does not necessarily produce a larger atomic operation. Making the operation
atomic is dependent on which context you’d like it to be atomic within. If your context is
a program with no concurrent processes, then this code is atomic within that context. If
your context is a goroutine that does not expose i to other goroutines, then this code is atomic.

So why do we care? Atomicity is important because if something is atomic, implicitly it is
safe within concurrent contexts. This allows us to compose logically correct programs.

                                   -----------------------------------------
                                          MEMORY ACCESS SYNCHRONIZATION
                                   -----------------------------------------

Let’s say we have a data race: two concurrent processes are attempting to access the same area
of memory, and the way they are accessing the memory is not atomic.

var data int
go func() { data++}()
if data == 0 {
    fmt.Println("the value is 0.")
} else {
    fmt.Printf("the value is %v.\n", data)
}

In fact, there’s a name for a section of your program that needs exclusive access to a shared resource.
This is called a critical section. In this example, we have three critical sections:
1. Our goroutine, which is incrementing the data variables.
2. Our if statement, which checks whether the value of data is 0.
3. Our fmt.Printf statement, which retrieves the value of data for output.

There are various ways to guard your program’s critical sections, and Go has some better ideas
on how to deal with this, but one way to solve this problem is to synchronize access to the
memory between your critical sections.
The following code is not idiomatic Go (and I don’t suggest you attempt to solve your data race
problems like this), but it very simply demonstrates memory access synchronization.

var memoryAccess sync.Mutex //1
var value int
go func() {
    memoryAccess.Lock() //2
    value++
    memoryAccess.Unlock() //3
}()

memoryAccess.Lock() //4
if value == 0 {
    fmt.Printf("the value is %v.\n", value)
} else {
    fmt.Printf("the value is %v.\n", value)
}
memoryAccess.Unlock() //5

1. Here we add a variable that will allow our code to synchronize access to the data variable’s memory.
2. Here we declare that until we declare otherwise, our goroutine should have exclusive access to this memory.
3. Here we declare that the goroutine is done with this memory.
4. Here we once again declare that the following conditional statements should have exclusive 
   access to the data variable’s memory.
5. Here we declare we’re once again done with this memory.

You may have noticed that while we have solved our data race, we haven’t actually solved our 
race condition! The order of operations in this program is still nondeterministic; we’ve just 
narrowed the scope of the nondeterminism a bit. In this example, either the goroutine will 
execute first, or both our if and else blocks will. We still don’t know which will occur first 
in any given execution of this program.
It is true that you can solve some problems by synchronizing access to the memory, 
but as we just saw, it does not automatically solve logical correctness. 
Further, it can also create maintenance and performance problems.

Synchronizing access to the memory in this manner also has performance ramifications.
But the calls to Lock you see can make our program slow. Every time we perform one of these operations,
our program pauses for a period of time. This brings up two questions:
1. Are my critical sections entered and exited repeatedly?
2. What size should my critical sections be?

The previous sections have all been about discussing program correctness in that if these issues 
are managed correctly, your program will never give an incorrect answer. Unfortunately, even 
if you successfully handle these classes of issues, there is another class of issues to contend 
with: deadlocks, livelocks, and starvation. These issues all concern ensuring your program has 
something useful to do at all times. If not handled properly, your program could enter a state 
in which it will stop functioning altogether.

                                            ----------------------
                                                  DEADLOCK
                                            ----------------------

A deadlocked program is one in which all concurrent processes are waiting on one another. 
In this state, the program will never recover without outside intervention.

type value struct {
    mu    sync.Mutex
    value int
}
var wg sync.WaitGroup
var printSum = func(v1, v2 *value) {
    defer wg.Done()
    v1.mu.Lock() // 1
    defer v1.mu.Unlock() // 2

    time.Sleep(2*time.Second) // 3
    v2.mu.Lock()
    defer v2.mu.Unlock()

    fmt.Printf("sum=%v\n", v1.value + v2.value)
}
var a, b value
wg.Add(2)
go printSum(&a, &b)
go printSum(&b, &a)
wg.Wait()

1. Here we attempt to enter the critical section for the incoming value.
2. Here we use the defer statement to exit the critical section before printSum returns.
3. Here we sleep for a period of time to simulate work (and trigger a deadlock).

If you were to try and run this code, you’d probably see:
fatal error: all goroutines are asleep - deadlock!

Essentially, we have created two gears that cannot turn together: our first call to printSum
locks a and then attempts to lock b, but in the meantime our second call to printSum has
locked b and has attempted to lock a. Both goroutines wait infinitely on each other.

It turns out there are a few conditions that must be present for deadlocks to arise. The conditions are
now known as the Coffman Conditions and are the basis for techniques that help detect, prevent, and correct deadlocks.
The Coffman Conditions are as follows:
1. Mutual Exclusion - A concurrent process holds exclusive rights to a resource at any one time.
2. Wait For Condition -A concurrent process must simultaneously hold a resource and be waiting for
an additional resource.
3. No Preemption - A resource held by a concurrent process can only be released by that process,
so it fulfills this condition.
4. Circular Wait - A concurrent process (P1) must be waiting on a chain of other concurrent
processes (P2), which are in turn waiting on it (P1), so it fulfills this final condition too.

These laws allow us to prevent deadlocks too. If we ensure that at least one of these conditions
is not true, we can prevent deadlocks from occurring.