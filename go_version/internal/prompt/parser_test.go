package prompt

import "testing"

func TestParseActionExtractsFirstJSON(t *testing.T) {
	raw := `noise before {"action":"done","dangerous":false,"why":"ok"} trailing noise`
	action, err := ParseAction(raw)
	if err != nil {
		t.Fatalf("ParseAction() error = %v", err)
	}
	if action.Action != "done" {
		t.Fatalf("Action = %q, want done", action.Action)
	}
}
