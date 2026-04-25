package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type EnvironmentProfile struct {
	OS            string
	Arch          string
	CPUs          int
	Tools         []string
	OllamaModels  []string
	InternetReady bool
}

func DiscoverEnvironment() *EnvironmentProfile {
	profile := &EnvironmentProfile{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
		CPUs: runtime.NumCPU(),
	}

	// Tools to check
	tools := []string{"git", "python", "node", "npm", "docker", "curl", "wget", "gcc", "make"}
	for _, tool := range tools {
		if _, err := exec.LookPath(tool); err == nil {
			profile.Tools = append(profile.Tools, tool)
		}
	}

	// Check internet
	cmd := exec.Command("ping", "-n", "1", "8.8.8.8")
	if err := cmd.Run(); err == nil {
		profile.InternetReady = true
	}

	// Ollama models (simplified)
	out, err := exec.Command("ollama", "list").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines[1:] { // skip header
			fields := strings.Fields(line)
			if len(fields) > 0 {
				profile.OllamaModels = append(profile.OllamaModels, fields[0])
			}
		}
	}

	return profile
}

func (p *EnvironmentProfile) String() string {
	return fmt.Sprintf(
		"OS: %s, Arch: %s, CPUs: %d, Tools: [%s], Internet: %v, Models: [%s]",
		p.OS, p.Arch, p.CPUs, strings.Join(p.Tools, ", "), p.InternetReady, strings.Join(p.OllamaModels, ", "),
	)
}
