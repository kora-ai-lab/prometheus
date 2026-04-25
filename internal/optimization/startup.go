package optimization

import (
	"sync"
	"sync/atomic"
)

type LazyLoader struct {
	loaded atomic.Bool
	fn      func()
	result any
	once   sync.Once
}

func LazyLoad(fn func()) *LazyLoader {
	return &LazyLoader{
		fn: fn,
	}
}

func (l *LazyLoader) Get() any {
	l.once.Do(func() {
		if l.fn != nil {
			l.fn()
			l.loaded.Store(true)
		}
	})
	return l.result
}

func (l *LazyLoader) SetResult(val any) {
	l.result = val
}

func ParallelInit(fns ...func()) *sync.WaitGroup {
	var wg sync.WaitGroup

	for _, fn := range fns {
		wg.Add(1)
		go func(f func()) {
			defer wg.Done()
			f()
		}(fn)
	}

	return &wg
}

type StartupOptimizer struct {
	loaders []*LazyLoader
	wg      sync.WaitGroup
}

func NewStartupOptimizer() *StartupOptimizer {
	return &StartupOptimizer{}
}

func (so *StartupOptimizer) AddLazy(fn func()) {
	loader := LazyLoad(fn)
	so.loaders = append(so.loaders, loader)
}

func (so *StartupOptimizer) Preload() {
	so.wg.Add(len(so.loaders))

	for _, loader := range so.loaders {
		go func(l *LazyLoader) {
			defer so.wg.Done()
			l.Get()
		}(loader)
	}

	so.wg.Wait()
}