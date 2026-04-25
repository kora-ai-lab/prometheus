package capabilities

import (
	"os/exec"
	"strings"
)

type Installer struct{}

func NewInstaller() *Installer { return &Installer{} }

var toolInstallCommands = map[string]struct {
	installCmds []string
	checkCmd    string
}{
	"git":           {installCmds: []string{"apt-get install -y git", "brew install git", "apk add git"}, checkCmd: "git --version"},
	"python3":       {installCmds: []string{"apt-get install -y python3", "brew install python3", "apk add python3"}, checkCmd: "python3 --version"},
	"npm":           {installCmds: []string{"apt-get install -y npm", "brew install npm", "apk add npm"}, checkCmd: "npm --version"},
	"node":          {installCmds: []string{"apt-get install -y nodejs", "brew install node", "apk add nodejs"}, checkCmd: "node --version"},
	"poppler-utils": {installCmds: []string{"apt-get install -y poppler-utils", "apk add poppler-utils"}, checkCmd: "pdftotext -v"},
	"playwright":    {installCmds: []string{"pip install playwright", "npm install -g playwright"}, checkCmd: "playwright --version"},
	"ffmpeg":        {installCmds: []string{"apt-get install -y ffmpeg", "brew install ffmpeg", "apk add ffmpeg"}, checkCmd: "ffmpeg -version"},
	"upx":           {installCmds: []string{"apt-get install -y upx", "brew install upx", "apk add upx"}, checkCmd: "upx --version"},
}

func (i *Installer) IsAvailable(tool string) bool {
	path, err := exec.LookPath(tool)
	return err == nil && path != ""
}

func (i *Installer) GetInstallCommands(tool string) []string {
	if entry, ok := toolInstallCommands[tool]; ok {
		return entry.installCmds
	}
	return nil
}

func (i *Installer) Install(tool string) (string, error) {
	cmds := i.GetInstallCommands(tool)
	if len(cmds) == 0 {
		return "", nil
	}

	var lastErr error
	for _, cmd := range cmds {
		parts := strings.Fields(cmd)
		if len(parts) < 2 {
			continue
		}
		out := exec.Command(parts[0], parts[1:]...)
		err := out.Run()
		if err == nil {
			if i.IsAvailable(tool) {
				return "Installed " + tool, nil
			}
		} else {
			lastErr = err
		}
	}
	if lastErr != nil {
		return "", lastErr
	}
	return "", nil
}

func (i *Installer) NeedsInstall(tool string) bool {
	if i.IsAvailable(tool) {
		return false
	}
	return i.GetInstallCommands(tool) != nil
}

func (i *Installer) SuggestInstall(neededTools []string) []string {
	var toInstall []string
	for _, tool := range neededTools {
		if i.NeedsInstall(tool) {
			toInstall = append(toInstall, tool)
		}
	}
	return toInstall
}

func (i *Installer) PromptMessage(tools []string) string {
	if len(tools) == 0 {
		return ""
	}
	return "I need: " + strings.Join(tools, ", ") + ". Install? [y/n]"
}