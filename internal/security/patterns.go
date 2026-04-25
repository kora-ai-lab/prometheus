package security

import (
	"strings"
)

type Pattern struct {
	Needle string
	Score  int
}

var dangerousPatterns = []Pattern{
	{Needle: "rm -rf /", Score: 100},
	{Needle: "curl", Score: 40},
	{Needle: "| sh", Score: 80},
	{Needle: "chmod 777 /", Score: 100},
	{Needle: "/etc/passwd", Score: 95},
	{Needle: ":(){ :|:& };:", Score: 100},
	{Needle: "sudo", Score: 30},
	{Needle: "chown", Score: 30},
	{Needle: "wget", Score: 40},
	{Needle: "mkfs", Score: 100},
	{Needle: "dd if=", Score: 60},
	{Needle: "> /dev/sda", Score: 100},
}

func calculateScore(command string) int {
	score := 0
	lower := strings.ToLower(command)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lower, strings.ToLower(pattern.Needle)) {
			score += pattern.Score
		}
	}
	return score
}