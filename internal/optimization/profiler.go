package optimization

import (
	"runtime"
	"sync"
	"time"
)

type MemSnapshot struct {
	Timestamp   time.Time
	AllocBytes uint64
	TotalAlloc uint64
	NumGC     uint32
	HeapAlloc uint64
}

type Profiler struct {
	mu       sync.RWMutex
	snapshots []MemSnapshot
	running  bool
}

func NewProfiler() *Profiler {
	return &Profiler{}
}

func (p *Profiler) Snapshot() MemSnapshot {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	
	return MemSnapshot{
		Timestamp:   time.Now(),
		AllocBytes:  stats.Alloc,
		TotalAlloc:  stats.TotalAlloc,
		NumGC:       stats.NumGC,
		HeapAlloc:   stats.HeapAlloc,
	}
}

func (p *Profiler) Record() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.snapshots = append(p.snapshots, p.Snapshot())
}

func (p *Profiler) GetSnapshots() []MemSnapshot {
	p.mu.RLock()
	defer p.mu.RUnlock()
	result := make([]MemSnapshot, len(p.snapshots))
	copy(result, p.snapshots)
	return result
}

func (p *Profiler) Start(interval time.Duration) {
	p.mu.Lock()
	p.running = true
	p.mu.Unlock()
	
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	
	for {
		p.mu.RLock()
		if !p.running {
			p.mu.RUnlock()
			return
		}
		p.mu.RUnlock()
		
		p.Record()
		
		select {
		case <-ticker.C:
		}
	}
}

func (p *Profiler) Stop() {
	p.mu.Lock()
	p.running = false
	p.mu.Unlock()
}