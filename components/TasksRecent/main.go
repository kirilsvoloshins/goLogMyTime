package TasksRecent

import (
	"goLogMyTime/common/styles"
	"goLogMyTime/components/TaskLine"
	"goLogMyTime/store"
	"goLogMyTime/types"
	"goLogMyTime/utils"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

const PREFIX = "R"

type Model struct {
	Store     *store.RootStore
	Width     int
	IsFocused bool
	Shortcut  types.SectionShortcut
	Zone      types.Zone
	taskLines []TaskLine.Model
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	m.taskLines = []TaskLine.Model{}

	taskDefs := m.Store.Task.TaskDefinitionsRecent
	maxTicketLength := styles.GetLongestTicketLen(taskDefs)

	for i, taskDef := range taskDefs {
		taskPadding := maxTicketLength - len(taskDef.Ticket)
		taskLine := TaskLine.Model{
			Store:       m.Store,
			ZoneParent:  m.Zone,
			IsSelected:  m.IsFocused && i == m.Store.Ui.TasksRecentCursor,
			TaskDef:     taskDef,
			TaskPadding: taskPadding,
		}
		taskLine.Update(msg)
		m.taskLines = append(m.taskLines, taskLine)
	}

	switch msg := msg.(type) {
	case types.FocusZoneMsg:
		if types.Zone(msg) == m.Zone {
			m.IsFocused = true
		}

	case types.BlurZoneMsg:
		if types.Zone(msg) == m.Zone {
			m.IsFocused = false
		}

	case tea.KeyPressMsg:
		if m.IsFocused {
			isToggleKey := utils.IsEnterKeyMsg(msg) || utils.IsSpaceKeyMsg(msg)

			if isToggleKey && len(m.Store.Task.TaskDefinitionsRecent) > 0 {
				taskDef := m.Store.Task.TaskDefinitionsRecent[m.Store.Ui.TasksRecentCursor]
				m.Store.Task.TaskToggle(taskDef.Ticket)
			}
		}
	}

	return m, nil
}

func (m *Model) View() tea.View {
	var style = m.Store.Ui.Style

	var sb strings.Builder
	sb.WriteString(
		style.SectionTitle.Render(
			lipgloss.Sprintf("Recent tasks [%s] ", m.Shortcut),
		),
	)
	for _, taskLine := range m.taskLines {
		sb.WriteString(taskLine.View().Content)
	}

	s := style.GetSectionStyle(m.IsFocused).Width(m.Width).Padding(0, 0)
	return tea.NewView(zone.Mark(
		string(m.Zone),
		s.Render(sb.String()),
	))
}
