package main

import (
	"fmt"
	"math/rand"
	"time"
)

const num_workers int64 = 4

type Request struct {
	fn func() int // a function to illustrate some work to be done
	response chan int // a channel to return the result
	id int
}

func requester(work chan<- Request, i int) {
	c := make(chan int)
	fmt.Printf("requester %v preparing request\n", i)
	time.Sleep(time.Duration(rand.Int63n(500)) * time.Millisecond)
	fmt.Printf("requester %v sending request\n", i)
	work <- Request{workFunc, c, i}
	<-c
	fmt.Printf("requester %v received response\n", i)
}

func workFunc() int {
	fmt.Printf("executing workFunc\n")
	time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(10)
}

type Worker struct {
	work <-chan Request

	done chan bool
}

func worker(work <-chan Request, i int) {
	for {
		request := <-work
		fmt.Printf("worker %v received work with request_id %v\n", i, request.id)
		request.response <- request.fn()
		fmt.Printf("worker %v sent response for request_id %v\n", i, request.id)
	}
}

/*func main() {
	fmt.Println("Hello, world.")
	fmt.Println(cmp.Diff("Hello World", "Hello Go"))
	work_channel := make(chan Request)
	for i := 0; i < 5; i++ {
		go worker(work_channel, i)
	}
	fmt.Println("Workers created\n")
	for i := 0; i < 50; i++ {
		go requester(work_channel, i)
	}
	fmt.Print("Requesters created\n")
	timeout := time.After(10 * time.Second)
	<-timeout
}*/

/*Architecture
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
	loops through and calls requests*/

func main() {
	worker_pool := PriorityQueue(make([]*Worker2, num_workers))
	incoming_work := make(chan Request)
	done_work := make(chan *Worker2)
	for i, _ := range worker_pool {
		worker_pool[i] = &Worker2{
			name: fmt.Sprintf("Worker %v", i),
			pending: 0,
			index: i,
			request_channel: make(chan Request, 10),  // buffer channel so we can put up to 10 items in it
			done_channel: done_work,
		}
	}
	worker_pool = PriorityQueue(worker_pool)
	fmt.Printf("worker pool is %v\n", worker_pool)
	lb := &LoadBalancer{SafePriorityQueue{pq: worker_pool}, incoming_work, done_work}
	fmt.Printf("starting workers\n")
	lb.StartWorkers()
	fmt.Printf("starting requests\n")
	for i := 0; i < 100; i++ {
		go requester(lb.incoming_work, i)
	}
	fmt.Printf("starting balancing\n")
	lb.Balance()
	fmt.Printf("got to end of main\n")
}

type Worker1 struct {
	request_channel chan Request
	pending int
	done_channel chan *Worker1
}

func (w *Worker2) DoWork() {
	for {
		request := <-w.request_channel
		request.response <-request.fn()
	}
}

type LoadBalancer struct {
	worker_pool SafePriorityQueue
	incoming_work chan Request
	done_notifications chan *Worker2
}

func (lb *LoadBalancer) StartWorkers() {
	for _, worker := range lb.worker_pool.pq {
		go worker.work()
	}
}

func (lb *LoadBalancer) Balance() {
	for {
		timeout_time := 1 * time.Second
		timeout := time.After(timeout_time)
		fmt.Printf("balancing with timeout time of %v\n", timeout_time)
		select {
		case worker := <-lb.done_notifications:
			lb.CompleteWork(worker)  // may need go routine to prevent deadlock
		case request := <-lb.incoming_work:
			go lb.DispatchWork(request)  // may need go routine to prevent deadlock
		case <-timeout:
			fmt.Printf("received no work for timeout %v, exiting balancing\n", timeout_time)
			return
		//default:
		//	fmt.Printf("got to ned of Balance\n")
		}
	}
}

func (lb *LoadBalancer) CompleteWork (wp *Worker2) {
	// reduce Worker2 priority
	wp.pending--
	// fix/update the heap
	// lock to prevent simultaneous pops and replacement
	lb.worker_pool.mux.Lock()
	lb.worker_pool.pq.updatePending(wp, wp.pending)
	lb.worker_pool.mux.Unlock()
}

func (lb *LoadBalancer) DispatchWork(request Request) {
	// lock to prevent simultaneous pops and replacement
	lb.worker_pool.mux.Lock()
	// pop from the heap
	worker_interface := lb.worker_pool.pq.Pop()
	lb.worker_pool.mux.Unlock()
	worker := worker_interface.(*Worker2)
	// send work
	worker.request_channel <- request
	// add to Worker2 pending count
	worker.pending++
	// push back on the heap
	lb.worker_pool.mux.Lock()
	lb.worker_pool.pq.Push(worker)
	lb.worker_pool.mux.Unlock()
}

