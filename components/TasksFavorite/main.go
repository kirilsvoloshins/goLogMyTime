package TasksFavorite

import (
	ButtonStartStop "goLogMyTime/components/ButtonStartStopTask"
	"goLogMyTime/store"
	"goLogMyTime/types"
	"goLogMyTime/utils"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

var idPrefix = "F_"

type Model struct {
	Store               *store.RootStore
	Shortcut            types.SectionShortcut
	Zone                types.Zone
	Width               int
	IsFocused           bool
	ButtonsStartStopMap map[string]ButtonStartStop.Model
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	taskDefinitionsFavorite := m.Store.Task.TaskDefinitionsFavorite

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

	case types.FocusZoneMsg:
		if types.Zone(msg) == m.Zone {
			m.IsFocused = true
		}

	case types.BlurZoneMsg:
		if types.Zone(msg) == m.Zone {
			m.IsFocused = false
		}

	}

	keyPressMsg, isKeyPress := msg.(tea.KeyPressMsg)
	if isKeyPress && m.IsFocused {
		isToggleKey := utils.IsEnterKeyMsg(keyPressMsg) || utils.IsSpaceKeyMsg(keyPressMsg)

		if isToggleKey && len(m.Store.Task.TaskDefinitionsFavorite) > 0 {
			taskDef := m.Store.Task.TaskDefinitionsFavorite[m.Store.Ui.TasksFavoriteCursor]
			m.Store.Task.TaskToggle(taskDef.Ticket)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() tea.View {
	var style = m.Store.Ui.Style
	var title = lipgloss.Sprintf("Favorite tasks [%s]", m.Shortcut)
	var titleRendered = style.SectionTitle.Render(title)

	var linesRendered []string
	for _, taskDefinition := range m.Store.Task.TaskDefinitionsFavorite {
		_, isButtonInMap := m.ButtonsStartStopMap[taskDefinition.Ticket]
		if !isButtonInMap {
			m.ButtonsStartStopMap[taskDefinition.Ticket] = ButtonStartStop.Model{
				Store:    m.Store,
				Ticket:   taskDefinition.Ticket,
				IdPrefix: idPrefix,
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
	if m.Store.Ui.TasksFavoriteCursor < len(linesRendered) {
		cursor = m.Store.Ui.TasksFavoriteCursor
	}

	for index, line := range linesRendered {
		sb.WriteString("\n")
		if index == cursor && m.IsFocused {
			sb.WriteString("> ")
		}
		sb.WriteString(line)
	}

	return tea.NewView(zone.Mark(
		string(m.Zone),
		style.GetSectionStyle(m.IsFocused).Width(m.Width).Render(
			lipgloss.JoinVertical(
				lipgloss.Right,
				sb.String(),
			),
		),
	))
}
