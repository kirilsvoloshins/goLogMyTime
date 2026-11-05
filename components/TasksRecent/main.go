package TasksRecent

import (
	"goLogMyTime/common/styles"
	ButtonStartStop "goLogMyTime/components/ButtonStartStopTask"
	"goLogMyTime/store"
	"goLogMyTime/types"
	"goLogMyTime/utils"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

type Model struct {
	Store               *store.RootStore
	Shortcut            types.SectionShortcut
	Zone                types.Zone
	Width               int
	IsFocused           bool
	ButtonsStartStopMap map[string]ButtonStartStop.Model
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd

	taskDefs := m.Store.Task.TaskDefinitionsRecent
	if len(taskDefs) > 0 && m.Store.Ui.TasksRecentCursor >= len(taskDefs) {
		m.Store.Ui.TasksRecentCursor = len(taskDefs) - 1
	}

	for ticket, buttonStartStop := range m.ButtonsStartStopMap {
		var actual = false
		for _, taskDefinition := range taskDefs {
			if ticket == taskDefinition.Ticket {
				actual = true
				break
			}
		}
		if !actual {
			delete(m.ButtonsStartStopMap, ticket)
		}

		cmd = buttonStartStop.Update(msg)
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

		if isToggleKey && len(m.Store.Task.TaskDefinitionsRecent) > 0 {
			taskDef := m.Store.Task.TaskDefinitionsRecent[m.Store.Ui.TasksRecentCursor]
			m.Store.Task.TaskToggle(taskDef.Ticket)
		}
	}

	return m, cmd
}

func (m *Model) View() tea.View {
	var style = m.Store.Ui.Style
	var title = lipgloss.Sprintf("Recent tasks [%s]", m.Shortcut)
	var titleRendered = style.SectionTitle.Render(title)

	taskDefs := m.Store.Task.TaskDefinitionsRecent

	for _, taskDef := range taskDefs {
		_, isButtonInMap := m.ButtonsStartStopMap[taskDef.Ticket]
		if !isButtonInMap {
			m.ButtonsStartStopMap[taskDef.Ticket] = ButtonStartStop.Model{
				Store:    m.Store,
				Ticket:   taskDef.Ticket,
				IdPrefix: "R_",
			}
		}
	}

	var linesRendered []string
	paddedDefLines := styles.GetPaddedDefLines(taskDefs)

	for id, taskDef := range taskDefs {
		line := lipgloss.JoinHorizontal(
			lipgloss.Bottom,
			paddedDefLines[id],
			m.ButtonsStartStopMap[taskDef.Ticket].View().Content,
		)
		linesRendered = append(linesRendered, line)
	}

	var sb strings.Builder
	sb.WriteString(titleRendered)

	for index, line := range linesRendered {
		sb.WriteString("\n")
		if index == m.Store.Ui.TasksRecentCursor && m.IsFocused {
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
