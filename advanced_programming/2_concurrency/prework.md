# Concurrency 2 - Shared Memory

## Part 1

Please read Chapter 9 of _The Go Programming Language_. As you read, please keep in mind the following discussion questions and prepare to discuss them in class:

* What's a "mutex", and why is it useful for writing concurrent programs?
* What's a "read/write mutex", and when would you want to use it instead of a (regular) mutex?
	* What are some performance tradeoffs to consider between `Mutex` and `RWMutex`?
* What does it mean for a mutex to be "re-entrant"?
	* What are some arguments for whether or not this is a good idea?
* What sort of errors can be caught by Go's "race detector"?
	* Can you implement a simple program where the race detector finds an error?

## Part 2

Each of the included code snippets contains one or more concurrency errors. What are the errors, and how would you fix them?

_Note: `example-4-optional` is a "stretch goal" because it's a bit more involved than the other examples. However, since it's (loosely) inspired by an actual production deadlock, you might nonetheless enjoy hunting down the issue._

## Part 3

Consider the following interface for an "ID service":

```go
type idService interface {
	// Returns an ID that hasn't been returned by a
	// previous call to the same idService.
	getUniqueId() uint64
}
```

Implement this interface three times using the following different strategies:

* Use a `sync.Mutex` to coordinate access to a shared counter
* Use a separate goroutine that has exclusive access to a private counter
* Use an atomic variable to coordinate access to a shared counter

How do you expect the three implementations to compare in terms of performance? What are the bottlenecks in each case?

_Note: If you don't have to return IDs in some "globally ascending order", one way to improve throughput is to have `idService` return a range of IDs rather than a single ID; callers would then be responsible for handing out individual IDs from among the most recent range they received._

## Stretch Goals

* You can think of a "semaphore" as generalizing a mutex from 2 states (unlocked / locked) to many states (0, 1, 2, ..., N "tickets" available). Semaphores are useful for coordinating concurrent access to a fixed pool of resources (e.g. N database connections). Since Go's standard library doesn't have a semaphore built in, _how can you simulate the behavior of a semaphore using Go's existing concurrency primitives_?

* Can you think of a way to implement a `RWMutex` yourself, using Go's other concurrency primitives as building blocks? (Hint: in addition to `sync.Mutex`, you might also find `sync.Cond` useful.)

## Additional Resources

[Russ Cox posted to the golang-nuts mailing list](https://groups.google.com/g/golang-nuts/c/XqW1qcuZgKg/m/Ui3nQkeLV80J) to describe why he thinks re-entrant mutexes are a bad idea.

The documentation for [`sync`](https://golang.org/pkg/sync/) and [`sync/atomic`](https://golang.org/pkg/sync/atomic/) provides more details about available concurrency primitives, as well as links to source code (if you're curious how they work under the hood).

For a fun way to practice concurrent programming, play [The Deadlock Empire](https://deadlockempire.github.io/) ("Slay dragons, master concurrency!")

[The Little Book of Semaphores](http://alumni.cs.ucr.edu/~kishore/papers/semaphores.pdf) contains a lot of problems and solutions for additional practice reasoning about concurrency primitives.
