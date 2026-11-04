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

func (m Model) Update(msg tea.Msg) tea.Cmd {
	mouseClickMsg, isMouseClick := msg.(tea.MouseClickMsg)

	if isMouseClick && utils.IsLeftClickInBoundsOf(m.getButtonZone(), mouseClickMsg) {
		m.Store.Task.TaskToggle(m.Ticket)
	}

	return nil
}

func (m Model) View() tea.View {
	isActive := m.Store.Task.IsTaskStarted(m.Ticket)
	buttonStyle := m.Store.Ui.Style.GetStartStopButtonStyle(isActive)
	buttonLabel := m.getButtonLabel(isActive)

	return tea.NewView(zone.Mark(
		string(m.getButtonZone()),
		buttonStyle.Render(buttonLabel),
	))
}

func (m Model) getButtonZone() types.Zone {
	return types.Zone(m.IdPrefix + m.Ticket)
}

func (m Model) getButtonLabel(active bool) string {
	if active {
		return "stop"
	}

	return "start"
}
