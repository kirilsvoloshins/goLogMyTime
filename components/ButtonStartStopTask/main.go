package ButtonStartStop

import (
	"goLogMyTime/domain/taskStore"
	"goLogMyTime/domain/types"
	"goLogMyTime/utils/assert"
	"goLogMyTime/utils/mouse"

	tea "charm.land/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

type Model struct {
	DataSource     *taskStore.TaskStore
	TaskDefinition *types.TaskDefinition
	IdPrefix       string
}

func (m Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.MouseClickMsg:
		assert.Assert(m.TaskDefinition != nil, "missing task def")

		if mouse.IsLeftClickInBoundsOf(m.getButtonId(), msg) {
			m.DataSource.TaskToggle(m.TaskDefinition)
		}
	}

	return nil
}

func (m Model) View() tea.View {
	active := m.DataSource.IsTaskStarted(m.TaskDefinition)
	buttonLabel := m.getButtonLabel(active)
	buttonStyle := m.DataSource.Style.GetStartStopButtonStyle(active)

	return tea.NewView(zone.Mark(
		m.getButtonId(),
		buttonStyle.Render(buttonLabel),
	))
}

func (m Model) getButtonId() string {
	return m.IdPrefix + m.TaskDefinition.Ticket
}

func (m Model) getButtonLabel(active bool) string {
	if active {
		return "stop"
	}

	return "start"
}
