package utils

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ProjectNameForm struct {
	input textinput.Model
}

func newProjectNameForm() *ProjectNameForm {
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
			m.input = textinput.Model{}
			return m, tea.Quit

		case "enter":
			return m, tea.Quit
		}
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *ProjectNameForm) View() string {
	title := titleStyle.Render("🚀 New Project Setup")
	inputBox := boxStyle.Render(m.input.View())

	return fmt.Sprintf("%s \n%s", title, inputBox)
}

func (m *ProjectNameForm) GetAppName() string {
	return m.input.Value()
}
