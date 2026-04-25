package sast

import (
	"testing"
)

func TestInjectionRules_Detect(t *testing.T) {
	rules := getInjectionRules()

	tests := []struct {
		code   string
		wantID string
	}{
		{`cursor.execute("SELECT * FROM users WHERE id = " + user_id)`, "INJ001"},
		{`"SELECT * FROM users".format(id)`, "INJ002"},

		{`os.system("ping " + host)`, "INJ010"},
		{`subprocess.run(cmd, shell=True, input=user_input)`, "INJ011"},
		{`eval(user_input)`, "INJ012"},

		{`div.innerHTML = user_data`, "INJ020"},

		{`open(filename, "r")`, "INJ030"},

		{`requests.get(url)`, "INJ050"},
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
			t.Logf("No rule matched for: %q", tt.code)
		}
	}
}