package ButtonStartStop

import (
	commonStyles "goLogMyTime/components/common/styles"
	"goLogMyTime/domain/taskStore"
	"goLogMyTime/domain/types"
	"goLogMyTime/utils/assert"

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
		buttonId := m.getButtonId()

		if isMouseLeftInBoundsOf(buttonId, msg) {
			m.DataSource.TaskToggle(m.TaskDefinition)
		}
	}

	return nil
}

func (m Model) View() tea.View {
	assert.Assert(m.TaskDefinition != nil, "more nils more")

	button := Button{
		id:    m.getButtonId(),
		label: "start",
		style: commonStyles.ButtonStyleNotStarted,
	}

	if m.DataSource.IsTaskStarted(m.TaskDefinition) {
		button.label = "stop"
		button.style = commonStyles.ButtonStyleStarted
	}

	button.style.Width(commonStyles.ButtonWidth)

	return tea.NewView(
		zone.Mark(button.id, button.style.Render(button.label)),
	)
}

func (m Model) getButtonId() string {
	return m.IdPrefix + m.TaskDefinition.Ticket
}

func isMouseLeftInBoundsOf(buttonId string, msg tea.MouseClickMsg) bool {
	return msg.Button == tea.MouseLeft && zone.Get(buttonId).InBounds(msg)
}

type Button struct {
	id    string
	label string
	style lipgloss.Style
}
