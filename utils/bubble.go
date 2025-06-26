package utils

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00BFFF")).
			Bold(true).
			PaddingTop(1).
			PaddingBottom(1).
			Align(lipgloss.Center)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00BFFF")).
			Padding(0, 1).
			Width(40).
			Align(lipgloss.Left)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#999999")).
			MarginTop(1).
			Align(lipgloss.Center)
)

type ProjectNameForm struct {
	input textinput.Model
}

func NewProjectNameForm() *ProjectNameForm {
	ti := textinput.New()
	ti.Placeholder = "Enter your project name"
	ti.Focus()
	ti.CharLimit = 64
	ti.Width = 38
	ti.Prompt = ""

	return &ProjectNameForm{
		input: ti,
	}
}

func (m *ProjectNameForm) Init() tea.Cmd {
	return textinput.Blink
}

func (m *ProjectNameForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			fmt.Println("✅ Project Name:", m.input.Value())
			return m, tea.Quit
		}
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *ProjectNameForm) View() string {
	title := titleStyle.Render("🚀 New Project Setup")
	inputBox := boxStyle.Render(m.input.View())
	footer := footerStyle.Render("Press Enter to confirm • Esc or Ctrl+C to cancel")

	return fmt.Sprintf("\n%s \n%s \n%s", title, inputBox, footer)
}

func (m *ProjectNameForm) GetAppName() string {
	return m.input.Value()
}
