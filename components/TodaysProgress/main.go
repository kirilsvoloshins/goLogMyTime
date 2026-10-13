package TodaysProgress

import (
	"fmt"
	"goLogMyTime/domain/taskStore"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Model struct {
	DataSource *taskStore.TaskStore
	Width      int
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() tea.View {
	style := m.DataSource.Style
	todaysProgressPercent := m.DataSource.GetTodaysProgreesPercent()
	title := fmt.Sprintf("Days progress: %0.f%%", todaysProgressPercent)
	content := "0 [" + m.getStringBar(todaysProgressPercent) + "] 8h"

	return tea.NewView(
		style.GetSectionStyle(false).Width(m.Width).Render(
			lipgloss.JoinVertical(
				lipgloss.Right,
				style.SectionTitle.Render(title),
				content,
			),
		),
	)
}

func (m Model) getStringBar(daysProgressPercent float64) string {
	charNumTotal := min(m.Width-12, 33)
	ratio := daysProgressPercent / 100.0
	charNumFilled := min(charNumTotal, int(ratio*float64(charNumTotal)))

	var sb strings.Builder

	for charId := range charNumTotal {
		if charId < charNumFilled {
			sb.WriteString(filled)
		} else {
			sb.WriteString(empty)
		}
	}

	return sb.String()
}

var filled = "#"
var empty = " "
