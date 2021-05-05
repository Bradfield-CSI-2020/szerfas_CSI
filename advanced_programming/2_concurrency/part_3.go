package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)


type idService interface {
	// Returns an ID that hasn't been returned by a
	// previous call to the same idService.
	getUniqueId() uint64
}

/*
Implement this interface three times using the following different strategies:

* Use a `sync.Mutex` to coordinate access to a shared counter
* Use a separate goroutine that has exclusive access to a private counter
* Use an atomic variable to coordinate access to a shared counter

How do you expect the three implementations to compare in terms of performance? What are the bottlenecks in each case?

_Note: If you don't have to return IDs in some "globally ascending order", one way to improve throughput is to have `idService`
return a range of IDs rather than a single ID; callers would then be responsible for handing out individual IDs from among
the most recent range they received._
*/

type IdMutex struct{
	count int
	mu sync.Mutex
}

func (id *IdMutex) getUniqueId() int {
	id.mu.Lock()
	defer id.mu.Unlock()
	id.count++
	return id.count
}

type IdMonitorRoutine struct{
	request chan bool
	response chan int
}

func (id *IdMonitorRoutine) getUniqueId() int {
	id.request <- true
	return <- id.response
}

// must be started as separate go routine
func (id *IdMonitorRoutine) startCount() int {
	count := 0
	for {
		<- id.request
		count++
		id.response <- count
	}
}


type IdAtomic struct{
	count int64
}

func (id *IdAtomic) getUniqueId() int64 {
	return atomic.AddInt64(&id.count, 1)
}

func main() {
	// call using mutex
	var (
		mutex IdMutex
		atom IdAtomic
	)
	// call using confinement to a monitor routine
	monitor := IdMonitorRoutine{make(chan bool), make(chan int)}
	go monitor.startCount()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {fmt.Printf("mutex id: %d\n", mutex.getUniqueId()); wg.Done()}()
		wg.Add(1)
		go func() {fmt.Printf("monitor id: %d\n", monitor.getUniqueId()); wg.Done()}()
		wg.Add(1)
		go func() {fmt.Printf("atomic id: %d\n", atom.getUniqueId()); wg.Done()}()
	}
	wg.Wait()
}
