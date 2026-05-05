package ButtonCreateStartStop

import (
	commonStyles "goLogMyTime/components/common/styles"
	"goLogMyTime/domain/taskStore"
	"goLogMyTime/domain/types"
	"goLogMyTime/utils/mouse"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

type Model struct {
	DataSource *taskStore.TaskStore
	IdPrefix   string
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.MouseClickMsg:
		// TODO: get task definition from the store data
		taskDefinition := m.DataSource.GetInputTaskDefinition()
		buttonId := m.getButtonId(taskDefinition)

		if mouse.IsLeftClickInBoundsOf(buttonId, msg) {
			if taskDefinition != nil {
				m.DataSource.TaskToggle(taskDefinition)
			} else if m.DataSource.IsInputTaskValidAndNew() {
				m.DataSource.TaskCreate()
			}
		}
	}

	return m, cmd
}

func (m Model) View() tea.View {
	taskDefinition := m.DataSource.GetInputTaskDefinition()

	button := Button{
		id:    m.getButtonId(taskDefinition),
		label: "<- input",
		style: commonStyles.ButtonStyleDisabled,
	}

	if taskDefinition != nil {
		if m.DataSource.IsTaskStarted(taskDefinition) {
			button.label = "stop"
			button.style = commonStyles.ButtonStyleStarted
		} else {
			button.label = "start"
			button.style = commonStyles.ButtonStyleNotStarted
		}
	} else if m.DataSource.IsInputTaskValidAndNew() {
		button.label = "create"
		button.style = commonStyles.ButtonStyleNotStarted
	}

	button.style.Width(commonStyles.ButtonWidth)

	return tea.NewView(zone.Mark(
		button.id,
		button.style.Render(button.label),
	))
}

func (m Model) getButtonId(taskDefinition *types.TaskDefinition) string {
	if taskDefinition == nil {
		return "CREATE_NEW_TASK_DEFINITION"
	}

	return m.IdPrefix + taskDefinition.Ticket
}

type Button struct {
	id    string
	label string
	style lipgloss.Style
}
