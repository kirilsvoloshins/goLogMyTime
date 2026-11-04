package ButtonCreateStartStop

import (
	"goLogMyTime/common/globals"
	"goLogMyTime/store"
	"goLogMyTime/utils"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

type Model struct {
	Store *store.RootStore
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.MouseClickMsg:
		if utils.IsLeftClickInBoundsOf(globals.ZoneCreateStartStop, msg) {
			taskStateChange := m.Store.Task.TaskToggleIfValid(m.Store.Ui.InputTicket, m.Store.Ui.InputTitle)
			if taskStateChange == globals.TASK_GOT_STOPPED {
				cmd = globals.FocusZone(globals.ZoneTextInputTask)
			}
		}
	}

	return m, cmd
}

func (m *Model) View() tea.View {
	buttonLabel, buttonStyle := m.GetButtonStyle()

	return tea.NewView(zone.Mark(
		globals.ZoneCreateStartStop,
		buttonStyle.Render(buttonLabel),
	))
}

func (m *Model) GetButtonStyle() (string, lipgloss.Style) {
	ticket := m.Store.Ui.InputTicket
	taskDefinition := m.Store.Task.GetTaskDefinition(ticket)
	if taskDefinition != nil {
		if m.Store.Task.IsTaskStarted(ticket) {
			return "stop", m.Store.Ui.Style.ButtonActive
		}
		return "start", m.Store.Ui.Style.ButtonInactive
	}

	title := m.Store.Ui.InputTitle
	if m.Store.Task.IsTaskValidAndNew(ticket, title) {
		return "create and start", m.Store.Ui.Style.ButtonInactive
	}

	return "input task and title", m.Store.Ui.Style.ButtonDisabled
}
