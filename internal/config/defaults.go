package config

import "strconv"

func DefaultConfig() *Config {
	return &Config{
		LLM: LLMConfig{
			Provider:   "local",
			ModelName:  "phi3:mini",
			Endpoint:   DefaultOllamaEndpoint,
			ServerPath: "",
		},
		Vision: VisionConfig{
			Enabled:     true,
			AutoCapture: true,
		},
		Browser: BrowserConfig{
			Enabled:  true,
			Level:    "cdp",
			Headless: true,
			Timeout:  30,
		},
		Security: SecurityConfig{
			RateLimitExecPerSec: 10,
			RateLimitLLMPerMin:  60,
			DangerousOpsConfirm: true,
			SandboxEnabled:      false,
		},
		Memory: MemoryConfig{
			CompactionThreshold: 0.70,
		},
		Logs: LogConfig{
			CompressAfterDays: 1,
			ArchiveAfterDays:  7,
			Format:            "jsonl",
		},
		UI: UIConfig{
			WebEnabled: false,
			WebPort:    DefaultWebPort,
			WebHost:    DefaultWebHost,
		},
	}
}

func DefaultConfigTOML() string {
	return `[llm]
provider = "local"
model_name = "phi3:mini"
endpoint = "` + DefaultOllamaEndpoint + `"
model_path = ""
vision_model_path = ""
server_path = ""

[vision]
enabled = true
auto_capture = true

[browser]
enabled = true
level = "cdp"
headless = true
timeout = 30

[security]
rate_limit_per_second = 10
rate_limit_llm_per_min = 60
dangerous_ops_confirmation = true
sandbox = false

[memory]
compaction_threshold = 0.70

[logs]
compress_after_days = 1
archive_after_days = 7
format = "jsonl"

[ui]
web_enabled = false
web_port = ` + strconv.Itoa(DefaultWebPort) + `
web_host = "` + DefaultWebHost + `"
`
}
