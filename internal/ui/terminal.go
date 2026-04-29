package ui

import (
	"fmt"
	"os"
)

func IsInteractive() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func printBanner() {
	fmt.Print(`
╔═════════════════════════════════════╗
║     Prometheus v1.0.2                 ║
║     AI-first agent runtime             ║
╚═════════════════════════════════════╝
`)
}