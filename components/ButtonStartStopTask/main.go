package ButtonStartStop

import (
	"goLogMyTime/store"
	"goLogMyTime/types"
	"goLogMyTime/utils"

	tea "charm.land/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

type Model struct {
	Store    *store.RootStore
	Ticket   string
	IdPrefix string
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.MouseClickMsg:
		if utils.IsLeftClickInBoundsOf(m.getButtonZone(), msg) {
			m.Store.Task.TaskToggle(m.Ticket)
		}
	}

	return nil
}

func (m *Model) View() tea.View {
	isActive := m.Store.Task.IsTaskStarted(m.Ticket)
	s := m.Store.Ui.Style.GetStartStopButtonStyle(isActive)

	return tea.NewView(zone.Mark(
		string(m.getButtonZone()),
		s.Render(
			m.getButtonLabel(isActive),
		),
	))
}

func (m *Model) getButtonZone() types.Zone {
	return types.Zone(m.IdPrefix + m.Ticket)
}

func (m *Model) getButtonLabel(active bool) string {
	if active {
		return "stop"
	}

	return "start"
}
