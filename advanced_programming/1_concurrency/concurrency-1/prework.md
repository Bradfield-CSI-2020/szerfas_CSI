# Concurrency 1 - Goroutines and Channels

## Part 1

Please work through the following articles:

* Sections 1-8 in the "Concurrency" section of [A Tour of Go](https://tour.golang.org/concurrency/1)
* The "Concurrency" section of [Effective Go](https://golang.org/doc/effective_go.html#concurrency)
* [The Go Memory Model](https://golang.org/ref/mem)

As you read, please keep in mind the following discussion questions and prepare to discuss them in class:

`TODO: Ankify the following`

* What's meant by the quote "Do not communicate by sharing memory; instead, share memory by communicating."? What are some advantages and disadvantages of each approach?
  >Concurrent programming often involves multiple threads sharing access to memory. This can create data races that (I believe) are typically resolved via locks. Go takes an alternative approach: go routines do NOT have access to shared memory, but instead pass data around through channels. In this way, it is only by communicating that go routines can share memory -- they CANNOT communicate by using shared memory. The primary advantage of this is that is abstracts away a lot of the complexity of managing concurrent threads (e.g., locking and unlocking shared memory). **I'm not sure what the disadvantages are?**
* What's the difference between a goroutine and a "thread"? How many goroutines would you expect a typical laptop to support (order of magnitude)?
  >A goroutine is a function run concurrently with other functions in the same address space. They're I/O multiplexed onto various OS threads, which while run in the context of a single process, are scheduled by the OS kernel (whereas goroutines are scheduled by the go application). A typical laptop can support tens of thousands or hundreds of thousands of goroutines. 
* What's the difference between a buffered and an unbuffered channel, and when would you want to use each?
  >An unbuffered channel blocks on send until a receiving goroutine receives at the same time. This effectively synchronizes concurrent programs. A buffered program only blocks on sending if the buffer is full, and blocks on receiving if the buffer is empty.
* When would you want to send on a channel vs. close a channel?
  >You only need to close on a channel when the receiving goroutine needs to know the channel is closed. This is true when using the `for i := range some_channel` for-loop syntax which pulls from a channel until it retrieves indication that it has been closed. This allows buffered channels to synchronize communication in a different way - a bit like a semaphore limiting throughput.
* What do `select`, `sync.Once`, and `sync.WaitGroup` do? What are some examples of when you would use each of these concurrency primitives?
  >`select` functions a bit like a `switch` statement by passing control to logic that is nonblocked by a channel communication statement (send or receive). For example, you a goroutine may function as a worker that can handle different jobs. A `select` statement will set the worker to work on any job that comes off any of its assigned queues.
  > 
* How does the idea of "happens before" help us reason about correctness in concurrent programs? What's a situation where neither "A happens before B" nor "B happens before A"?

* Questions for class
> * What are the disadvantages of goroutines over threads?
> * How do goroutines have access to different CPU cores if they're not scheduled by the OS? Is Go multiplexing these on to different threads, and indirectly CPU cores, in the background?
> * Example 5 below: how do we guarantee non-sequential process?
> * Example 5 below: why wouldn't an unbuffered channel work for parallelism? Since a receive on the channel completes before a send finishes, doesn't that actually guarantee that the first response is passes control back to the original routine?

## Part 2

Each of the included code snippets (`example-?.go`) contains one or more concurrency errors. What are the errors, and how would you fix them?
> Example 1:
> There are a couple problems with this example: First, the `i` within the goroutine function literal is a closure, and so all goroutines will share a loop variable.
> Second, the loop will proceed through launching all goroutines before any of them receive control and run the statement "launched goroutine `i`". If by `launch` we mean the moment in which the goroutine is created, then this statement will run at a very different point in time from what it means to communicate.
> To fix, you could pass `i` into the go routine or redeclare the variable (e.g., `i := i`) prior to line 11.
> 
> Example 2:
> line 10 declares a channel variable with `var done chan struct{}` but a channel must be created before use using the `make` keyword. I believe sending and receiving on a nil channel block forever, so after all of the goroutines are created by the loop starting on line 11, the first `<-done` reached by the program on line 22 will block the initial routine, then all three other goroutines will proceed until `done <- struct{}{}` where they will block forever. 
> After the original routine and three new routines are blocked, a deadlock error is thrown. 
> To fix, we simply need to change line 10 to `done := make(chan struct{})`
> 
> Example 3: `wg.Add(1)` needs to be moved from line 18 to before line 17. Otherwise line 32 `wg.Wait()` might be called before anything is added to wg.Add.
> 
> Example 4: Here we create a buffered channel which will no block on line 14. So we'll the original routine will carry through to line 15 and then exit before blocking and passing control to the goroutine created on line 9. 
> To fix, we can make this an unbuffered channel, changing `done := make(chan struct{}, 1)` to `done := make(chan struct{})` on line 8.
> 
> Example 5: I *think* the problem is: our `parallelQuery` function spins up several goroutines, then blocks hoping to unblock as soon as one returns, but does not ensure that it will unblock as soon as one returns.
> It is possible that instead of control being returned to the original routine, it instead is passed to any of the other routines, even all of them, before returned to the original routine, eliminating the benefit we hope to see via parallelization.
> I'm not sure how to fix yet... we might be able to restructure to so that each query is passed an unbuffered channel, then we use a select statement to block until one of those are channels returns. If we use an unbuffered channel then the send will block until the receive completes... but that may still one the risk of control being passed to each query sequentially...
## Part 3

Finally, please complete Exercise 8.8 from The Go Programming Language:

> Using a select statement, add a timeout to the echo server from Section 8.3 so that it disconnects any client that shouts nothing within 10 seconds.

Note that code from The Go Programming Language can be found on [Github](https://github.com/adonovan/gopl.io/a).

## Additional Resources

If you would like more practice, Chapter 8 of The Go Programming Language contains many more examples and exercises.

[Time, Clocks, and the Ordering of Events in a Distributed System](https://lamport.azurewebsites.net/pubs/time-clocks.pdf) discusses the "happens before" relationship in the context of distributed systems. Internalizing this logical (rather than physical) notion of "time" will help you reason about distributed systems much more effectively.

[Communicating Sequential Processes](https://www.cs.cmu.edu/~crary/819-f09/Hoare78.pdf) is the original source of many concurrency ideas used in Go (channels, goroutines, select).