package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ProjectNameForm struct {
	input        textinput.Model
	errorMessage string
}

func newProjectNameForm() *ProjectNameForm {
	ti := textinput.New()
	ti.Placeholder = "Enter your project name"
	ti.Focus()

	ti.CharLimit = 64
	ti.Width = 38
	ti.Prompt = ""

	return &ProjectNameForm{
		input:        ti,
		errorMessage: "",
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
			m.input = textinput.Model{}
			return m, tea.Quit

		case "enter":
			if m.isValidProjectName(m.input.Value()) {
				m.errorMessage = ""
				return m, tea.Quit
			} else {
				m.errorMessage = "Invalid project name. Use only letters, numbers, hyphens and underscores (no spaces)."
			}
		}
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *ProjectNameForm) View() string {
	title := titleStyle.Render("🚀 New Project Setup")
	inputBox := boxStyle.Render(m.input.View())

	result := fmt.Sprintf("%s \n%s", title, inputBox)

	if m.errorMessage != "" {
		errorBox := errorStyle.Render(m.errorMessage)
		result += fmt.Sprintf("\n%s", errorBox)
	}

	return result
}

func (m *ProjectNameForm) GetAppName() string {
	return m.input.Value()
}

// isValidProjectName validates that the project name is suitable for go mod init
func (m *ProjectNameForm) isValidProjectName(name string) bool {
	if strings.TrimSpace(name) == "" {
		return false
	}

	if strings.Contains(name, " ") {
		return false
	}

	validName := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return validName.MatchString(name)
}
