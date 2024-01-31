package lock

import (
	"errors"
	"sync"
	"time"
)

// LockMap structure stores the lock state for each key in the sync.Map
type LockMap struct {
	m              sync.Map // Stores the lock status
	LockWaitSecond int
}

// NewLockMap creates a new LockMap
func NewLockMap(lockWaitSecond int) *LockMap {
	return &LockMap{
		m:              sync.Map{},
		LockWaitSecond: lockWaitSecond,
	}
}

// Lock attempts to lock a specific key within the given timeout duration.
// If the key is already locked, it waits until it is released or the timeout is reached.
func (lm *LockMap) Lock(key interface{}) error {
	var mu *sync.Mutex
	for {
		actual, loaded := lm.m.LoadOrStore(key, &sync.Mutex{})
		mu = actual.(*sync.Mutex)

		if !loaded {
			mu.Lock()
			return nil
		}

		// Use a select statement to wait for either the lock to be available
		// or for the timeout to expire.
		select {
		case <-time.After(time.Duration(lm.LockWaitSecond)):
			// Timeout reached, return an error
			return errors.New("lock timeout")
		default:
			// Attempt to lock
			mu.Lock()
			return nil
		}
	}
}

// Unlock releases the lock for a specific key
func (lm *LockMap) Unlock(key interface{}) {
	if actual, ok := lm.m.Load(key); ok {
		mu := actual.(*sync.Mutex)
		mu.Unlock()
	}
}
