package TodaysProgress

import (
	"fmt"
	"goLogMyTime/store"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Model struct {
	Store *store.RootStore
	Width int
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	return m, nil
}

func (m *Model) View() tea.View {
	style := m.Store.Ui.Style
	workTimeMin := m.Store.Task.GetTodaysWorkTimeMin(m.Store.Ui.TodaysDate)
	todaysProgressPercent := float64((100 * workTimeMin) / WORK_DAY_IN_MINUTES)
	title := fmt.Sprintf("Days progress: %0.f%%", todaysProgressPercent)
	content := ""

	if todaysProgressPercent >= 100 {
		content = "time to go home :)"
	} else {
		content = "0 [" + m.getStringBar(todaysProgressPercent) + "] 8h"
	}

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

const WORK_DAY_IN_MINUTES = 8 * 60

func (m *Model) getStringBar(daysProgressPercent float64) string {
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

const filled = "#"
const empty = " "
