package main

import (
	"container/heap"
	"fmt"
)

// A PriorityQueue implements heap.Interface and holds Workers.
type PriorityQueue []*Worker

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
	item := x.(*Worker)
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

// updatePending modifies the pending items and name of a Worker in the queue.
func (pq *PriorityQueue) updatePending(item *Worker, pending int) {
	item.pending = pending
	heap.Fix(pq, item.index)
}

// This example creates a PriorityQueue with some items, adds and manipulates an item,
// and then removes the items in priority order.
func main2() {

	num_workers := 4
	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(PriorityQueue, num_workers)
	for i := 0; i < num_workers; i++ {
		pq[i] = &Worker{
			index:    i,
		}
	}
	heap.Init(&pq)

	// Insert a new item and then modify its pending item count.
	item := &Worker{
		name: "Worker5",
		pending: 1,
	}
	heap.Push(&pq, item)
	pq.updatePending(item, 5)

	// Take the items out; they arrive in increasing priority order.
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Worker)
		fmt.Printf("%d items pending for worker: %s\n", item.pending, item.name)
	}
	// heap.Pop(&pq)    <---- this returns an error, illustrating this heap is not safe to too many pops
}
