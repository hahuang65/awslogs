package tui

import (
	"fmt"

	"git.sr.ht/~hwrd/awslogs/internal/aws"
	"git.sr.ht/~hwrd/awslogs/internal/tui/view"
	"git.sr.ht/~hwrd/awslogs/internal/tui/view/list"
	"git.sr.ht/~hwrd/awslogs/internal/tui/view/stream"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	service     aws.Service
	listView    tea.Model
	streamView  tea.Model
	spinner     spinner.Model
	currentView view.View
}

func newModel(service aws.Service) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		service:     service,
		listView:    list.New(),
		streamView:  stream.New(),
		spinner:     s,
		currentView: view.Spinner,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.listView.Init(),
		m.streamView.Init(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		if m.currentView == view.Spinner {
			newSpinner, cmd := m.spinner.Update(msg)
			m.spinner = newSpinner
			cmds = append(cmds, cmd)
		}

	case view.SetMsg:
		m.currentView = view.View(msg)

	case aws.ListLogGroupsMsg:
		cmds = append(cmds, m.service.ListLogGroups)

	case aws.LoadLogStreamsMsg:
		cmds = append(cmds, m.service.LoadLogStreams(aws.LogGroup(msg)))
	}

	// Only update the sub-models if it's the currently focused one, or the msg is not a keypress
	// This prevents keypresses in one sub-model from triggering actions in another sub-model
	_, isKeyMsg := msg.(tea.KeyMsg)

	if !isKeyMsg || m.currentView == view.List {
		newListView, cmd := m.listView.Update(msg)
		m.listView = newListView
		cmds = append(cmds, cmd)
	}

	if !isKeyMsg || m.currentView == view.Stream {
		newStreamView, cmd := m.streamView.Update(msg)
		m.streamView = newStreamView
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.currentView == view.Spinner {
		return fmt.Sprintf("\n\n   %s Loading logs\n\n", m.spinner.View())
	} else if m.currentView == view.List {
		return m.listView.View()
	} else if m.currentView == view.Stream {
		return m.streamView.View()
	} else {
		return ""
	}
}

func Start(service aws.Service) error {
	p := tea.NewProgram(
		newModel(service),
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	return p.Start()
}
