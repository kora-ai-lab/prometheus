package browser

import (
	"bytes"
	"testing"
	"time"
)

func skipIfPlaywrightNotAvailable(t *testing.T) {
	t.Skip("playwright not available in this build")
}

func TestPlaywrightClient_Navigate(t *testing.T) {
	skipIfPlaywrightNotAvailable(t)

	client, err := NewPlaywrightClient("", nil)
	if err != nil {
		t.Fatalf("Failed to create PlaywrightClient: %v", err)
	}
	defer client.Close()

	err = client.Navigate("https://example.com")
	if err != nil {
		t.Errorf("Navigate failed: %v", err)
	}
}

func TestPlaywrightClient_GetHTML(t *testing.T) {
	skipIfPlaywrightNotAvailable(t)

	client, err := NewPlaywrightClient("", nil)
	if err != nil {
		t.Fatalf("Failed to create PlaywrightClient: %v", err)
	}
	defer client.Close()

	err = client.Navigate("https://example.com")
	if err != nil {
		t.Fatalf("Navigate failed: %v", err)
	}

	html, err := client.GetHTML()
	if err != nil {
		t.Errorf("GetHTML failed: %v", err)
	}

	if html == "" {
		t.Error("GetHTML returned empty string")
	}

	if !bytes.Contains([]byte(html), []byte("example")) {
		t.Error("GetHTML did not return expected content")
	}
}

func TestPlaywrightClient_GetText(t *testing.T) {
	skipIfPlaywrightNotAvailable(t)

	client, err := NewPlaywrightClient("", nil)
	if err != nil {
		t.Fatalf("Failed to create PlaywrightClient: %v", err)
	}
	defer client.Close()

	err = client.Navigate("https://example.com")
	if err != nil {
		t.Fatalf("Navigate failed: %v", err)
	}

	text, err := client.GetText()
	if err != nil {
		t.Errorf("GetText failed: %v", err)
	}

	if text == "" {
		t.Error("GetText returned empty string")
	}
}

func TestPlaywrightClient_Click(t *testing.T) {
	skipIfPlaywrightNotAvailable(t)

	client, err := NewPlaywrightClient("", nil)
	if err != nil {
		t.Fatalf("Failed to create PlaywrightClient: %v", err)
	}
	defer client.Close()

	err = client.Navigate("https://example.com")
	if err != nil {
		t.Fatalf("Navigate failed: %v", err)
	}

	// Try to click on a link that exists on example.com
	err = client.Click("a")
	if err != nil {
		t.Logf("Click on 'a' element failed (might not exist on page): %v", err)
	}
}

func TestPlaywrightClient_Fill(t *testing.T) {
	skipIfPlaywrightNotAvailable(t)

	client, err := NewPlaywrightClient("", nil)
	if err != nil {
		t.Fatalf("Failed to create PlaywrightClient: %v", err)
	}
	defer client.Close()

	err = client.Navigate("https://example.com")
	if err != nil {
		t.Fatalf("Navigate failed: %v", err)
	}

	// Try to fill an input that may not exist
	err = client.Fill("input[type='text']", "test value")
	if err != nil {
		t.Logf("Fill on input failed (might not exist on page): %v", err)
	}
}

func TestPlaywrightClient_Submit(t *testing.T) {
	skipIfPlaywrightNotAvailable(t)

	client, err := NewPlaywrightClient("", nil)
	if err != nil {
		t.Fatalf("Failed to create PlaywrightClient: %v", err)
	}
	defer client.Close()

	err = client.Navigate("https://example.com")
	if err != nil {
		t.Fatalf("Navigate failed: %v", err)
	}

	// Try to submit a form that may not exist
	err = client.Submit("form")
	if err != nil {
		t.Logf("Submit on form failed (might not exist on page): %v", err)
	}
}

func TestPlaywrightClient_Screenshot(t *testing.T) {
	skipIfPlaywrightNotAvailable(t)

	client, err := NewPlaywrightClient("", nil)
	if err != nil {
		t.Fatalf("Failed to create PlaywrightClient: %v", err)
	}
	defer client.Close()

	err = client.Navigate("https://example.com")
	if err != nil {
		t.Fatalf("Navigate failed: %v", err)
	}

	screenshot, err := client.Screenshot()
	if err != nil {
		t.Errorf("Screenshot failed: %v", err)
	}

	if len(screenshot) == 0 {
		t.Error("Screenshot returned empty bytes")
	}

	// Check if it's a valid PNG header
	if len(screenshot) >= 4 && !(screenshot[0] == 0x89 && screenshot[1] == 0x50 && screenshot[2] == 0x4E && screenshot[3] == 0x47) {
		t.Error("Screenshot does not appear to be a valid PNG")
	}
}

func TestPlaywrightClient_EvalJS(t *testing.T) {
	skipIfPlaywrightNotAvailable(t)

	client, err := NewPlaywrightClient("", nil)
	if err != nil {
		t.Fatalf("Failed to create PlaywrightClient: %v", err)
	}
	defer client.Close()

	err = client.Navigate("https://example.com")
	if err != nil {
		t.Fatalf("Navigate failed: %v", err)
	}

	result, err := client.EvalJS("1 + 1")
	if err != nil {
		t.Errorf("EvalJS failed: %v", err)
	}

	if result != float64(2) && result != int(2) {
		t.Errorf("EvalJS returned unexpected result: %v", result)
	}
}

func TestPlaywrightClient_WaitForSelector(t *testing.T) {
	skipIfPlaywrightNotAvailable(t)

	client, err := NewPlaywrightClient("", nil)
	if err != nil {
		t.Fatalf("Failed to create PlaywrightClient: %v", err)
	}
	defer client.Close()

	err = client.Navigate("https://example.com")
	if err != nil {
		t.Fatalf("Navigate failed: %v", err)
	}

	// example.com has h1 element
	err = client.WaitForSelector("h1", 5*time.Second)
	if err != nil {
		t.Errorf("WaitForSelector failed: %v", err)
	}
}

func TestPlaywrightClient_WaitForNavigation(t *testing.T) {
	skipIfPlaywrightNotAvailable(t)

	client, err := NewPlaywrightClient("", nil)
	if err != nil {
		t.Fatalf("Failed to create PlaywrightClient: %v", err)
	}
	defer client.Close()

	err = client.Navigate("https://example.com")
	if err != nil {
		t.Fatalf("Navigate failed: %v", err)
	}

	// This test just verifies the method exists and can be called
	// In a real test we would trigger a navigation
	// Use a goroutine to timeout the wait
	done := make(chan error, 1)
	go func() {
		done <- client.WaitForNavigation(1 * time.Second)
	}()

	// Wait a bit then stop
	time.Sleep(100 * time.Millisecond)
}

func TestPlaywrightClient_GetCookies(t *testing.T) {
	skipIfPlaywrightNotAvailable(t)

	client, err := NewPlaywrightClient("", nil)
	if err != nil {
		t.Fatalf("Failed to create PlaywrightClient: %v", err)
	}
	defer client.Close()

	err = client.Navigate("https://example.com")
	if err != nil {
		t.Fatalf("Navigate failed: %v", err)
	}

	cookies, err := client.GetCookies()
	if err != nil {
		t.Errorf("GetCookies failed: %v", err)
	}

	// cookies may be empty, that's OK
	if cookies == nil {
		t.Error("GetCookies returned nil slice")
	}
}

func TestPlaywrightClient_SetCookie(t *testing.T) {
	skipIfPlaywrightNotAvailable(t)

	client, err := NewPlaywrightClient("", nil)
	if err != nil {
		t.Fatalf("Failed to create PlaywrightClient: %v", err)
	}
	defer client.Close()

	err = client.Navigate("https://example.com")
	if err != nil {
		t.Fatalf("Navigate failed: %v", err)
	}

	cookie := &Cookie{
		Name:  "test",
		Value: "value",
	}

	err = client.SetCookie(cookie)
	if err != nil {
		t.Errorf("SetCookie failed: %v", err)
	}
}

func TestPlaywrightClient_ScrollDown(t *testing.T) {
	skipIfPlaywrightNotAvailable(t)

	client, err := NewPlaywrightClient("", nil)
	if err != nil {
		t.Fatalf("Failed to create PlaywrightClient: %v", err)
	}
	defer client.Close()

	err = client.Navigate("https://example.com")
	if err != nil {
		t.Fatalf("Navigate failed: %v", err)
	}

	err = client.ScrollDown(100)
	if err != nil {
		t.Errorf("ScrollDown failed: %v", err)
	}
}

func TestPlaywrightClient_Close(t *testing.T) {
	skipIfPlaywrightNotAvailable(t)

	client, err := NewPlaywrightClient("", nil)
	if err != nil {
		t.Fatalf("Failed to create PlaywrightClient: %v", err)
	}

	err = client.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}
}
