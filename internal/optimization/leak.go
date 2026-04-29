package optimization

import (
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

// LeakEntry represents a detected memory leak
type LeakEntry struct {
	Goroutines   int
	AllocBytes   uint64
	TotalAlloc   uint64
	SysBytes     uint64
	NumGC        uint32
	Timestamp    time.Time
	StackTraces  []string
}

// LeakDetector monitors goroutine and memory leaks
type LeakDetector struct {
	mu              sync.Mutex
	baseline        *LeakEntry
	history         []LeakEntry
	maxHistory      int
	goroutineThreshold int
	memoryGrowthThreshold float64
}

// NewLeakDetector creates a detector with configurable thresholds
func NewLeakDetector(maxHistory int, goroutineThreshold int, memoryGrowthThreshold float64) *LeakDetector {
	return &LeakDetector{
		maxHistory:           maxHistory,
		goroutineThreshold:   goroutineThreshold,
		memoryGrowthThreshold: memoryGrowthThreshold,
	}
}

// SetBaseline captures the current state as the baseline for comparison
func (ld *LeakDetector) SetBaseline() {
	ld.mu.Lock()
	defer ld.mu.Unlock()
	ld.baseline = ld.snapshot()
}

// DetectLeaks checks for memory and goroutine leaks compared to baseline
func (ld *LeakDetector) DetectLeaks() (*LeakEntry, bool) {
	ld.mu.Lock()
	defer ld.mu.Unlock()

	forceGC()
	current := ld.snapshot()
	ld.history = append(ld.history, *current)
	if len(ld.history) > ld.maxHistory {
		ld.history = ld.history[len(ld.history)-ld.maxHistory:]
	}

	if ld.baseline == nil {
		return current, false
	}

	hasLeak := ld.analyzeLeak(current)
	return current, hasLeak
}

// GoroutineCount returns the current number of active goroutines
func (ld *LeakDetector) GoroutineCount() int {
	return runtime.NumGoroutine()
}

// MemoryStats returns current memory statistics
func (ld *LeakDetector) MemoryStats() runtime.MemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m
}

// History returns the leak detection history
func (ld *LeakDetector) History() []LeakEntry {
	ld.mu.Lock()
	defer ld.mu.Unlock()
	result := make([]LeakEntry, len(ld.history))
	copy(result, ld.history)
	return result
}

// ClearHistory resets the detection history
func (ld *LeakDetector) ClearHistory() {
	ld.mu.Lock()
	defer ld.mu.Unlock()
	ld.history = nil
}

func (ld *LeakDetector) snapshot() *LeakEntry {
	forceGC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &LeakEntry{
		Goroutines: runtime.NumGoroutine(),
		AllocBytes:   m.Alloc,
		TotalAlloc:   m.TotalAlloc,
		SysBytes:     m.Sys,
		NumGC:        m.NumGC,
		Timestamp:    time.Now(),
		StackTraces:  collectGoroutineStacks(),
	}
}

func (ld *LeakDetector) analyzeLeak(current *LeakEntry) bool {
	goroutineDelta := current.Goroutines - ld.baseline.Goroutines
	if goroutineDelta > ld.goroutineThreshold {
		return true
	}

	if ld.baseline.AllocBytes > 0 {
		memoryGrowth := float64(current.AllocBytes-ld.baseline.AllocBytes) / float64(ld.baseline.AllocBytes)
		if memoryGrowth > ld.memoryGrowthThreshold {
			return true
		}
	}

	return false
}

func forceGC() {
	runtime.GC()
	debug.FreeOSMemory()
}

func collectGoroutineStacks() []string {
	buf := make([]byte, 1<<16)
	n := runtime.Stack(buf, true)
	if n >= len(buf) {
		return []string{"stack trace truncated"}
	}
	return []string{string(buf[:n])}
}
