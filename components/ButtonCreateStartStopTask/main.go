package ButtonCreateStartStop

import (
	"goLogMyTime/components/common/globals"
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
	buttonStyle := m.DataSource.Style.GetCreateStartStopButtonStyle(m.DataSource.TaskInputState)
	buttonLabel, ok := buttonMap[m.DataSource.TaskInputState]
	assert.Assert(ok, "Invalid task input state")

	return tea.NewView(zone.Mark(
		globals.CreateStartStopZone,
		buttonStyle.Render(buttonLabel),
	))
}

type Button struct {
	label string
	style lipgloss.Style
}

type ButtonLabelMap map[int]string

var buttonMap = ButtonLabelMap{
	globals.TASK_INPUT_EMPTY:          "< input",
	globals.TASK_INPUT_VALID_NEW:      "create",
	globals.TASK_INPUT_VALID_INACTIVE: "start",
	globals.TASK_INPUT_VALID_ACTIVE:   "stop",
}
