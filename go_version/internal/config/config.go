package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	LLM      LLMConfig
	Vision   VisionConfig
	Browser  BrowserConfig
	Security SecurityConfig
	Memory   MemoryConfig
	Logs     LogConfig
	UI       UIConfig
}

type LLMConfig struct {
	Provider        string
	ModelPath       string
	VisionModelPath string
	ModelName       string
	Endpoint        string
	ServerPath      string
}

type VisionConfig struct {
	Enabled     bool
	AutoCapture bool
}

type BrowserConfig struct {
	Enabled  bool
	Level    string
	Headless bool
	Timeout  int
}

type SecurityConfig struct {
	RateLimitExecPerSec int
	RateLimitLLMPerMin  int
	DangerousOpsConfirm bool
	SandboxEnabled      bool
}

type MemoryConfig struct {
	CompactionThreshold float64
}

type LogConfig struct {
	CompressAfterDays int
	ArchiveAfterDays  int
	Format            string
}

type UIConfig struct {
	WebEnabled bool
	WebPort    int
	WebHost    string
}

func PrometheusHome() string {
	if v := os.Getenv("PROMETHEUS_HOME"); v != "" {
		return v
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ".prometheus"
	}
	return filepath.Join(home, ".prometheus")
}

func EnsureHome() (string, error) {
	home := PrometheusHome()
	dirs := []string{
		home,
		filepath.Join(home, "runtime"),
		filepath.Join(home, "models"),
		filepath.Join(home, "logs"),
		filepath.Join(home, "security"),
		filepath.Join(home, "prompts"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return "", err
		}
	}

	cfgPath := filepath.Join(home, "prometheus.toml")
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		if err := os.WriteFile(cfgPath, []byte(DefaultConfigTOML()), 0o600); err != nil {
			return "", err
		}
	}

	return home, nil
}

func Load(home string) (*Config, error) {
	cfg := DefaultConfig()
	path := filepath.Join(home, "prometheus.toml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := parseMinimalTOML(cfg, string(data)); err != nil {
		return nil, err
	}
	applyEnvOverrides(cfg)
	return cfg, nil
}

func parseMinimalTOML(cfg *Config, raw string) error {
	section := ""
	lines := strings.Split(raw, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.Trim(line, "[]")
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid config line %d: %q", i+1, line)
		}
		key := strings.TrimSpace(parts[0])
		rawVal := strings.TrimSpace(parts[1])
		val := rawVal
		if strings.HasPrefix(rawVal, `"`) {
			unquoted, err := strconv.Unquote(rawVal)
			if err != nil {
				return fmt.Errorf("invalid quoted value on line %d: %w", i+1, err)
			}
			val = unquoted
		}
		if err := setValue(cfg, section, key, val); err != nil {
			return fmt.Errorf("config field %s.%s: %w", section, key, err)
		}
	}
	return nil
}

func setValue(cfg *Config, section, key, val string) error {
	switch section + "." + key {
	case "llm.provider":
		cfg.LLM.Provider = val
	case "llm.model_name":
		cfg.LLM.ModelName = val
	case "llm.endpoint":
		cfg.LLM.Endpoint = val
	case "llm.model_path":
		cfg.LLM.ModelPath = val
	case "llm.vision_model_path":
		cfg.LLM.VisionModelPath = val
	case "llm.server_path":
		cfg.LLM.ServerPath = val
	case "vision.enabled":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		cfg.Vision.Enabled = b
	case "vision.auto_capture":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		cfg.Vision.AutoCapture = b
	case "browser.enabled":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		cfg.Browser.Enabled = b
	case "browser.level":
		cfg.Browser.Level = val
	case "browser.headless":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		cfg.Browser.Headless = b
	case "browser.timeout":
		n, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		cfg.Browser.Timeout = n
	case "security.rate_limit_per_second":
		n, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		cfg.Security.RateLimitExecPerSec = n
	case "security.rate_limit_llm_per_min":
		n, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		cfg.Security.RateLimitLLMPerMin = n
	case "security.dangerous_ops_confirmation":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		cfg.Security.DangerousOpsConfirm = b
	case "security.sandbox":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		cfg.Security.SandboxEnabled = b
	case "memory.compaction_threshold":
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		cfg.Memory.CompactionThreshold = f
	case "logs.compress_after_days":
		n, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		cfg.Logs.CompressAfterDays = n
	case "logs.archive_after_days":
		n, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		cfg.Logs.ArchiveAfterDays = n
	case "logs.format":
		cfg.Logs.Format = val
	case "ui.web_enabled":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		cfg.UI.WebEnabled = b
	case "ui.web_port":
		n, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		cfg.UI.WebPort = n
	case "ui.web_host":
		cfg.UI.WebHost = val
	default:
		return nil
	}

	return nil
}

func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("PROMETHEUS_LLM_PROVIDER"); v != "" {
		cfg.LLM.Provider = v
	}
	if v := os.Getenv("PROMETHEUS_LLM_MODEL_PATH"); v != "" {
		cfg.LLM.ModelPath = v
	}
	if v := os.Getenv("PROMETHEUS_LLM_SERVER_PATH"); v != "" {
		cfg.LLM.ServerPath = v
	}
	if v := os.Getenv("PROMETHEUS_WEB_ENABLED"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.UI.WebEnabled = b
		}
	}
}

func Save(home string, cfg *Config) error {
	return os.WriteFile(filepath.Join(home, "prometheus.toml"), []byte(renderTOML(cfg)), 0o600)
}

func UpdateLLM(home string, update func(*LLMConfig)) error {
	cfg, err := Load(home)
	if err != nil {
		return err
	}
	update(&cfg.LLM)
	return Save(home, cfg)
}

func renderTOML(cfg *Config) string {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	return fmt.Sprintf(`[llm]
provider = %q
model_name = %q
endpoint = %q
model_path = %q
vision_model_path = %q
server_path = %q

[vision]
enabled = %t
auto_capture = %t

[browser]
enabled = %t
level = %q
headless = %t
timeout = %d

[security]
rate_limit_per_second = %d
rate_limit_llm_per_min = %d
dangerous_ops_confirmation = %t
sandbox = %t

[memory]
compaction_threshold = %.2f

[logs]
compress_after_days = %d
archive_after_days = %d
format = %q

[ui]
web_enabled = %t
web_port = %d
web_host = %q
`,
		cfg.LLM.Provider,
		cfg.LLM.ModelName,
		cfg.LLM.Endpoint,
		cfg.LLM.ModelPath,
		cfg.LLM.VisionModelPath,
		cfg.LLM.ServerPath,
		cfg.Vision.Enabled,
		cfg.Vision.AutoCapture,
		cfg.Browser.Enabled,
		cfg.Browser.Level,
		cfg.Browser.Headless,
		cfg.Browser.Timeout,
		cfg.Security.RateLimitExecPerSec,
		cfg.Security.RateLimitLLMPerMin,
		cfg.Security.DangerousOpsConfirm,
		cfg.Security.SandboxEnabled,
		cfg.Memory.CompactionThreshold,
		cfg.Logs.CompressAfterDays,
		cfg.Logs.ArchiveAfterDays,
		cfg.Logs.Format,
		cfg.UI.WebEnabled,
		cfg.UI.WebPort,
		cfg.UI.WebHost,
	)
}
