package security

import (
	"testing"
)

func TestDAST_Scan(t *testing.T) {
	cfg := DASTConfig{
		TargetURL: "http://httpbin.org/get",
		Timeout:   10 * 1000000000,
	}

	dast := NewDAST(cfg)
	result, err := dast.Scan()

	if err != nil {
		t.Logf("Scan error (network): %v", err)
	} else {
		t.Logf("Status: %d, Server: %s", result.StatusCode, result.Server)
	}
}

func TestDAST_Fuzz(t *testing.T) {
	cfg := DASTConfig{
		TargetURL: "http://httpbin.org",
	}

	dast := NewDAST(cfg)
	endpoints := dast.FuzzEndpoints()

	t.Logf("Found %d endpoints", len(endpoints))
	for _, e := range endpoints {
		t.Logf("  %s", e)
	}
}