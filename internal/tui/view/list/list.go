package list

import (
	"time"

	"git.sr.ht/~hwrd/awslogs/internal/aws"
	"git.sr.ht/~hwrd/awslogs/internal/tui/view"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	listStyle  = lipgloss.NewStyle().Padding(1, 2)
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#036B46")).
			Padding(0, 1)
)

type model struct {
	list         list.Model
	delegateKeys *delegateKeyMap
}

type item struct {
	logGroup aws.LogGroup
}

func New() model {
	var (
		delegateKeys = newDelegateKeyMap()
	)

	delegate := newItemDelegate(delegateKeys)
	logGroupList := list.New([]list.Item{}, delegate, 0, 0)
	logGroupList.Title = "AWS Logs"
	logGroupList.Styles.Title = titleStyle
	logGroupList.StatusMessageLifetime = time.Second * 5

	return model{
		list:         logGroupList,
		delegateKeys: delegateKeys,
	}
}

func (m model) Init() tea.Cmd {
	return aws.ListLogGroups
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		topGap, rightGap, bottomGap, leftGap := listStyle.GetPadding()
		m.list.SetSize(msg.Width-leftGap-rightGap, msg.Height-topGap-bottomGap)

	case aws.LogGroupsMsg:
		m.list.SetItems(itemize(msg))
		cmds = append(cmds, view.SetView(view.List))

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return listStyle.Render(m.list.View())
}

func itemize(lgs []aws.LogGroup) []list.Item {
	list_items := []list.Item{}

	for _, lg := range lgs {
		list_items = append(list_items, item{logGroup: lg})
	}

	return list_items
}

func (i item) Title() string {
	return i.logGroup.Name
}

func (i item) Description() string {
	return ""
}

func (i item) FilterValue() string {
	return i.Title()
}
