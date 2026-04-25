package metrics

import "testing"

func TestHistogram(t *testing.T) {
	h := NewHistogram([]float64{1, 5, 10})

	h.Observe(0.5)
	h.Observe(3)
	h.Observe(7)

	if h.Count() != 3 {
		t.Errorf("Expected count 3, got %d", h.Count())
	}

	if h.Sum() != 10.5 {
		t.Errorf("Expected sum 10.5, got %f", h.Sum())
	}
}