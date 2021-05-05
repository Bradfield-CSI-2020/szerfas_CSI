package main

import (
	"fmt"
	"sync"
)

type coordinator struct {
	lock   sync.RWMutex
	leader string
}

func newCoordinator(leader string) *coordinator {
	return &coordinator{
		lock:   sync.RWMutex{},
		leader: leader,
	}
}

func (c *coordinator) logStateLocked() {
	c.lock.RLock()
	defer c.lock.RUnlock()

	c.logState()  // if we keep this as fmt.Printf("leader = %q\n", c.leader) we hit a deadlock because our mutex is not re-entrant -- you this c.lock.Rlock() call on line 21 will be attempted on a mutex already locked on line 32
}

func (c *coordinator) logState() {
	fmt.Printf("leader = %q\n", c.leader)
}

func (c *coordinator) setLeader(leader string, shouldLog bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.leader = leader

	if shouldLog {
		c.logState()
	}
}

func main() {
	c := newCoordinator("us-east")
	c.logStateLocked()
	c.setLeader("us-west", true)
}
