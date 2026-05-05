package NowDoing

import (
	ButtonCreateStartStop "goLogMyTime/components/ButtonCreateStartStopTask"
	"goLogMyTime/components/TextInputWrapper"
	commonStyles "goLogMyTime/components/common/styles"
	"goLogMyTime/domain/taskStore"
	"goLogMyTime/domain/types"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

type Model struct {
	DataSource            *taskStore.TaskStore
	IsFocused             bool
	Shortcut              string
	Zone                  string
	TextInputTask         *TextInputWrapper.Model
	TextInputDescription  *TextInputWrapper.Model
	ButtonCreateStartStop *ButtonCreateStartStop.Model
	Width                 int
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	u1, cmd := m.TextInputTask.Update(msg)
	m.TextInputTask = &u1
	cmds = append(cmds, cmd)

	u2, cmd := m.TextInputDescription.Update(msg)
	m.TextInputDescription = &u2
	cmds = append(cmds, cmd)

	u3, cmd := m.ButtonCreateStartStop.Update(msg)
	m.ButtonCreateStartStop = &u3
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	return tea.NewView(zone.Mark(
		m.Zone,
		commonStyles.GetSectionStyle(m.IsFocused).Width(m.Width).Render(
			lipgloss.JoinVertical(
				lipgloss.Right,
				commonStyles.SectionTitle.Render(m.getTitle()),
				lipgloss.JoinHorizontal(
					lipgloss.Bottom,
					m.TextInputDescription.View().Content,
					m.TextInputTask.View().Content,
					m.ButtonCreateStartStop.View().Content,
				),
			),
		),
	))
}

func (m Model) getTitle() string {
	taskRecordActive := m.DataSource.GetTaskRecordActive()

	if taskRecordActive == nil {
		return lipgloss.Sprintf("Start doing [%s]", m.Shortcut)
	}

	return lipgloss.Sprintf("Now doing [%s]", m.Shortcut)
}

var DESCRIPTION_MAX_LENGTH = 6 // "review"
var DESCRIPTION_MAX_INDEX = DESCRIPTION_MAX_LENGTH - 1

func getStringNowDoing(taskRecordActive *types.TaskRecord) string {
	var sb strings.Builder
	if len(taskRecordActive.Description) > 0 {
		sb.WriteString("\"")
		if len(taskRecordActive.Description) > DESCRIPTION_MAX_LENGTH {
			runes := []rune(taskRecordActive.Description)
			sb.WriteString(string(runes[0:DESCRIPTION_MAX_INDEX]) + "...")
		} else {
			sb.WriteString(taskRecordActive.Description)
		}
		sb.WriteString("\" ")
	}

	return sb.String() + taskRecordActive.Ticket + " "
}
