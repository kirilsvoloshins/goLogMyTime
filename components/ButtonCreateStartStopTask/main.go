package ButtonCreateStartStop

import (
	"goLogMyTime/components/common/globals"
	commonStyles "goLogMyTime/components/common/styles"
	"goLogMyTime/domain/taskStore"
	"goLogMyTime/utils/assert"
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
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.MouseClickMsg:
		if mouse.IsLeftClickInBoundsOf(globals.CreateStartStopZone, msg) {
			taskDefinition := m.DataSource.GetInputTaskDefinition()

			if taskDefinition != nil {
				m.DataSource.TaskToggle(taskDefinition)
				return m, tea.Batch(cmds...)
			}

			if m.DataSource.IsInputTaskValidAndNew() {
				m.DataSource.TaskCreate()
				return m, tea.Batch(cmds...)
			}
		}

	}

	return m, nil
}

func (m Model) View() tea.View {
	button, ok := buttonMap[m.DataSource.TaskInputState]
	assert.Assert(ok, "Invalid task input state")

	return tea.NewView(zone.Mark(
		globals.CreateStartStopZone,
		button.style.Width(commonStyles.ButtonWidth).Render(button.label),
	))
}

type Button struct {
	label string
	style lipgloss.Style
}

type ButtonMap map[int]Button

var buttonMap = ButtonMap{
	globals.TASK_INPUT_EMPTY: Button{
		label: "<- input",
		style: commonStyles.ButtonStyleDisabled,
	},
	globals.TASK_INPUT_VALID_NEW: Button{
		label: "create",
		style: commonStyles.ButtonStyleNotStarted,
	},
	globals.TASK_INPUT_VALID_INACTIVE: Button{
		label: "start",
		style: commonStyles.ButtonStyleNotStarted,
	},
	globals.TASK_INPUT_VALID_ACTIVE: Button{
		label: "stop",
		style: commonStyles.ButtonStyleStarted,
	},
}
