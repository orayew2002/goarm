package utils

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/MH-KodaCore/goarm/domain"
)

const sffListHeight = 4

// FrameworkItem represents a framework option.
type FrameworkItem string

func (i FrameworkItem) FilterValue() string { return "" }

// Delegate rendering each item
type sffItemDelegate struct{}

func (d sffItemDelegate) Height() int                             { return 1 }
func (d sffItemDelegate) Spacing() int                            { return 0 }
func (d sffItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d sffItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(FrameworkItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	if index == m.Index() {
		fmt.Fprint(w, sffSelectedItemStyle.Render("➤ "+str))
		return
	}

	fmt.Fprint(w, sffItemStyle.Render("  "+str))
}

type FrameworkSelectForm struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m *FrameworkSelectForm) Init() tea.Cmd {
	return nil
}

func (m *FrameworkSelectForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			if i, ok := m.list.SelectedItem().(FrameworkItem); ok {
				m.choice = string(i)
			}

			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *FrameworkSelectForm) View() string {
	title := titleStyle.Render("🚀 Choose web framework")
	return fmt.Sprintf("%s \n%s", title, m.list.View())
}

func (m FrameworkSelectForm) GetChoice() string {
	return m.choice
}

func newFrameworkSelectForm() FrameworkSelectForm {
	items := make([]list.Item, len(domain.SupportedFrameworkTypes))
	for index := range domain.SupportedFrameworkTypes {
		items[index] = FrameworkItem(domain.SupportedFrameworkTypes[index])
	}

	const defaultWidth = 40

	l := list.New(items, sffItemDelegate{}, defaultWidth, sffListHeight)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	l.SetShowTitle(false)

	return FrameworkSelectForm{
		list: l,
	}
}
