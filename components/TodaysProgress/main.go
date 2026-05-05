package TodaysProgress

import (
	"fmt"
	commonStyles "goLogMyTime/components/common/styles"
	"goLogMyTime/domain/taskStore"
	"goLogMyTime/utils/layout"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Model struct {
	DataSource *taskStore.TaskStore
	width      int
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = layout.GetSectionWidth(msg.Width)
	}

	return m, nil
}

func (m Model) View() tea.View {
	todaysProgressPercent := m.DataSource.GetTodaysProgreesPercent()
	title := fmt.Sprintf("Days progress: %0.f%%", todaysProgressPercent)
	content := m.getStringBar(todaysProgressPercent)

	return tea.NewView(
		commonStyles.BorderedCentered.Width(m.width).Render(
			lipgloss.JoinVertical(
				lipgloss.Right,
				commonStyles.SectionTitle.Render(title),
				content,
			),
		),
	)
}

func (m Model) getStringBar(daysProgressPercent float64) string {
	var charNumTotal int
	var charNumFilled int
	charNumTotal = min(m.width-12, 33)
	ratio := daysProgressPercent / 100.0

	charNumFilled = min(charNumTotal, int(ratio*float64(charNumTotal)))

	var sb strings.Builder
	sb.WriteString("0 [")
	for charId := range charNumTotal {
		if charId < charNumFilled {
			sb.WriteString("#")
		} else {
			sb.WriteString(" ")
		}
	}
	sb.WriteString("] 8h")

	return sb.String()
}
