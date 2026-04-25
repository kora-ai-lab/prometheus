package metrics

import (
	"sync"
	"sync/atomic"
)

type Histogram struct {
	buckets []float64
	counts  []uint64
	sum    uint64
	mu     sync.Mutex
}

func NewHistogram(buckets []float64) *Histogram {
	return &Histogram{
		buckets: buckets,
		counts:  make([]uint64, len(buckets)+1),
	}
}

func (h *Histogram) Observe(value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	atomic.AddUint64(&h.sum, uint64(value*1000))

	bucket := len(h.buckets)
	for i, b := range h.buckets {
		if value < b {
			bucket = i
			break
		}
	}
	atomic.AddUint64(&h.counts[bucket], 1)
}

func (h *Histogram) Sum() float64 {
	return float64(atomic.LoadUint64(&h.sum)) / 1000
}

func (h *Histogram) Count() uint64 {
	var total uint64
	for _, c := range h.counts {
		total += c
	}
	return total
}