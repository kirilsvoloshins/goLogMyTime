package ButtonStartStop

import (
	commonStyles "goLogMyTime/components/common/styles"
	"goLogMyTime/domain/taskStore"
	"goLogMyTime/domain/types"
	"goLogMyTime/utils/assert"
	"goLogMyTime/utils/mouse"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
	button := m.getButton()

	return tea.NewView(zone.Mark(
		m.getButtonId(),
		button.style.Render(button.label),
	))
}

func (m Model) getButtonId() string {
	return m.IdPrefix + m.TaskDefinition.Ticket
}

func (m Model) getButton() Button {
	if m.DataSource.IsTaskStarted(m.TaskDefinition) {
		return Button{
			label: "stop",
			style: commonStyles.ButtonStyleStarted,
		}
	}

	return Button{
		label: "start",
		style: commonStyles.ButtonStyleNotStarted,
	}
}

type Button struct {
	label string
	style lipgloss.Style
}
