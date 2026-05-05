package TimeTable

import (
	"fmt"
	commonStyles "goLogMyTime/components/common/styles"
	"goLogMyTime/domain/taskStore"
	"goLogMyTime/utils/assert"
	"goLogMyTime/utils/date"
	"goLogMyTime/utils/layout"
	"goLogMyTime/utils/mouse"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

var buttonDayBeforeId = "dayBefore"
var buttonDayAfterId = "dayAfter"

type Model struct {
	DataSource *taskStore.TaskStore
	IsFocused  bool
	Shortcut   string
	Zone       string
	width      int
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseClickMsg:
		m.IsFocused = mouse.IsLeftClickInBoundsOf(m.Zone, msg)
		if mouse.IsLeftClickInBoundsOf(buttonDayBeforeId, msg) {
			m.DataSource.TimeTableDaysOffsetFromToday--
		} else if mouse.IsLeftClickInBoundsOf(buttonDayAfterId, msg) {
			m.DataSource.TimeTableDaysOffsetFromToday++
		}
	case tea.WindowSizeMsg:
		m.width = layout.GetSectionWidth(msg.Width)
	}

	return m, nil
}

func (m Model) View() tea.View {
	dateBeingChecked := m.DataSource.TodaysDate.AddDate(0, 0, m.DataSource.TimeTableDaysOffsetFromToday)
	dateString := date.DateToString(dateBeingChecked)
	var title = lipgloss.Sprintf("Timetable for %s [%s]", dateString, m.Shortcut)

	timetable := m.DataSource.GetTimeTable(dateString)

	var sb strings.Builder

	if timetable != nil {
		for _, record := range *timetable {
			sb.WriteString("\n")
			taskDefinition, ok := m.DataSource.TaskDefinitionMap[record.Name]
			assert.Assert(ok, "very not ok")

			hours := int(record.Minutes / 60)
			minutes := record.Minutes - hours*60
			sb.WriteString(fmt.Sprintf("%s - %s - %dh %dmin", taskDefinition.Info, record.Name, hours, minutes))
		}
	}

	ButtonDayBefore := tea.NewView(zone.Mark(
		buttonDayBeforeId,
		commonStyles.ButtonStyleNotStarted.Render("<< day before"),
	))

	ButtonDayAfter := tea.NewView(zone.Mark(
		buttonDayAfterId,
		commonStyles.ButtonStyleNotStarted.Render("day after >>"),
	))

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		ButtonDayBefore.Content,
		"   ",
		ButtonDayAfter.Content,
	)

	content := lipgloss.JoinVertical(
		lipgloss.Right,
		buttons,
		sb.String(),
	)

	return tea.NewView(zone.Mark(
		m.Zone,
		commonStyles.GetSectionStyle(m.IsFocused).Width(m.width).Render(
			lipgloss.JoinVertical(
				lipgloss.Right,
				commonStyles.SectionTitle.Render(title),
				content,
			),
		),
	))
}
