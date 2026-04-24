package llm

type ModelEntry struct {
	Name          string
	Filename      string
	SizeBytes     int64
	MinRAMMb      int
	ContextWindow int
	URL           string
	SHA256        string
	Quality       string
	IsVision      bool
}

var TextModels = []ModelEntry{
	{
		Name:          "Qwen2.5 0.5B (minimal)",
		Filename:      "qwen2.5-0.5b-q4_k_m.gguf",
		SizeBytes:     390_000_000,
		MinRAMMb:      1500,
		ContextWindow: 32_768,
		Quality:       "minimal",
	},
	{
		Name:          "Phi-3 Mini 4K (recommended)",
		Filename:      "phi-3-mini-4k-instruct-q4_k_m.gguf",
		SizeBytes:     2_200_000_000,
		MinRAMMb:      3500,
		ContextWindow: 4096,
		Quality:       "recommended",
	},
	{
		Name:          "Llama 3.2 3B (balanced)",
		Filename:      "llama-3.2-3b-q4_k_m.gguf",
		SizeBytes:     2_000_000_000,
		MinRAMMb:      6000,
		ContextWindow: 8192,
		Quality:       "balanced",
	},
	{
		Name:          "Mistral 7B (high)",
		Filename:      "mistral-7b-q4_k_m.gguf",
		SizeBytes:     4_100_000_000,
		MinRAMMb:      12_000,
		ContextWindow: 32_768,
		Quality:       "high",
	},
}

var VisionModels = []ModelEntry{
	{
		Name:      "Moondream2 (minimal vision)",
		Filename:  "moondream2-q4_k_m.gguf",
		SizeBytes: 900_000_000,
		MinRAMMb:  2000,
		IsVision:  true,
	},
	{
		Name:      "Phi-3 Vision (recommended)",
		Filename:  "phi-3-vision-q4_k_m.gguf",
		SizeBytes: 2_500_000_000,
		MinRAMMb:  4000,
		IsVision:  true,
	},
}

func SelectModel(ramMb int) *ModelEntry {
	switch {
	case ramMb < 3000:
		return &TextModels[0]
	case ramMb < 6000:
		return &TextModels[1]
	case ramMb < 12000:
		return &TextModels[2]
	default:
		return &TextModels[3]
	}
}
