package test

import (
	"crypto/rand"
	"sync"
	"testing"
)

func TestAtomicMap(t *testing.T) {
	var (
		wg      sync.WaitGroup
		setSize = 1000
		data    = randDataSet(setSize)
		aMap    = NewTypeAtomicMap()
	)

	for i := 0; i < len(data); i++ {
		wg.Add(1)
		go func(d Type) {
			aMap.Set(d.S, &d)

			fakeKey := randStr()
			fake, ok := aMap.Get(fakeKey)
			if fake != nil || ok {
				t.Errorf("found non-existing key %s\n", fakeKey)
			}

			back, ok := aMap.Get(d.S)
			if back.S != d.S || !ok {
				t.Errorf("Invalid key returned actual: %s expected: %s\n", back.S, d.S)
			}

			len := aMap.Len()
			if len < 1 || len > setSize {
				t.Errorf("Impossible map length: %d", len)
			}

			aMap.Delete(d.S)

			back, ok = aMap.Get(d.S)
			if back != nil || ok {
				t.Errorf("found deleted key %s\n", back.S)
			}

			wg.Done()
		}(data[i])
	}

	wg.Wait()
}

func randDataSet(size int) []Type {
	set := make([]Type, size)
	for k := range set {
		set[k] = Type{S: randStr()}
	}

	return set
}

var dictionary = "0123456789abcdefghijklmnopqrstuvwxyz"

func randStr() string {
	var bytes = make([]byte, 50)
	_, _ = rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}

	return string(bytes)
}
