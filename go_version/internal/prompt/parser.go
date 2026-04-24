package prompt

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Action struct {
	Thinking      string            `json:"thinking"`
	Action        string            `json:"action"`
	Command       string            `json:"command"`
	CreateFile    *CreateFileAction `json:"create_file,omitempty"`
	BrowserAction string            `json:"browser_action,omitempty"`
	BrowserArgs   map[string]string `json:"browser_args,omitempty"`
	VisionTarget  string            `json:"vision_target,omitempty"`
	VisionFile    string            `json:"vision_file,omitempty"`
	Question      string            `json:"question,omitempty"`
	Dangerous     bool              `json:"dangerous"`
	Why           string            `json:"why"`
}

type CreateFileAction struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type ParseError struct {
	Raw      string
	Msg      string
	ParseErr error
}

func (e *ParseError) Error() string {
	switch {
	case e.ParseErr != nil && e.Msg != "":
		return e.Msg + ": " + e.ParseErr.Error()
	case e.ParseErr != nil:
		return e.ParseErr.Error()
	default:
		return e.Msg
	}
}

func ParseAction(raw string) (*Action, error) {
	jsonStr := ExtractFirstJSON(raw)
	if jsonStr == "" {
		return nil, &ParseError{Raw: raw, Msg: "no JSON object found"}
	}

	var action Action
	if err := json.Unmarshal([]byte(jsonStr), &action); err != nil {
		return nil, &ParseError{Raw: raw, Msg: "invalid JSON action", ParseErr: err}
	}

	switch action.Action {
	case "exec":
		if action.Command == "" {
			return nil, &ParseError{Raw: raw, Msg: "exec action missing command"}
		}
	case "create":
		if action.CreateFile == nil || action.CreateFile.Path == "" {
			return nil, &ParseError{Raw: raw, Msg: "create action missing create_file"}
		}
	case "browser":
		if action.BrowserAction == "" {
			return nil, &ParseError{Raw: raw, Msg: "browser action missing browser_action"}
		}
	case "ask":
		if action.Question == "" {
			return nil, &ParseError{Raw: raw, Msg: "ask action missing question"}
		}
	case "vision", "done", "error":
	default:
		return nil, &ParseError{Raw: raw, Msg: "unknown action: " + action.Action}
	}

	return &action, nil
}

func ExtractFirstJSON(s string) string {
	start := strings.Index(s, "{")
	if start == -1 {
		return ""
	}

	depth := 0
	inString := false
	escaped := false
	for i := start; i < len(s); i++ {
		ch := s[i]
		if inString {
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == '"' {
				inString = false
			}
			continue
		}

		switch ch {
		case '"':
			inString = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return s[start : i+1]
			}
		}
	}
	return ""
}

func FormatParseRepair(err error, raw string) string {
	return fmt.Sprintf(
		"Ta réponse n'était pas du JSON valide. Réponds uniquement avec un objet JSON.\nErreur: %v\nExtrait: %s",
		err,
		truncateForRepair(raw),
	)
}

func truncateForRepair(raw string) string {
	if len(raw) <= 200 {
		return raw
	}
	return raw[:200]
}
