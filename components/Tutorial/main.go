package Tutorial

import (
	commonStyles "goLogMyTime/components/common/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Model struct {
	Width int
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() tea.View {
	return tea.NewView(
		commonStyles.Tutorial.Width(m.Width).Render(
			lipgloss.JoinHorizontal(
				lipgloss.Right,
				commonStyles.SectionTitle.Render(tutorialString),
			),
		),
	)
}

var tutorialString = `ctrl+c = exit
esc = unfocus input
tab/shift+tab = switch between inputs
1,2,3,4,5 = shortcuts to focus section
j,down,tab = ↓ go down
k,up,shift+tab = ↑ go up
Enter = start task from input
Enter/space = toggle selected > task
tab/l = next timetable date
shift+tab/h = prev timetable date`
