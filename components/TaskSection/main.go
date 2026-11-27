package TaskSection

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

type Model struct {
	Store     *store.RootStore
	Zone      types.Zone
	Shortcut  types.SectionShortcut
	Title     string
	Cursor    *int
	TaskDefs  *types.TaskDefinitions
	Width     int
	IsFocused bool
	taskLines []TaskLine.Model
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	m.taskLines = []TaskLine.Model{}

	maxTicketLength := styles.GetLongestTicketLen(*m.TaskDefs)

	for i, taskDef := range *m.TaskDefs {
		taskPadding := maxTicketLength - len(taskDef.Ticket)
		taskLine := TaskLine.Model{
			Store:       m.Store,
			ZoneParent:  m.Zone,
			IsSelected:  m.IsFocused && i == *m.Cursor,
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

			if isToggleKey && len(*m.TaskDefs) > 0 {
				taskDef := (*m.TaskDefs)[*m.Cursor]
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
			lipgloss.Sprintf(m.Title+" [%s]", m.Shortcut),
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
