package memoize

import (
	"reflect"
	"sync"
	"sync/atomic"
)

// An EvictionPolicy controls the eviction behavior of keys in a Store when
// they no longer have any references.
type EvictionPolicy int

const (
	// ImmediatelyEvict evicts keys as soon as they no longer have references.
	ImmediatelyEvict EvictionPolicy = iota

	// NeverEvict does not evict keys.
	NeverEvict
)

type Store struct {
	evictionPolicy EvictionPolicy

	promisesMu sync.Mutex
	promises   map[interface{}]*Promise
}

// Promise returns a reference-counted promise for the future result of
// calling the specified function.
//
// Calls to Promise with the same key return the same promise, incrementing its
// reference count.  The caller must call the returned function to decrement
// the promise's reference count when it is no longer needed. The returned
// function must not be called more than once.
//
// Once the last reference has been released, the promise is removed from the
// store.
func (store *Store) Promise(key interface{}, function Function) (*Promise, func()) {
	store.promisesMu.Lock()
	p, ok := store.promises[key]
	if !ok {
		p = NewPromise(reflect.TypeOf(key).String(), function)
		if store.promises == nil {
			store.promises = map[interface{}]*Promise{}
		}
		store.promises[key] = p
	}
	p.refcount++
	store.promisesMu.Unlock()

	var released int32
	release := func() {
		if !atomic.CompareAndSwapInt32(&released, 0, 1) {
			panic("release called more than once")
		}
		store.promisesMu.Lock()

		p.refcount--
		if p.refcount == 0 && store.evictionPolicy != NeverEvict {
			// Inv: if p.refcount > 0, then store.promises[key] == p.
			delete(store.promises, key)
		}
		store.promisesMu.Unlock()
	}

	return p, release
}

// Stats returns the number of each type of key in the store.
func (s *Store) Stats() map[reflect.Type]int {
	result := map[reflect.Type]int{}

	s.promisesMu.Lock()
	defer s.promisesMu.Unlock()

	for k := range s.promises {
		result[reflect.TypeOf(k)]++
	}
	return result
}

// DebugOnlyIterate iterates through the store and, for each completed
// promise, calls f(k, v) for the map key k and function result v.  It
// should only be used for debugging purposes.
func (s *Store) DebugOnlyIterate(f func(k, v interface{})) {
	s.promisesMu.Lock()
	defer s.promisesMu.Unlock()

	for k, p := range s.promises {
		if v := p.Cached(); v != nil {
			f(k, v)
		}
	}
}
