package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#255ea0")).
			Padding(0, 1)
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"})
	selectKey = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "choose namespace"),
	)
)

type item struct {
	title        string
	creationTime string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.creationTime }
func (i item) FilterValue() string { return i.title }

func NewList() list.Model {
	l := list.New(nil, namespaceSwitcherDelegate(), 0, 0)
	l.Styles.Title = titleStyle
	l.Title = "📦 Namespaces"
	l.SetShowStatusBar(true)
	l.StartSpinner()
	return l
}

type namespaceChangedMsg struct {
	ns string
}

func namespaceSwitcherDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var selectedNamespace string
		item, ok := m.SelectedItem().(item)
		if !ok {
			return nil
		}

		selectedNamespace = item.Title()
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, selectKey):
				return tea.Batch(m.NewStatusMessage(statusMessageStyle.Render("✓ Set namespace to "+selectedNamespace)),
					func() tea.Msg {
						return namespaceChangedMsg{
							ns: selectedNamespace,
						}
					})
			}
		}
		return nil
	}

	help := []key.Binding{selectKey}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}
