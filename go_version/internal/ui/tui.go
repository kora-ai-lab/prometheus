package ui

import "fmt"

type TUI struct{}

func New() *TUI {
	return &TUI{}
}

func (t *TUI) Println(msg string) {
	fmt.Println(msg)
}
