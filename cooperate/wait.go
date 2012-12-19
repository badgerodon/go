package main

import (
	"sync"
)

type (
	Waiter struct {
		signal chan bool
		lock sync.Mutex
		cnt int
	}
)

func NewWaiter() *Waiter {
	return &Waiter{
		signal: make(chan bool, 1),
	}
}

func (w Waiter) Inc() {
	w.lock.Lock()
	w.cnt++
	w.lock.Unlock()
}

func (w Waiter) Dec() {
	w.lock.Lock()
	w.cnt--
	if w.cnt == 0 {
		w.signal <- true
	}
	w.lock.Unlock()
}

func (w Waiter) Wait() {
	<- w.signal
}