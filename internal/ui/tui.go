package ui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true).
			Padding(0, 1)

	conversationStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")).
				Padding(0, 1)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("165")).
			Bold(true)

	statusDoneStyle_   = lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Bold(true)
	statusRunning_ = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
	statusBlocked_ = lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true)
	statusError_   = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)

	indicatorDone_   = "✓"
	indicatorRun_   = "⟳"
	indicatorBlock_ = "⊙"
	indicatorWarn_  = "⚠"
	indicatorError_ = "✗"
	indicatorVision_  = "👁"
	indicatorBrowser_ = "🌐"
)

type Model struct {
	width        int
	height      int
	messages     []Message
	input       string
	term        string
	status      string
	goal        string
	lastOutput  string
	showLogs    bool
	blocking    bool
	blockReason string
}

type Message struct {
	Role    string
	Content string
}

func NewTUI() *Model {
	m := &Model{
		messages: []Message{},
		status:  "running",
		term:    os.Getenv("TERM"),
	}
	return m
}

func (m *Model) SetGoal(goal string) {
	m.goal = goal
}

func (m *Model) AddMessage(role, content string) {
	m.messages = append(m.messages, Message{Role: role, Content: content})
}

func (m *Model) SetStatus(status string) {
	m.status = status
}

func (m *Model) SetLastOutput(output string) {
	m.lastOutput = output
}

func (m *Model) SetBlocking(reason string) {
	m.blocking = true
	m.blockReason = reason
}

func (m *Model) ToggleLogs() {
	m.showLogs = !m.showLogs
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+l":
			m.showLogs = !m.showLogs
		case "enter":
			if m.blocking {
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m *Model) View() string {
	if m.term == "dumb" {
		return m.plainView()
	}
	return m.terminalView()
}

func (m *Model) plainView() string {
	var b strings.Builder
	b.WriteString("PROMETHEUS\n")
	b.WriteString(strings.Repeat("-", 40))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Goal: %s\n", m.goal))
	b.WriteString(fmt.Sprintf("Status: %s\n", m.status))
	for _, msg := range m.messages {
		b.WriteString(fmt.Sprintf("%s: %s\n", msg.Role, msg.Content))
	}
	b.WriteString("> ")
	return b.String()
}

func (m *Model) terminalView() string {
	var b strings.Builder

	header := fmt.Sprintf("⟪ PROMETHEUS — %s ⟫", m.statusIndicator())
	if m.blocking {
		header = fmt.Sprintf("⊙ BLOqué — %s", m.blockReason)
	}
	b.WriteString(headerStyle.Render(header))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", max(0, m.width)))
	b.WriteString("\n")

	goalLen := min(len(m.goal), m.width-10)
	b.WriteString(fmt.Sprintf("Goal: %s\n", m.goal[:goalLen]))
	b.WriteString("\n")

	msgStart := max(0, len(m.messages)-m.height+8)
	visible := m.messages[msgStart:]
	for _, msg := range visible {
		roleStyle := conversationStyle
		if msg.Role == "assistant" {
			roleStyle = conversationStyle.Foreground(lipgloss.Color("229"))
		}
		content := msg.Content
		if len(content) > m.width-4 {
			content = content[:m.width-7] + "..."
		}
		b.WriteString(roleStyle.Render(fmt.Sprintf("%s: %s", msg.Role, content)))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(inputStyle.Render("> "))
	b.WriteString(m.input)
	return b.String()
}

func (m *Model) statusIndicator() string {
	switch m.status {
	case "done":
		return "✓ Terminé"
	case "running":
		return "⟳ En cours"
	case "blocked":
		return "⊙ Bloqué"
	case "failed":
		return "✗ Échec"
	default:
		return m.status
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (m *Model) Run() error {
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}