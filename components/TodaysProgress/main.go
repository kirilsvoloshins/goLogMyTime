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
	progress := m.Store.Task.GetTodaysProgress(m.Store.Ui.TodaysDate)
	content := ""

	if progress >= 1 {
		content = "time to go home :)"
	} else {
		content = m.getProgressBar(progress)
	}

	s := style.GetSectionStyle(false).Width(m.Width)

	return tea.NewView(
		s.Render(
			lipgloss.JoinVertical(
				lipgloss.Right,
				style.SectionTitle.Render(
					fmt.Sprintf("Days progress: %0.f%%", progress*100),
				),
				content,
			),
		),
	)
}

func (m *Model) getProgressBar(progress float32) string {
	charNumTotal := min(m.Width-12, 33)
	numFilled := max(0, min(charNumTotal, int(progress*float32(charNumTotal))))
	numEmpty := max(0, charNumTotal-numFilled)

	var sb strings.Builder
	sb.WriteString("0 [")
	sb.WriteString(strings.Repeat("#", numFilled))
	sb.WriteString(strings.Repeat(" ", numEmpty))
	sb.WriteString("] 8h")
	return sb.String()
}
