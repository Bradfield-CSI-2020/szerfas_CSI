package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var l sync.RWMutex
	l.RLock()
	go func() {
		l.Lock()
		l.Unlock()
	}()
	//l.RLock()  // if instead do the unlock up here, then the Lock() call on line 13 won't get added to the front of the line during sleep on line 17, and fairness won't create a deadlock
	time.Sleep(1)
	l.RLock()
	l.RUnlock()
	l.RUnlock()
	fmt.Println("all good!")
}
