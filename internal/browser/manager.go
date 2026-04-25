package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/kora-ai-lab/prometheus/internal/capabilities"
	"github.com/kora-ai-lab/prometheus/internal/prompt"
	"github.com/kora-ai-lab/prometheus/internal/vision"
)

type Manager struct {
	client    Client
	capEngine *capabilities.Engine
	vision    vision.VisionProvider
}

func NewManager(capEngine *capabilities.Engine, visionProvider vision.VisionProvider) *Manager {
	return &Manager{
		client:    NewStubClient(),
		capEngine: capEngine,
		vision:    visionProvider,
	}
}

func (m *Manager) Do(ctx context.Context, action *prompt.Action) string {
	_ = ctx
	switch action.BrowserAction {
	case "navigate":
		if err := m.client.Navigate(action.BrowserArgs["url"]); err != nil {
			return "ERROR: " + err.Error()
		}
		return "ACTION OK"
	case "get_html":
		html, err := m.client.GetHTML()
		if err != nil {
			return "ERROR: " + err.Error()
		}
		return html
	case "screenshot":
		img, err := m.client.Screenshot()
		if err != nil {
			return "ERROR: " + err.Error()
		}
		if m.vision != nil && m.vision.HasVision() {
			analysis, _ := m.vision.Analyze(ctx, img, "Describe the screenshot.")
			return analysis
		}
		return "SCREENSHOT_CAPTURED"
	case "eval_js":
		out, err := m.client.EvalJS(action.BrowserArgs["script"])
		if err != nil {
			return "ERROR: " + err.Error()
		}
		raw, _ := json.Marshal(out)
		return string(raw)
	case "wait_for":
		timeout := 10 * time.Second
		if v := action.BrowserArgs["timeout"]; v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				timeout = time.Duration(n) * time.Second
			}
		}
		if err := m.client.WaitForSelector(action.BrowserArgs["selector"], timeout); err != nil {
			return "ERROR: " + err.Error()
		}
		return "ACTION OK"
	}
	return fmt.Sprintf("ERROR: browser action %q not implemented", action.BrowserAction)
}

func (m *Manager) Screenshot(ctx context.Context) ([]byte, error) {
	_ = ctx
	return m.client.Screenshot()
}

func (m *Manager) VisionResult(ctx context.Context, action *prompt.Action) string {
	var img []byte
	var err error
	switch action.VisionTarget {
	case "browser":
		img, err = m.client.Screenshot()
	default:
		err = fmt.Errorf("vision target %q is not implemented", action.VisionTarget)
	}
	if err != nil {
		return "VISION ERROR: " + err.Error()
	}
	analysis, err := m.vision.Analyze(ctx, img, action.Why)
	if err != nil {
		return "VISION ERROR: " + err.Error()
	}
	return "[VISION]\n" + analysis
}

func (m *Manager) Close() error {
	if m.client == nil {
		return nil
	}
	return m.client.Close()
}
