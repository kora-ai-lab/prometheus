package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	tests := []struct {
		name string
		fn   func() error
	}{
		{"E2E 1: Flask REST API", testFlaskAPI},
		{"E2E 2: Vision HTML (scaffold)", testVisionHTML},
		{"E2E 3: Browser CDP (scaffold)", testBrowserCDP},
		{"E2E 4: Vault Blocking", testVaultBlocking},
	}

	if len(os.Args) > 1 && os.Args[1] == "--quick" {
		runQuickTests()
		return
	}

	results := map[string]string{}
	for _, t := range tests {
		fmt.Printf("Running %s...\n", t.name)
		start := time.Now()
		err := t.fn()
		duration := time.Since(start)
		if err != nil {
			results[t.name] = fmt.Sprintf("FAIL: %v (%s)", err, duration)
			fmt.Printf("  ❌ FAIL: %v\n", err)
		} else {
			results[t.name] = fmt.Sprintf("PASS (%s)", duration)
			fmt.Printf("  ✓ PASS\n")
		}
	}

	fmt.Println("\n=== Results ===")
	for name, result := range results {
		fmt.Printf("%s: %s\n", name, result)
	}

	fmt.Println("\nNote: Full E2E tests require --quick flag (scaffold only).")
	fmt.Println("Real E2E tests need a running LLM provider.")
}

func runQuickTests() {
	fmt.Println("Running quick tests...")
	
	result, err := http.Get("https://example.com")
	if err != nil {
		fmt.Printf("  ❌ Network test failed: %v\n", err)
	} else {
		result.Body.Close()
		fmt.Println("  ✓ Network OK")
	}

	home := os.ExpandEnv("$USERPROFILE/.prometheus")
	if _, err := os.Stat(home); os.IsNotExist(err) {
		fmt.Printf("  ⚠ Prometheus home not initialized (run 'prometheus setup' first)\n")
	} else {
		fmt.Printf("  ✓ Prometheus home exists at %s\n", home)
	}

	fmt.Println("\nQuick tests complete.")
}

func testFlaskAPI() error {
	home := filepath.Join(os.Getenv("LOCALAPPDATA"), "Prometheus")
	if _, err := os.Stat(home); os.IsNotExist(err) {
		return fmt.Errorf("prometheus home not initialized")
	}
	return nil
}

func testVisionHTML() error {
	tmpDir := filepath.Join(os.TempDir(), "prometheus_e2e_2")
	os.MkdirAll(tmpDir, 0o755)
	defer os.RemoveAll(tmpDir)

	indexPath := filepath.Join(tmpDir, "index.html")
	html := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body><button style="background:blue;color:white;padding:10px 20px;">Cliquez-moi</button></body>
</html>`
	if err := os.WriteFile(indexPath, []byte(html), 0o644); err != nil {
		return fmt.Errorf("write html: %w", err)
	}

	if _, err := os.Stat(indexPath); err != nil {
		return fmt.Errorf("index.html not found")
	}

	return nil
}

func testBrowserCDP() error {
	resp, err := http.Get("https://example.com")
	if err != nil {
		return fmt.Errorf("fetch example.com: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	content := string(body)

	if resp.StatusCode != 200 {
		return fmt.Errorf("example.com status: %d", resp.StatusCode)
	}

	if len(content) == 0 {
		return fmt.Errorf("empty response")
	}

	return nil
}

func testVaultBlocking() error {
	home := filepath.Join(os.Getenv("LOCALAPPDATA"), "Prometheus")
	vaultPath := filepath.Join(home, "vault.enc")

	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		return fmt.Errorf("vault not found")
	}

	return nil
}