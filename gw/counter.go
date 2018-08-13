package gw

import (
	"fmt"
	"sync/atomic"
)

type Counter int64

// Get returns the value of int64 atomically.
func (a *Counter) Get() int64 {
	return int64(*a)
}

// Set sets the value of int64 atomically.
func (a *Counter) Set(newValue int64) {
	atomic.StoreInt64((*int64)(a), newValue)
}

// GetAndSet sets new value and returns the old atomically.
func (a *Counter) GetAndSet(newValue int64) int64 {
	for {
		current := a.Get()
		if a.CompareAndSet(current, newValue) {
			return current
		}
	}
}

// CompareAndSet compares int64 with expected value, if equals as expected
// then sets the updated value, this operation performs atomically.
func (a *Counter) CompareAndSet(expect, update int64) bool {
	return atomic.CompareAndSwapInt64((*int64)(a), expect, update)
}

// GetAndIncrement gets the old value and then increment by 1, this operation
// performs atomically.
func (a *Counter) GetAndIncrement() int64 {
	for {
		current := a.Get()
		next := current + 1
		if a.CompareAndSet(current, next) {
			return current
		}
	}

}

// GetAndDecrement gets the old value and then decrement by 1, this operation
// performs atomically.
func (a *Counter) GetAndDecrement() int64 {
	for {
		current := a.Get()
		next := current - 1
		if a.CompareAndSet(current, next) {
			return current
		}
	}
}

// GetAndAdd gets the old value and then add by delta, this operation
// performs atomically.
func (a *Counter) GetAndAdd(delta int64) int64 {
	for {
		current := a.Get()
		next := current + delta
		if a.CompareAndSet(current, next) {
			return current
		}
	}
}

// IncrementAndGet increments the value by 1 and then gets the value, this
// operation performs atomically.
func (a *Counter) IncrementAndGet() int64 {
	for {
		current := a.Get()
		next := current + 1
		if a.CompareAndSet(current, next) {
			return next
		}
	}
}

// DecrementAndGet decrements the value by 1 and then gets the value, this
// operation performs atomically.
func (a *Counter) DecrementAndGet() int64 {
	for {
		current := a.Get()
		next := current - 1
		if a.CompareAndSet(current, next) {
			return next
		}
	}
}

// AddAndGet adds the value by delta and then gets the value, this operation
// performs atomically.
func (a *Counter) AddAndGet(delta int64) int64 {
	for {
		current := a.Get()
		next := current + delta
		if a.CompareAndSet(current, next) {
			return next
		}
	}
}

func (a *Counter) String() string {
	return fmt.Sprintf("%d", a.Get())
}
