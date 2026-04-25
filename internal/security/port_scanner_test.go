package security

import (
	"testing"
)

func TestPortScanner_CommonPorts(t *testing.T) {
	results := ScanCommonPorts("127.0.0.1")

	t.Logf("Found %d open ports", len(results))

	for _, r := range results {
		t.Logf("Port %d: %s (risk: %s)", r.Port, r.Service, r.Risk)
	}
}

func TestScanPort_NotOpen(t *testing.T) {
	_, err := ScanPort("127.0.0.1", 59999)
	if err == nil {
		t.Error("Expected error for closed port")
	}
}