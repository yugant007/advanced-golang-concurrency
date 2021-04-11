STARVATION

Starvation is any situation where a concurrent process cannot get all the resources it needs
to perform work.

When we discussed livelocks, the resource each goroutine was starved of was a shared lock. Livelocks
warrant discussion separate from starvation because in a livelock, all the concurrent processes are
starved equally, and no work is accomplished. More broadly, starvation usually implies that there are
one or more greedy concurrent process that are unfairly preventing one or more concurrent processes
from accomplishing work as efficiently as possible, or maybe at all.


FINDING A BALANCE
It is worth mentioning that the previous code example can also serve as an example of the
performance ramifications of memory access synchronization. Because synchronizing access to
the memory is expensive, it might be advantageous to broaden our lock beyond our critical sections.
On the other hand, by doing so—as we saw—we run the risk of starving other concurrent processes.

If you utilize memory access synchronization, you’ll have to find a balance between preferring
coarse-grained synchronization for performance, and fine-grained synchronization for fairness.
When it comes time to performance tune your application, to start with, I highly recommend you
constrain memory access synchronization only to critical sections; if the synchronization becomes
a performance problem, you can always broaden the scope. It’s much harder to go the other way.

We should also consider the case where the starvation is coming from outside the Go process. Keep
in mind that starvation can also apply to CPU, memory, file handles, database connections:
any resource that must be shared is a candidate for starvation.