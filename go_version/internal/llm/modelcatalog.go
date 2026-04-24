package llm

type ModelEntry struct {
	Name          string
	Filename     string
	SizeBytes    int64
	MinRAMMb     int
	ContextWindow int
	URL          string
	SHA256       string
	Quality      string
	IsVision     bool
}

var TextModels = []ModelEntry{
	{
		Name:           "Qwen2.5 0.5B (minimal)",
		Filename:      "qwen2.5-0.5b-instruct-q4_k_m.gguf",
		SizeBytes:     398_000_000,
		MinRAMMb:     1500,
		ContextWindow: 32_768,
		URL:          "https://huggingface.co/Qwen/Qwen2.5-0.5B-Instruct-GGUF/resolve/main/qwen2.5-0.5b-instruct-q4_k_m.gguf",
		SHA256:       "bfb8508e3f15c0aa6c5c5599af03681d953e87044c2e2e031d52cd8d07f85d21",
		Quality:      "minimal",
	},
	{
		Name:           "Phi-3 Mini 4K (recommended)",
		Filename:      "phi-3-mini-4k-instruct-q4.gguf",
		SizeBytes:     2_200_000_000,
		MinRAMMb:     3500,
		ContextWindow: 4096,
		URL:          "https://huggingface.co/microsoft/Phi-3-mini-4k-instruct-gguf/resolve/main/Phi-3-mini-4k-instruct-q4.gguf",
		SHA256:       "8a83c7fb9049a9b2e92266fa7ad0493a3a1e",
		Quality:      "recommended",
	},
	{
		Name:           "Llama 3.2 3B (balanced)",
		Filename:      "llama-3.2-3b-instruct-q4_k_m.gguf",
		SizeBytes:    2_000_000_000,
		MinRAMMb:     6000,
		ContextWindow: 8192,
		URL:          "https://huggingface.co/unsloth/Llama-3.2-3B-Instruct-GGUF/resolve/main/Llama-3.2-3B-Instruct-q4_k_m.gguf",
		SHA256:       "",
		Quality:      "balanced",
	},
	{
		Name:           "Mistral 7B (high)",
		Filename:      "mistral-7b-v0.3-q4_k_m.gguf",
		SizeBytes:     4_100_000_000,
		MinRAMMb:     12_000,
		ContextWindow: 32_768,
		URL:          "https://huggingface.co/TheBloke/Mistral-7B-v0.3-GGUF/resolve/main/mistral-7b-v0.3.Q4_K_M.gguf",
		SHA256:       "",
		Quality:      "high",
	},
}

var VisionModels = []ModelEntry{
	{
		Name:       "Moondream2 (minimal vision)",
		Filename:   "moondream2-q4_k_m.gguf",
		SizeBytes: 900_000_000,
		MinRAMMb:  2000,
		URL:        "https://huggingface.co/vikhyatk/moondream2-gguf/resolve/main/moondream2-q4_k_m.gguf",
		SHA256:     "",
		IsVision:   true,
	},
	{
		Name:       "Phi-3 Vision (recommended)",
		Filename:   "phi-3-vision-instruct-q4.gguf",
		SizeBytes: 2_500_000_000,
		MinRAMMb:  4000,
		URL:       "https://huggingface.co/microsoft/Phi-3-vision-gguf/resolve/main/Phi-3-vision-instruct-q4.gguf",
		SHA256:     "",
		IsVision:   true,
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

func SelectVisionModel(ramMb int) *ModelEntry {
	switch {
	case ramMb < 4000:
		return &VisionModels[0]
	default:
		return &VisionModels[1]
	}
}