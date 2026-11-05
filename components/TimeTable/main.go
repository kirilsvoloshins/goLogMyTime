package TimeTable

import (
	"fmt"
	"goLogMyTime/common/globals"
	"goLogMyTime/common/styles"
	"goLogMyTime/store"
	"goLogMyTime/types"
	"goLogMyTime/utils"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

type Model struct {
	Store     *store.RootStore
	Zone      types.Zone
	Shortcut  types.SectionShortcut
	Width     int
	IsFocused bool
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case types.FocusZoneMsg:
		if types.Zone(msg) == m.Zone {
			m.IsFocused = true
		}

	case types.BlurZoneMsg:
		if types.Zone(msg) == m.Zone {
			m.IsFocused = false
		}

	case tea.MouseClickMsg:
		if utils.IsLeftClickInBoundsOf(globals.ZoneNextDate, msg) {
			m.Store.Ui.ChangeDaysOffset(true)
		}

		if utils.IsLeftClickInBoundsOf(globals.ZonePrevDate, msg) {
			m.Store.Ui.ChangeDaysOffset(false)
		}

	case tea.KeyPressMsg:
		if !m.IsFocused {
			break
		}

		prev := utils.IsShiftTabKeyMsg(msg)
		next := utils.IsTabKeyMsg(msg)

		if prev || next {
			m.Store.Ui.ChangeDaysOffset(next)
		}
	}

	return m, nil
}

func (m *Model) View() tea.View {
	style := m.Store.Ui.Style
	dateString, timeTable := m.Store.Task.GetDateAndTimeTable(
		m.Store.Ui.TodaysDate,
		m.Store.Ui.TimeTableDaysOffsetFromToday,
	)
	title := lipgloss.Sprintf("Timetable for %s [%s]", dateString, m.Shortcut)

	var sb strings.Builder

	if timeTable != nil {
		paddedLines := styles.GetPaddedTimeTableLines(timeTable, &m.Store.Task.TaskDefinitionMap)
		for id, record := range *timeTable {
			hours, minutes := utils.MinutesToHoursMinutes(record.Minutes)
			fmt.Fprintf(&sb, "\n%s - %dh %02dmin", paddedLines[id], hours, minutes)
		}
	}

	ButtonDayBefore := tea.NewView(zone.Mark(
		globals.ZonePrevDate,
		style.GenericButtonStyle.Render("<< day before"),
	))

	ButtonDayAfter := tea.NewView(zone.Mark(
		globals.ZoneNextDate,
		style.GenericButtonStyle.Render("day after >>"),
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
		string(m.Zone),
		style.GetSectionStyle(m.IsFocused).Width(m.Width).Render(
			lipgloss.JoinVertical(
				lipgloss.Right,
				style.SectionTitle.Render(title),
				content,
			),
		),
	))
}
