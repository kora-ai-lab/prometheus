package sast

import (
	"testing"
)

func TestSecretsRules_Detect(t *testing.T) {
	rules := getSecretsRules()
	
	tests := []struct {
		code   string
		wantID string
	}{
		{`api_key = "sk-abcdef123456789"`, "SEC001"},
		{`password = "secret123"`, "SEC002"},
		{`jwt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`, "SEC003"},
		{`-----BEGIN RSA PRIVATE KEY-----`, "SEC004"},
		{`ghp_abcdefghijklmnopqrstuvwxyz123456789012`, "SEC005"},
		{`hashlib.md5(password)`, "SEC006"},
		{`verify = False`, "SEC007"},
	}
	
	for _, tt := range tests {
		matched := false
		for _, rule := range rules {
			if rule.Pattern.MatchString(tt.code) {
				if rule.ID != tt.wantID {
					t.Errorf("Rule %s matched for %q, want %s", rule.ID, tt.code, tt.wantID)
				}
				matched = true
				break
			}
		}
		if !matched {
			t.Errorf("No rule matched for: %q", tt.code)
		}
	}
}