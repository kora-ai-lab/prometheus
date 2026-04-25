package optimization

import (
	"testing"
	"time"
)

func TestProfiler_Snapshot(t *testing.T) {
	p := NewProfiler()
	snapshot := p.Snapshot()

	if snapshot.AllocBytes == 0 {
		t.Error("Expected non-zero allocation")
	}
}

func TestProfiler_Record(t *testing.T) {
	p := NewProfiler()
	p.Record()

	snapshots := p.GetSnapshots()
	if len(snapshots) != 1 {
		t.Errorf("Expected 1 snapshot, got %d", len(snapshots))
	}
}

func TestProfiler_StartStop(t *testing.T) {
	p := NewProfiler()
	go p.Start(10 * time.Millisecond)
	time.Sleep(50 * time.Millisecond)
	p.Stop()

	snapshots := p.GetSnapshots()
	if len(snapshots) == 0 {
		t.Error("Expected snapshots after start")
	}
}