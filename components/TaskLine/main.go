package TaskLine

import (
	"goLogMyTime/common/styles"
	"goLogMyTime/components/ButtonFavorite"
	ButtonStartStop "goLogMyTime/components/ButtonStartStopTask"
	"goLogMyTime/store"
	"goLogMyTime/types"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Model struct {
	Store           *store.RootStore
	ZoneParent      types.Zone
	IsSelected      bool
	TaskDef         types.TaskDefinition
	TaskPadding     int
	buttonStartStop *ButtonStartStop.Model
	buttonFavorite  *ButtonFavorite.Model
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	// TODO: probably not performant
	m.buttonStartStop = &ButtonStartStop.Model{
		Store:    m.Store,
		Ticket:   m.TaskDef.Ticket,
		IdPrefix: string(m.ZoneParent) + "FS_",
	}
	m.buttonStartStop.Update(msg)

	m.buttonFavorite = &ButtonFavorite.Model{
		Store:    m.Store,
		Ticket:   m.TaskDef.Ticket,
		IdPrefix: string(m.ZoneParent) + "FF_",
	}
	m.buttonFavorite.Update(msg)

	return m, nil
}

func (m *Model) View() tea.View {
	var sb strings.Builder
	sb.WriteString("\n")
	if m.IsSelected {
		sb.WriteString("> ")
	}
	sb.WriteString(
		lipgloss.JoinHorizontal(
			lipgloss.Bottom,
			styles.GetTaskInfo(m.TaskDef, m.TaskPadding),
			m.buttonStartStop.View().Content,
			m.buttonFavorite.View().Content,
		),
	)

	s := lipgloss.NewStyle()

	return tea.NewView(
		s.Render(sb.String()),
	)
}
