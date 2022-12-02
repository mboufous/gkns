package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mboufous/gkns/app"
)

var (
	exitKeys = key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "Quit internal.app"))

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusBarItemStyle = statusBarStyle.Copy().
				Foreground(lipgloss.Color("#FFFDF5")).
				Background(lipgloss.Color("#255ea0")).
				Padding(0, 1)
	currentContextStyle = statusBarItemStyle.Copy().
				MarginLeft(0)

	currentNamespaceStyle = statusBarItemStyle.Copy().
				MarginRight(0).
				Align(lipgloss.Right)
	listViewStyle = lipgloss.NewStyle().
			Padding(1, 0)
)

type Model struct {
	client           *app.ClientProvider
	list             list.Model
	statusMessage    error
	currentStatus    string
	frameWidth       int
	currentContext   string
	currentNamespace string
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadNamespaces,
		m.list.StartSpinner(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := listViewStyle.GetFrameSize()
		statusBarHeight := lipgloss.Height(m.statusView())
		m.list.SetSize(msg.Width-h, msg.Height-v-statusBarHeight)
		m.frameWidth = msg.Width
	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}
		if key.Matches(msg, exitKeys) {
			return m, tea.Quit
		}
	case errMsg:
		m.statusMessage = msg.err
		return m, nil
	case namespacesReadyMsg:
		m.list.StopSpinner()
		m.list.SetItems(msg.namespaces)
		m.currentContext = msg.currentContext
		m.currentNamespace = msg.currentNamespace
	case namespaceChangedMsg:
		var oldNamespace string
		oldNamespace, m.currentNamespace = m.currentNamespace, msg.ns
		if m.currentNamespace != oldNamespace {
			m.SwitchNamespace()
		}

	}

	return m, cmd
}

func (m Model) View() string {
	if m.statusMessage != nil {
		return fmt.Sprintf("\n Error: %v\n", m.statusMessage)
	}
	return lipgloss.JoinVertical(lipgloss.Top, listViewStyle.Render(m.list.View()), m.statusView())
}

func (m Model) statusView() string {
	var currentNamespace string
	var currentContext string

	if len(m.list.VisibleItems()) == 0 {
		currentNamespace = "Loading ..."
		currentContext = "Loading ..."
	} else {
		currentNamespace = "⚡ Namespace: " + m.currentNamespace
		currentContext = "⚡ Context: " + m.currentContext
	}

	currentContextView := currentContextStyle.Render(currentContext)
	currentNamespaceView := currentNamespaceStyle.Render(currentNamespace)
	spaceHolderView := statusBarStyle.Copy().
		Width(m.frameWidth - lipgloss.Width(currentContextView) - lipgloss.Width(currentNamespaceView)).
		Render("")

	bar := lipgloss.JoinHorizontal(lipgloss.Left, currentContextView, spaceHolderView, currentNamespaceView)

	return statusBarStyle.Width(m.frameWidth).Render(bar)
}

// New creates new bubble tea model
func New(client *app.ClientProvider) (*Model, error) {
	return &Model{
		client: client,
		list:   NewList(),
	}, nil

}
