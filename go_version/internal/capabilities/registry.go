package capabilities

type Capability struct {
	Name         string
	CheckCmd     string
	InstallCmds map[string]string
	Type        string
	SizeMb      int
	Description string
}

var Registry = map[string]Capability{
	// Dev tools
	"python3": {
		Name:        "python3",
		CheckCmd:    "python3 --version",
		Type:        "system",
		SizeMb:      30,
		Description: "Python interpreter",
		InstallCmds: map[string]string{
			"apt":    "sudo apt-get install -y python3",
			"brew":  "brew install python3",
			"winget": "winget install Python.Python.3.11",
		},
	},
	"python": {
		Name:        "python",
		CheckCmd:    "python --version",
		Type:        "system",
		SizeMb:      30,
		Description: "Python interpreter (alias)",
		InstallCmds: map[string]string{
			"apt":    "sudo apt-get install -y python3",
			"brew":  "brew install python3",
		},
	},
	"node": {
		Name:        "node",
		CheckCmd:    "node --version",
		Type:        "system",
		SizeMb:      30,
		Description: "Node.js runtime",
		InstallCmds: map[string]string{
			"apt":    "curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash - && sudo apt-get install -y nodejs",
			"brew":  "brew install node",
			"winget": "winget install OpenJS.NodeJS.LTS",
		},
	},
	"npm": {
		Name:        "npm",
		CheckCmd:    "npm --version",
		Type:        "npm",
		SizeMb:      10,
		Description: "Node package manager",
		InstallCmds: map[string]string{
			"apt":    "sudo apt-get install -y npm",
			"brew":  "brew install npm",
		},
	},
	"go": {
		Name:        "go",
		CheckCmd:    "go version",
		Type:        "system",
		SizeMb:      100,
		Description: "Go compiler",
		InstallCmds: map[string]string{
			"apt":    "sudo apt-get install -y golang-go",
			"brew":  "brew install go",
		},
	},
	"rust": {
		Name:        "rust",
		CheckCmd:    "rustc --version",
		Type:        "system",
		SizeMb:      200,
		Description: "Rust toolchain",
		InstallCmds: map[string]string{
			"apt":    "curl --proto =https -fsSL https://sh.rustup.rs | sh",
			"brew":  "brew install rustup-init && rustup-init",
		},
	},
	"cargo": {
		Name:        "cargo",
		CheckCmd:    "cargo --version",
		Type:        "cargo",
		SizeMb:      100,
		Description: "Rust package manager",
	},
	"java": {
		Name:        "java",
		CheckCmd:    "java -version",
		Type:        "system",
		SizeMb:      200,
		Description: "Java runtime",
		InstallCmds: map[string]string{
			"apt":    "sudo apt-get install -y openjdk-17-jdk",
			"brew":  "brew install openjdk",
		},
	},
	"maven": {
		Name:        "maven",
		CheckCmd:    "mvn --version",
		Type:        "system",
		SizeMb:      20,
		Description: "Java build tool",
		InstallCmds: map[string]string{
			"apt":    "sudo apt-get install -y maven",
			"brew":  "brew install maven",
		},
	},

	// Web tools
	"curl": {
		Name:        "curl",
		CheckCmd:    "curl --version",
		Type:        "system",
		SizeMb:      10,
		Description: "HTTP client",
		InstallCmds: map[string]string{
			"apt":    "sudo apt-get install -y curl",
			"brew":  "brew install curl",
		},
	},
	"wget": {
		Name:        "wget",
		CheckCmd:    "wget --version",
		Type:        "system",
		SizeMb:      10,
		Description: "File downloader",
		InstallCmds: map[string]string{
			"apt":    "sudo apt-get install -y wget",
			"brew":  "brew install wget",
		},
	},
	"httpie": {
		Name:        "httpie",
		CheckCmd:    "http --version",
		Type:        "pip",
		SizeMb:      5,
		Description: "User-friendly HTTP client",
		InstallCmds: map[string]string{
			"pip":   "pip install httpie",
			"brew":  "brew install httpie",
		},
	},

	// Data tools
	"jq": {
		Name:        "jq",
		CheckCmd:    "jq --version",
		Type:        "system",
		SizeMb:      5,
		Description: "JSON processor",
		InstallCmds: map[string]string{
			"apt":    "sudo apt-get install -y jq",
			"brew":  "brew install jq",
			"winget": "winget install jqlang.jq",
		},
	},
	"sqlite3": {
		Name:        "sqlite3",
		CheckCmd:    "sqlite3 --version",
		Type:        "system",
		SizeMb:      10,
		Description: "SQLite CLI",
		InstallCmds: map[string]string{
			"apt":    "sudo apt-get install -y sqlite3",
			"brew":  "brew install sqlite3",
		},
	},
	"redis-cli": {
		Name:        "redis-cli",
		CheckCmd:    "redis-cli --version",
		Type:        "system",
		SizeMb:      10,
		Description: "Redis client",
		InstallCmds: map[string]string{
			"apt":    "sudo apt-get install -y redis-tools",
		},
	},
	"pandas": {
		Name:        "pandas",
		CheckCmd:    "python3 -c 'import pandas'",
		Type:        "pip",
		SizeMb:      30,
		Description: "Python data analysis",
		InstallCmds: map[string]string{
			"pip":   "pip install pandas",
		},
	},
	"requests": {
		Name:        "requests",
		CheckCmd:    "python3 -c 'import requests'",
		Type:        "pip",
		SizeMb:      5,
		Description: "Python HTTP library",
		InstallCmds: map[string]string{
			"pip":   "pip install requests",
		},
	},

	// DevOps tools
	"docker": {
		Name:        "docker",
		CheckCmd:    "docker --version",
		Type:        "system",
		SizeMb:      100,
		Description: "Container runtime",
		InstallCmds: map[string]string{
			"apt":    "sudo apt-get install -y docker.io && sudo systemctl start docker",
			"brew":  "brew install docker",
		},
	},
	"kubectl": {
		Name:        "kubectl",
		CheckCmd:    "kubectl version --client",
		Type:        "system",
		SizeMb:      50,
		Description: "Kubernetes CLI",
		InstallCmds: map[string]string{
			"curl":  "curl -LO 'https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl'",
			"brew":  "brew install kubectl",
		},
	},

	// Security tools
	"semgrep": {
		Name:        "semgrep",
		CheckCmd:    "semgrep --version",
		Type:        "pip",
		SizeMb:      50,
		Description: "Static analysis security scanner",
		InstallCmds: map[string]string{
			"pip":   "pip install semgrep",
			"brew":  "brew install semgrep",
		},
	},

	// Media tools
	"ffmpeg": {
		Name:        "ffmpeg",
		CheckCmd:    "ffmpeg -version",
		Type:        "system",
		SizeMb:      50,
		Description: "Audio/video converter",
		InstallCmds: map[string]string{
			"apt":    "sudo apt-get install -y ffmpeg",
			"brew":  "brew install ffmpeg",
		},
	},
	"yt-dlp": {
		Name:        "yt-dlp",
		CheckCmd:    "yt-dlp --version",
		Type:        "pip",
		SizeMb:      10,
		Description: "YouTube video downloader",
		InstallCmds: map[string]string{
			"pip":   "pip install yt-dlp",
			"brew":  "brew install yt-dlp",
		},
	},

	// Vision tools
	"tesseract": {
		Name:        "tesseract",
		CheckCmd:    "tesseract --version",
		Type:        "system",
		SizeMb:      20,
		Description: "OCR engine",
		InstallCmds: map[string]string{
			"apt":    "sudo apt-get install -y tesseract-ocr",
			"brew":  "brew install tesseract",
		},
	},

	// Already existing
	"git": {
		Name:        "git",
		CheckCmd:    "git --version",
		Type:        "system",
		SizeMb:      50,
		Description: "Version control",
		InstallCmds: map[string]string{
			"apt":    "sudo apt-get update && sudo apt-get install -y git",
			"brew":  "brew install git",
			"winget": "winget install Git.Git -e",
		},
	},
	"chromium": {
		Name:        "chromium",
		CheckCmd:    "chromium --version",
		Type:        "system",
		SizeMb:      200,
		Description: "Browser runtime",
	},
}