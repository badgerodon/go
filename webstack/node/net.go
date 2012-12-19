package main

import (
	"sync"
)

var (
	_nextPort int
	netLock sync.Mutex
)

func nextPort() int {
	netLock.Lock()
	defer netLock.Unlock()

	_nextPort += 1
	return 49152 + (_nextPort % (65535 - 49152))
}