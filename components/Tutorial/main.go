package Tutorial

import (
	commonStyles "goLogMyTime/components/common/styles"
	"goLogMyTime/utils/layout"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Model struct {
	width int
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = layout.GetSectionWidth(msg.Width)
	}

	return m, nil
}

func (m Model) View() tea.View {
	var sb strings.Builder
	sb.WriteString("ctrl+c = exit\n")
	sb.WriteString("esc = unfocus input\n")
	sb.WriteString("tab/shift+tab = switch between inputs\n")
	sb.WriteString("1,2,3,4,5 = shortcuts to focus section\n")
	sb.WriteString("j,down,tab = ↓ go down\n")
	sb.WriteString("k,up,shift+tab = ↑ go up\n")
	sb.WriteString("Enter = start task from input\n")
	sb.WriteString("Enter/space = toggle selected > task\n")
	sb.WriteString("tab/l = next timetable date\n")
	sb.WriteString("shift+tab/h = prev timetable date\n")

	return tea.NewView(
		commonStyles.Tutorial.Width(m.width).Render(
			lipgloss.JoinHorizontal(
				lipgloss.Right,
				commonStyles.SectionTitle.Render(sb.String()),
			),
		),
	)
}
