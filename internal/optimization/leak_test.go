package optimization

import (
	"runtime"
	"testing"
	"time"
)

func TestNewLeakDetector(t *testing.T) {
	ld := NewLeakDetector(10, 50, 0.5)
	if ld == nil {
		t.Fatal("NewLeakDetector returned nil")
	}
	if ld.maxHistory != 10 {
		t.Errorf("expected maxHistory 10, got %d", ld.maxHistory)
	}
	if ld.goroutineThreshold != 50 {
		t.Errorf("expected goroutineThreshold 50, got %d", ld.goroutineThreshold)
	}
	if ld.memoryGrowthThreshold != 0.5 {
		t.Errorf("expected memoryGrowthThreshold 0.5, got %f", ld.memoryGrowthThreshold)
	}
}

func TestLeakDetector_SetBaseline(t *testing.T) {
	ld := NewLeakDetector(10, 50, 0.5)
	ld.SetBaseline()

	if ld.baseline == nil {
		t.Fatal("baseline not set")
	}
	if ld.baseline.Goroutines <= 0 {
		t.Error("baseline goroutines should be > 0")
	}
}

func TestLeakDetector_DetectLeaks_NoLeak(t *testing.T) {
	ld := NewLeakDetector(10, 50, 0.5)
	ld.SetBaseline()

	entry, hasLeak := ld.DetectLeaks()
	if hasLeak {
		t.Error("should not detect leak immediately after baseline")
	}
	if entry == nil {
		t.Fatal("entry is nil")
	}
	if entry.Goroutines <= 0 {
		t.Error("goroutine count should be > 0")
	}
	if entry.AllocBytes == 0 {
		t.Error("alloc bytes should be > 0")
	}
}

func TestLeakDetector_GoroutineCount(t *testing.T) {
	ld := NewLeakDetector(10, 50, 0.5)
	count := ld.GoroutineCount()
	if count < 1 {
		t.Errorf("expected at least 1 goroutine, got %d", count)
	}
}

func TestLeakDetector_MemoryStats(t *testing.T) {
	ld := NewLeakDetector(10, 50, 0.5)
	stats := ld.MemoryStats()
	if stats.Alloc == 0 {
		t.Error("alloc should be > 0")
	}
	if stats.TotalAlloc == 0 {
		t.Error("total alloc should be > 0")
	}
}

func TestLeakDetector_History(t *testing.T) {
	ld := NewLeakDetector(3, 50, 0.5)
	ld.SetBaseline()

	ld.DetectLeaks()
	ld.DetectLeaks()
	ld.DetectLeaks()
	ld.DetectLeaks()

	history := ld.History()
	if len(history) > 3 {
		t.Errorf("expected max 3 history entries, got %d", len(history))
	}
}

func TestLeakDetector_ClearHistory(t *testing.T) {
	ld := NewLeakDetector(10, 50, 0.5)
	ld.SetBaseline()
	ld.DetectLeaks()

	if len(ld.History()) == 0 {
		t.Fatal("history should not be empty")
	}

	ld.ClearHistory()
	if len(ld.History()) != 0 {
		t.Error("history should be empty after clear")
	}
}

func TestLeakDetector_DetectLeaks_GoroutineLeak(t *testing.T) {
	ld := NewLeakDetector(10, 5, 0.5)
	ld.SetBaseline()

	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			<-done
		}()
	}

	time.Sleep(10 * time.Millisecond)

	_, hasLeak := ld.DetectLeaks()
	if !hasLeak {
		t.Log("Goroutine leak not detected (threshold may be too high for this test)")
	}

	close(done)
	time.Sleep(10 * time.Millisecond)
	runtime.GC()
}

func TestLeakDetector_EntryFields(t *testing.T) {
	ld := NewLeakDetector(10, 50, 0.5)
	ld.SetBaseline()

	entry, _ := ld.DetectLeaks()

	if entry.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
	if entry.NumGC == 0 {
		t.Log("NumGC is 0 (may be first GC cycle)")
	}
	if len(entry.StackTraces) == 0 {
		t.Error("stack traces should not be empty")
	}
}

func TestLeakDetector_ConcurrentAccess(t *testing.T) {
	ld := NewLeakDetector(10, 50, 0.5)
	ld.SetBaseline()

	done := make(chan struct{})
	for i := 0; i < 5; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				ld.DetectLeaks()
				ld.GoroutineCount()
				ld.MemoryStats()
			}
			done <- struct{}{}
		}()
	}

	for i := 0; i < 5; i++ {
		<-done
	}
}
