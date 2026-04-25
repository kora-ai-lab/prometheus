package browser

import (
	"testing"
)

// TestCAPTCHA_Detect_recaptcha tests detection of Google reCAPTCHA
func TestCAPTCHA_Detect_recaptcha(t *testing.T) {
	html := `<div class="g-recaptcha" data-sitekey="test-key"></div>`
	if !DetectCAPTCHA(html) {
		t.Errorf("DetectCAPTCHA() should detect g-recaptcha, got false")
	}
}

// TestCAPTCHA_Detect_hcaptcha tests detection of hCaptcha
func TestCAPTCHA_Detect_hcaptcha(t *testing.T) {
	html := `<div class="h-captcha" data-sitekey="test-key"></div>`
	if !DetectCAPTCHA(html) {
		t.Errorf("DetectCAPTCHA() should detect h-captcha, got false")
	}
}

// TestCAPTCHA_Detect_cfchallenge tests detection of Cloudflare challenge
func TestCAPTCHA_Detect_cfchallenge(t *testing.T) {
	html := `<div class="cf-challenge"></div>`
	if !DetectCAPTCHA(html) {
		t.Errorf("DetectCAPTCHA() should detect cf-challenge, got false")
	}
}

// TestCAPTCHA_Detect_false tests that HTML without CAPTCHA returns false
func TestCAPTCHA_Detect_false(t *testing.T) {
	html := `<div class="normal-content">
		<p>This is regular HTML without any bot protection</p>
		<button>Submit</button>
	</div>`
	if DetectCAPTCHA(html) {
		t.Errorf("DetectCAPTCHA() should return false for HTML without CAPTCHA, got true")
	}
}

// TestCAPTCHA_Detect_caseinsensitive tests case-insensitive detection
func TestCAPTCHA_Detect_caseinsensitive(t *testing.T) {
	tests := []struct {
		name string
		html string
	}{
		{
			name: "uppercase G-RECAPTCHA",
			html: `<div class="G-RECAPTCHA"></div>`,
		},
		{
			name: "mixed case H-CaPtChA",
			html: `<div class="H-CaPtChA"></div>`,
		},
		{
			name: "uppercase CF-CHALLENGE",
			html: `<div class="CF-CHALLENGE"></div>`,
		},
		{
			name: "uppercase CHALLENGE",
			html: `<div class="CHALLENGE" data-type="captcha"></div>`,
		},
		{
			name: "uppercase RECAPTCHA",
			html: `<div class="RECAPTCHA"></div>`,
		},
		{
			name: "uppercase CAPTCHA",
			html: `<div class="CAPTCHA"></div>`,
		},
		{
			name: "uppercase DATA-SITEKEY",
			html: `<div DATA-SITEKEY="test-key"></div>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !DetectCAPTCHA(tt.html) {
				t.Errorf("DetectCAPTCHA() should detect %s (case-insensitive), got false", tt.name)
			}
		})
	}
}
