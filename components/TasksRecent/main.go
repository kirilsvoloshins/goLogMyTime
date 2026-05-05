package TasksRecent

import (
	ButtonStartStop "goLogMyTime/components/ButtonStartStopTask"
	commonStyles "goLogMyTime/components/common/styles"
	"goLogMyTime/domain/taskStore"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

type Model struct {
	DataSource          *taskStore.TaskStore
	IsFocused           bool
	Shortcut            string
	Zone                string
	ButtonsStartStopMap map[string]ButtonStartStop.Model
	Width               int
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	taskDefinitionsRecent := m.DataSource.GetTaskDefinitionsRecent(true)

	for ticket, buttonStartStop := range m.ButtonsStartStopMap {
		var actual = false
		for _, taskDefinition := range taskDefinitionsRecent {
			if ticket == taskDefinition.Ticket {
				actual = true
				break
			}
		}
		if !actual {
			delete(m.ButtonsStartStopMap, ticket)
		}

		cmd = buttonStartStop.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	var title = lipgloss.Sprintf("Recent tasks [%s]", m.Shortcut)
	var titleRendered = commonStyles.SectionTitle.Render(title)

	var linesRendered []string
	for _, taskDefinition := range m.DataSource.GetTaskDefinitionsRecent(true) {
		_, isButtonInMap := m.ButtonsStartStopMap[taskDefinition.Ticket]
		if !isButtonInMap {
			m.ButtonsStartStopMap[taskDefinition.Ticket] = ButtonStartStop.Model{
				DataSource:     m.DataSource,
				TaskDefinition: &taskDefinition,
				IdPrefix:       "R_",
			}
		}

		line := lipgloss.JoinHorizontal(
			lipgloss.Bottom,
			lipgloss.Sprintf("%s - %s ", taskDefinition.Info, taskDefinition.Ticket),
			m.ButtonsStartStopMap[taskDefinition.Ticket].View().Content,
		)
		linesRendered = append(linesRendered, line)
	}

	var sb strings.Builder
	sb.WriteString(titleRendered)

	cursor := 0
	if m.DataSource.TasksRecentCursor < len(linesRendered) {
		cursor = m.DataSource.TasksRecentCursor
	}

	for index, line := range linesRendered {
		sb.WriteString("\n")
		if index == cursor && m.IsFocused {
			sb.WriteString("> " + line)
		} else {
			sb.WriteString(line)
		}
	}

	return tea.NewView(zone.Mark(
		m.Zone,
		commonStyles.GetSectionStyle(m.IsFocused).Width(m.Width).Render(
			lipgloss.JoinVertical(
				lipgloss.Right,
				sb.String(),
			),
		),
	))
}
