// +build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	user32             = syscall.NewLazyDLL("user32.dll")
	procGetConsoleWindow = kernel32.NewProc("GetConsoleWindow")
	procShowWindow       = user32.NewProc("ShowWindow")
)

const (
	SW_HIDE = 0
	SW_SHOW = 5
)

func hideConsole() {
	hwnd, _, _ := procGetConsoleWindow.Call()
	if hwnd != 0 {
		procShowWindow.Call(hwnd, SW_HIDE)
	}
}

func showConsole() {
	hwnd, _, _ := procGetConsoleWindow.Call()
	if hwnd != 0 {
		procShowWindow.Call(hwnd, SW_SHOW)
	}
}

func main() {
	// Hide console window for desktop mode
	hideConsole()

	// Get the directory of this executable
	exePath, err := os.Executable()
	if err != nil {
		showErrorAndExit("Failed to get executable path: " + err.Error())
		return
	}
	exeDir := getDir(exePath)

	// Path to the main prometheus.exe
	prometheusExe := exeDir + "\\prometheus.exe"

	// Check if prometheus.exe exists
	if _, err := os.Stat(prometheusExe); os.IsNotExist(err) {
		showErrorAndExit("prometheus.exe not found in: " + exeDir)
		return
	}

	// Launch prometheus.exe with --web flag
	cmd := exec.Command(prometheusExe, "--web")
	cmd.Dir = exeDir

	// Start the process
	if err := cmd.Start(); err != nil {
		showErrorAndExit("Failed to start Prometheus: " + err.Error())
		return
	}

	// Wait a moment for the server to start
	time.Sleep(2 * time.Second)

	// Open browser to localhost:8080
	openBrowser("http://localhost:8080")

	// The prometheus process runs in background
	// We could add a system tray icon here for control
	// For now, just exit - the prometheus process continues independently
}

func getDir(path string) string {
	// Remove the executable name
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '\\' || path[i] == '/' {
			return path[:i]
		}
	}
	return path
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	cmd.Start()
}

func showErrorAndExit(message string) {
	showConsole()
	fmt.Println("Error:", message)
	fmt.Println("\nPress Enter to exit...")
	fmt.Scanln()
	os.Exit(1)
}
