//go:build !windows

package service

import "fmt"

func Install() error {
	return fmt.Errorf("not supported on this platform")
}

func Uninstall() error {
	return fmt.Errorf("not supported on this platform")
}

func Status() (string, error) {
	return "", fmt.Errorf("not supported on this platform")
}

func Start() error {
	return fmt.Errorf("not supported on this platform")
}

func Stop() error {
	return fmt.Errorf("not supported on this platform")
}

func Run() error {
	return fmt.Errorf("not supported on this platform")
}
