package utils

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/orayew2002/goarm/domain"
)

const sdfListHeight = 6

// DatabaseItem represents a database option.
type DatabaseItem string

func (i DatabaseItem) FilterValue() string { return "" }

// Delegate rendering each item
type sdfItemDelegate struct{}

func (d sdfItemDelegate) Height() int                             { return 1 }
func (d sdfItemDelegate) Spacing() int                            { return 0 }
func (d sdfItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d sdfItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(DatabaseItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	if index == m.Index() {
		fmt.Fprint(w, sdfSelectedItemStyle.Render("➤ "+str))
		return
	}

	fmt.Fprint(w, sdfItemStyle.Render("  "+str))
}

type DatabaseSelectForm struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m *DatabaseSelectForm) Init() tea.Cmd {
	return nil
}

func (m *DatabaseSelectForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if i, ok := m.list.SelectedItem().(DatabaseItem); ok {
				m.choice = string(i)
			}

			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *DatabaseSelectForm) View() string {
	title := titleStyle.Render("🚀 Choose database")
	return fmt.Sprintf("%s \n%s", title, m.list.View())
}

func (m DatabaseSelectForm) GetChoice() string {
	return m.choice
}

func newDatabaseSelectForm() DatabaseSelectForm {
	items := make([]list.Item, len(domain.SupportedDatabaseTypes))
	for index, _ := range domain.SupportedDatabaseTypes {
		items[index] = DatabaseItem(domain.SupportedDatabaseTypes[index])
	}

	const defaultWidth = 40

	l := list.New(items, sdfItemDelegate{}, defaultWidth, sdfListHeight)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	l.SetShowTitle(false)

	return DatabaseSelectForm{
		list: l,
	}
}
