package TasksFavorite

import (
	ButtonStartStop "goLogMyTime/components/ButtonStartStopTask"
	"goLogMyTime/components/common/globals"
	commonStyles "goLogMyTime/components/common/styles"
	"goLogMyTime/domain/taskStore"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

var idPrefix = "F_"

type Model struct {
	DataSource          *taskStore.TaskStore
	IsFocused           bool
	Shortcut            string
	Zone                string
	ButtonsStartStopMap map[string]ButtonStartStop.Model
	Width               int
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	taskDefinitionsFavorite := m.DataSource.GetTaskDefinitionsFavorite()

	for ticket, buttonStartStop := range m.ButtonsStartStopMap {
		var actual = false
		for _, taskDefinition := range taskDefinitionsFavorite {
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

	switch msg := msg.(type) {

	case globals.FocusSectionMsg:
		if string(msg) == m.Shortcut {
			m.focus()
		}

	case globals.BlurSectionMsg:
		if string(msg) == m.Shortcut {
			m.blur()
		}

	}

	return m, tea.Batch(cmds...)
}

func (m *Model) focus() {
	m.IsFocused = true
}

func (m *Model) blur() {
	m.IsFocused = false
}

func (m Model) View() tea.View {
	var title = lipgloss.Sprintf("Favorite tasks [%s]", m.Shortcut)
	var titleRendered = commonStyles.SectionTitle.Render(title)

	var linesRendered []string
	for _, taskDefinition := range m.DataSource.GetTaskDefinitionsFavorite() {
		_, isButtonInMap := m.ButtonsStartStopMap[taskDefinition.Ticket]
		if !isButtonInMap {
			m.ButtonsStartStopMap[taskDefinition.Ticket] = ButtonStartStop.Model{
				DataSource:     m.DataSource,
				TaskDefinition: &taskDefinition,
				IdPrefix:       idPrefix,
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
	if m.DataSource.TasksFavoriteCursor < len(linesRendered) {
		cursor = m.DataSource.TasksFavoriteCursor
	}

	for index, line := range linesRendered {
		sb.WriteString("\n")
		if index == cursor && m.IsFocused {
			sb.WriteString("> ")
		}
		sb.WriteString(line)
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
