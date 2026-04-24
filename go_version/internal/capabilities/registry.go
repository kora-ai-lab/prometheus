package capabilities

type Capability struct {
	Name        string
	CheckCmd    string
	InstallCmds map[string]string
	Type        string
	SizeMb      int
	Description string
}

var Registry = map[string]Capability{
	"git": {
		Name:        "git",
		CheckCmd:    "git --version",
		Type:        "system",
		SizeMb:      50,
		Description: "Version control",
		InstallCmds: map[string]string{
			"apt":    "sudo apt-get update && sudo apt-get install -y git",
			"brew":   "brew install git",
			"winget": "winget install --id Git.Git -e",
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
