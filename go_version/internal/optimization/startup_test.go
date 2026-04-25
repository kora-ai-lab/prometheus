package optimization

import (
	"sync"
	"testing"
)

func TestLazyLoad(t *testing.T) {
	called := false
	loader := LazyLoad(func() {
		called = true
	})

	loader.Get()

	if !called {
		t.Error("Expected lazy function to be called")
	}
}

func TestLazyLoad_Once(t *testing.T) {
	count := 0
	loader := LazyLoad(func() {
		count++
	})

	loader.Get()
	loader.Get()
	loader.Get()

	if count != 1 {
		t.Errorf("Expected 1 call, got %d", count)
	}
}

func TestParallelInit(t *testing.T) {
	var results []int

	var mu sync.Mutex
	wg := ParallelInit(
		func() { mu.Lock(); results = append(results, 1); mu.Unlock() },
		func() { mu.Lock(); results = append(results, 2); mu.Unlock() },
		func() { mu.Lock(); results = append(results, 3); mu.Unlock() },
	)

	wg.Wait()

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}
}