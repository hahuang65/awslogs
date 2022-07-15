package list

import (
	"git.sr.ht/~hwrd/awslogs/internal/aws"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	statusMessageStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#036B46", Dark: "#036B46"}).
		Render
)

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var (
			pi item
		)

		if i, ok := m.SelectedItem().(item); ok {
			pi = i
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.open):
				return tea.Batch(
					m.NewStatusMessage(statusMessageStyle("Opening "+pi.Title())),
					aws.LoadLogStreams(pi.logGroup),
				)

			case key.Matches(msg, keys.refresh):
				return tea.Batch(
					m.NewStatusMessage(statusMessageStyle("Refreshing Log Groups")),
					aws.ListLogGroups,
				)
			}
		}

		return nil
	}

	help := []key.Binding{keys.open, keys.refresh}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type delegateKeyMap struct {
	open    key.Binding
	refresh key.Binding
}

func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.open,
		d.refresh,
	}
}

func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.open,
			d.refresh,
		},
	}
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		open: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "view logs"),
		),
		refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
	}
}
