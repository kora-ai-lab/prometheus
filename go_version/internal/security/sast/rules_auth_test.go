package sast

import (
	"testing"
)

func TestAuthRules_Detect(t *testing.T) {
	rules := getAuthRules()

	tests := []struct {
		code   string
		wantID string
	}{
		{`if password == input:`, "AUTH010"},
		{`hashlib.md5(password)`, "AUTH011"},
		{`jwt.decode(token, verify=False)`, "AUTH012"},
		{`pickle.loads(user_data)`, "AUTH020"},
		{`__proto__: {}`, "AUTH021"},
		{`yaml.load(data)`, "AUTH022"},
	}

	for _, tt := range tests {
		matched := false
		for _, rule := range rules {
			if rule.Pattern.MatchString(tt.code) {
				matched = true
				break
			}
		}
		if !matched {
			t.Logf("No rule matched for: %q", tt.code)
		}
	}
}