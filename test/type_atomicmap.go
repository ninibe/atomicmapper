// generated file - DO NOT EDIT
// command: atomicmapper -pointer -type Type


package test

import (
	"sync"
	"sync/atomic"
)

// TypeAtomicMap is a copy-on-write thread-safe map of pointers to Type
type TypeAtomicMap struct {
	mu sync.Mutex
	val atomic.Value
}

type _TypeMap map[string]*Type

// NewTypeAtomicMap returns a new initialized TypeAtomicMap
func NewTypeAtomicMap() *TypeAtomicMap {
	am := &TypeAtomicMap{}
	am.val.Store(make(_TypeMap, 0))
	return am
}

// Get returns a pointer to Type for a given key
func (am *TypeAtomicMap) Get(key string) (value *Type, ok bool) {
	value, ok = am.val.Load().(_TypeMap)[key]
	return value, ok
}

// Len returns the number of elements in the map
func (am *TypeAtomicMap) Len() int {
	return len(am.val.Load().(_TypeMap))
}

// Set inserts in the map a pointer to Type under a given key
func (am *TypeAtomicMap) Set(key string, value *Type) {
	am.mu.Lock()
	defer am.mu.Unlock()

	m1 := am.val.Load().(_TypeMap)
	m2 := make(_TypeMap, len(m1)+1)
	for k, v := range m1 {
		m2[k] = v
	}

	m2[key] = value
	am.val.Store(m2)
	return
}

// Delete removes the pointer to Type under key from the map
func (am *TypeAtomicMap) Delete(key string) {
	am.mu.Lock()
	defer am.mu.Unlock()

	m1 := am.val.Load().(_TypeMap)
	_, ok := m1[key]
	if !ok {
		return
	}

	m2 := make(_TypeMap, len(m1)-1)
	for k, v := range m1 {
		if k != key {
			m2[k] = v
		}
	}

	am.val.Store(m2)
	return
}
