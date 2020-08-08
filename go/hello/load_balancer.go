package testy

import (
	"container/heap"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
Architecture
	load balancer {pool of workers}
		maintains workers in a heap
		select statement for incoming requests
			sends incoming tasks to top of heap
			notes updates to asks as cause to re-arrange heap
	Workers {channel to receive work, channel to indicate done with work, # of pending tasks}
	requests{workFunc, response chan} sent to workers via request func
main
	workers created
	load balancer created
	loops through and calls requests
*/

const num_workers int64 = 4
var jobs_completed = 0

func main() {
	// create workers and channels
	worker_pool := PriorityQueue(make([]*Worker, num_workers))
	incoming_work := make(chan Request)
	done_work := make(chan *Worker)
	for i, _ := range worker_pool {
		worker_pool[i] = &Worker{
			name: fmt.Sprintf("Worker %v", i),
			pending: 0,
			index: i,
			request_channel: make(chan Request, 10),  // buffer channel so we can put up to 10 items in it
			done_channel: done_work,
		}
	}

	// make safe and initialize into a heap
	safe_worker_pool := SafePriorityQueue{pq: worker_pool}
	heap.Init(&safe_worker_pool.pq)
	fmt.Printf("worker pool is %v\n", worker_pool)

	//create loadbalancers, start workers, make requests, and begin to balance
	lb := &LoadBalancer{safe_worker_pool, incoming_work, done_work}
	fmt.Printf("starting workers\n")
	lb.StartWorkers()
	fmt.Printf("starting requests\n")
	for i := 0; i < 1000; i++ {
		go requester(lb.incoming_work, i)
	}
	fmt.Printf("starting balancing\n")
	start := time.Now()
	lb.Balance()
	fmt.Printf("got to end of main\n")
	time_elapsed := time.Since(start)
	fmt.Println(time_elapsed)
}

// An Worker is something we manage in a priority queue.
type Worker struct {
	name string // The name of the worker; arbitrary.
	pending int    // The number of pending items item in the queue.
	// The index is needed by updatePending and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
	request_channel chan Request
	done_channel chan *Worker
}

func (w *Worker) work() {
	for {
		request := <- w.request_channel
		request.response <- request.fn()
		w.done_channel <- w
	}
}

// SafePriorityQueue is a PriorityQueue safe to access concurrently
type SafePriorityQueue struct {
	pq PriorityQueue
	mux sync.Mutex
}


type Request struct {
	fn func() int // a function to illustrate some work to be done
	response chan int // a channel to return the result
	id int
}

func requester(work chan<- Request, i int) {
	c := make(chan int)
	fmt.Printf("requester %v preparing request\n", i)
	//time.Sleep(time.Duration(rand.Int63n(500)) * time.Millisecond)
	fmt.Printf("requester %v sending request\n", i)
	work <- Request{workFunc, c, i}
	<-c
	jobs_completed++
	fmt.Printf("requester %v received response; jobs completed now: %v\n", i, jobs_completed)
}

func workFunc() int {
	fmt.Printf("executing workFunc\n")
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(10)
}

type LoadBalancer struct {
	worker_pool SafePriorityQueue
	incoming_work chan Request
	done_notifications chan *Worker
}

func (lb *LoadBalancer) StartWorkers() {
	for _, worker := range lb.worker_pool.pq {
		go worker.work()
	}
}

func (lb *LoadBalancer) Balance() {
	for {
		timeout_time := 3 * time.Second
		timeout := time.After(timeout_time)
		fmt.Printf("balancing with timeout time of %v\n", timeout_time)
		select {
		case worker := <-lb.done_notifications:
			go lb.CompleteWork(worker)  // may need go routine to prevent deadlock
		case request := <-lb.incoming_work:
			go lb.DispatchWork(request)  // may need go routine to prevent deadlock
		case <-timeout:
			fmt.Printf("received no work for timeout %v, exiting balancing\n", timeout_time)
			return
		//default:
		//	fmt.Printf("got to end of Balance\n")
		}
	}
}

func (lb *LoadBalancer) CompleteWork (wp *Worker) {
	// lock to prevent too many go routines grabbing the heap and moving indices out of range
	lb.worker_pool.mux.Lock()
	// fix/update the heap
	lb.worker_pool.pq.updatePending(wp, wp.pending - 1)
	lb.worker_pool.mux.Unlock()
}

func (lb *LoadBalancer) DispatchWork(request Request) {
	// lock to prevent too many go routines grabbing the heap and moving indices out of range
	lb.worker_pool.mux.Lock()
	lock_flag := true
	worker := heap.Pop(&lb.worker_pool.pq).(*Worker)
	if lb.worker_pool.pq.Len() > 0 {
		lb.worker_pool.mux.Unlock()
		lock_flag = false
	}
	// Note: we cannot unlock while sending request to worker because that may lead to too many go routines Popping the heap, leaving us with lb.worker_pool.mux.Unlock()
	// send work
	worker.request_channel <- request
	// add to Worker pending count
	worker.pending++
	// push back on the heap
	if lock_flag == false {
		lb.worker_pool.mux.Lock()
	}
	heap.Push(&lb.worker_pool.pq, worker)
	lb.worker_pool.mux.Unlock()
}

