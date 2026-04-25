package update

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

var currentBinary string

func init() {
	currentBinary, _ = os.Executable()
}

func ApplyUpdate(tmpPath string) error {
	if _, err := os.Stat(tmpPath); err != nil {
		return fmt.Errorf("update binary not found: %w", err)
	}

	backupPath := currentBinary + ".backup"
	if _, err := os.Stat(currentBinary); err == nil {
		if err := copyFile(backupPath, currentBinary); err != nil {
			return fmt.Errorf("create backup: %w", err)
		}
	}

	if err := copyFile(currentBinary, tmpPath); err != nil {
		if _, err := os.Stat(backupPath); err == nil {
			os.Rename(backupPath, currentBinary)
		}
		return fmt.Errorf("replace binary: %w", err)
	}

	if err := os.Chmod(currentBinary, 0755); err != nil {
		return fmt.Errorf("chmod: %w", err)
	}

	return restart()
}

func restart() error {
	args := os.Args
	cmd := exec.Command(currentBinary, args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

func copyFile(dst, src string) error {
	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	return err
}