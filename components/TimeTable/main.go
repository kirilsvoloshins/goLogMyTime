package TimeTable

import (
	"fmt"
	"goLogMyTime/components/common/globals"
	commonStyles "goLogMyTime/components/common/styles"
	"goLogMyTime/domain/taskStore"
	"goLogMyTime/utils/assert"
	"goLogMyTime/utils/date"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

type Model struct {
	DataSource *taskStore.TaskStore
	IsFocused  bool
	Shortcut   string
	Zone       string
	Width      int
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
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

	return m, nil
}

func (m *Model) focus() {
	m.IsFocused = true
}

func (m *Model) blur() {
	m.IsFocused = false
}

func (m Model) View() tea.View {
	dateString, timeTable := m.DataSource.GetDateAndTimeTable()
	title := lipgloss.Sprintf("Timetable for %s [%s]", dateString, m.Shortcut)

	var sb strings.Builder

	if timeTable != nil {
		for _, record := range *timeTable {
			taskDefinition, ok := m.DataSource.TaskDefinitionMap[record.Name]
			assert.Assert(ok, "very not ok")

			hours, minutes := date.MinutesToHoursMinutes(record.Minutes)
			sb.WriteString(fmt.Sprintf("\n%s - %s - %dh %dmin", taskDefinition.Info, record.Name, hours, minutes))
		}
	}

	ButtonDayBefore := tea.NewView(zone.Mark(
		globals.PrevDateZone,
		commonStyles.ButtonStyleNotStarted.Render("<< day before"),
	))

	ButtonDayAfter := tea.NewView(zone.Mark(
		globals.NextDateZone,
		commonStyles.ButtonStyleNotStarted.Render("day after >>"),
	))

	content := lipgloss.JoinVertical(
		lipgloss.Right,
		lipgloss.JoinHorizontal(
			lipgloss.Bottom,
			ButtonDayBefore.Content,
			"   ",
			ButtonDayAfter.Content,
		),
		sb.String(),
	)

	return tea.NewView(zone.Mark(
		m.Zone,
		commonStyles.GetSectionStyle(m.IsFocused).Width(m.Width).Render(
			lipgloss.JoinVertical(
				lipgloss.Right,
				commonStyles.SectionTitle.Render(title),
				content,
			),
		),
	))
}
