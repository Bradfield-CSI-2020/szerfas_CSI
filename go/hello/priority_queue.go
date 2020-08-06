// This example demonstrates a priority queue built using the heap interface.
package main

import (
	"container/heap"
	"sync"
)

// An Worker2 is something we manage in a priority queue.
type Worker2 struct {
	name string // The name of the worker; arbitrary.
	pending int    // The number of pending items item in the queue.
	// The index is needed by updatePending and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
	request_channel chan Request
	done_channel chan *Worker2
}

func (w *Worker2) work() {
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

// A PriorityQueue implements heap.Interface and holds Worker2s.
type PriorityQueue []*Worker2

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the worker with the lowest number pending so we use less than here.
	return pq[i].pending < pq[j].pending
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Worker2)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// updatePending modifies the pending items and name of a Worker2 in the queue.
func (pq *PriorityQueue) updatePending(item *Worker2, pending int) {
	item.pending = pending
	heap.Fix(pq, item.index)
}

// This example creates a PriorityQueue with some items, adds and manipulates an item,
// and then removes the items in priority order.
/*func main() {

	num_workers := 4
	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(PriorityQueue, num_workers)
	for i := 0; i < num_workers; i++ {
		pq[i] = &Worker2{
			index:    i,
		}
	}
	heap.Init(&pq)

	// Insert a new item and then modify its pending item count.
	item := &Worker2{
		name: "Worker5",
		pending: 1,
	}
	heap.Push(&pq, item)
	pq.updatePending(item, 5)

	// Take the items out; they arrive in increasing priority order.
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Worker2)
		fmt.Printf("%d items pending for worker: %s\n", item.pending, item.name)
	}
}*/
