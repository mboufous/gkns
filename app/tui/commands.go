package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type namespacesReadyMsg struct {
	namespaces       []list.Item
	currentContext   string
	currentNamespace string
}

// errMsg holds the error msg
type errMsg struct {
	err error
}

func (e errMsg) Error() string {
	return e.err.Error()
}

// loadNamespaces loads namespace from the current context
func (m Model) loadNamespaces() tea.Msg {
	//time.Sleep(2 * time.Second)
	namespaces, err := m.client.Namespaces()
	if err != nil {
		return errMsg{
			err: fmt.Errorf("error getting namespaces. %w\n", err),
		}
	}

	items := make([]list.Item, len(namespaces))
	for i, ns := range namespaces {
		items[i] = item{
			title:        ns.Name,
			creationTime: ns.CreationTimestamp.String(),
		}
	}

	currentContext, err := m.client.CurrentContext()
	if err != nil {
		return errMsg{
			err: fmt.Errorf("error getting current context. %w\n", err),
		}
	}

	currentNamespace, err := m.client.CurrentNamespace()
	if err != nil {
		return errMsg{
			err: fmt.Errorf("error getting current namespace. %w\n", err),
		}
	}

	return namespacesReadyMsg{
		namespaces:       items,
		currentContext:   currentContext,
		currentNamespace: currentNamespace,
	}
}

// SwitchNamespace switches current context's namespace
func (m Model) SwitchNamespace() {
	err := m.client.SwitchNamespace(m.currentNamespace)
	if err != nil {
		panic(err)
	}
}
