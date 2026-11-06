package TasksFavorite

import (
	"goLogMyTime/common/styles"
	"goLogMyTime/components/ButtonFavorite"
	ButtonStartStop "goLogMyTime/components/ButtonStartStopTask"
	"goLogMyTime/store"
	"goLogMyTime/types"
	"goLogMyTime/utils"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

const (
	prefixBtnStartStop = "FS_"
	prefixBtnFavorite  = "FF_"
)

// TODO: refactor, so whole line is a component - no per feature button maps :)
type Model struct {
	Store               *store.RootStore
	Shortcut            types.SectionShortcut
	Zone                types.Zone
	Width               int
	IsFocused           bool
	ButtonsStartStopMap map[types.Ticket]ButtonStartStop.Model
	ButtonsFavoriteMap  map[types.Ticket]ButtonFavorite.Model
}

func (m *Model) Init() {
	m.syncButtons()
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.syncButtons()

	for ticket := range m.ButtonsStartStopMap {
		if btn, ok := m.ButtonsStartStopMap[ticket]; ok {
			cmd = btn.Update(msg)
			cmds = append(cmds, cmd)
		}

		if btn, ok := m.ButtonsFavoriteMap[ticket]; ok {
			cmd = btn.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	taskDefs := m.Store.Task.TaskDefinitionsFavorite

	if len(taskDefs) > 0 && m.Store.Ui.TasksFavoriteCursor >= len(taskDefs) {
		m.Store.Ui.TasksFavoriteCursor = len(taskDefs) - 1
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

	case tea.KeyPressMsg:
		if m.IsFocused {
			isToggleKey := utils.IsEnterKeyMsg(msg) || utils.IsSpaceKeyMsg(msg)

			if isToggleKey && len(m.Store.Task.TaskDefinitionsFavorite) > 0 {
				taskDef := m.Store.Task.TaskDefinitionsFavorite[m.Store.Ui.TasksFavoriteCursor]
				m.Store.Task.TaskToggle(taskDef.Ticket)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() tea.View {
	var style = m.Store.Ui.Style
	var title = lipgloss.Sprintf("Favorite tasks [%s]", m.Shortcut)
	var titleRendered = style.SectionTitle.Render(title)

	taskDefs := m.Store.Task.TaskDefinitionsFavorite

	var linesRendered []string
	paddedDefLines := styles.GetPaddedDefLines(taskDefs)

	for id, taskDef := range taskDefs {
		btnStartStop := m.ButtonsStartStopMap[taskDef.Ticket]
		btnFavorite := m.ButtonsFavoriteMap[taskDef.Ticket]
		line := lipgloss.JoinHorizontal(
			lipgloss.Bottom,
			paddedDefLines[id],
			btnStartStop.View().Content,
			btnFavorite.View().Content,
		)
		linesRendered = append(linesRendered, line)
	}

	var sb strings.Builder
	sb.WriteString(titleRendered)

	for index, line := range linesRendered {
		sb.WriteString("\n")
		if index == m.Store.Ui.TasksFavoriteCursor && m.IsFocused {
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

func (m *Model) syncButtons() {
	taskDefs := m.Store.Task.TaskDefinitionsRecent

	validTickets := make(map[types.Ticket]bool, len(taskDefs))
	for _, taskDef := range taskDefs {
		validTickets[taskDef.Ticket] = true

		if _, exists := m.ButtonsStartStopMap[taskDef.Ticket]; !exists {
			m.ButtonsStartStopMap[taskDef.Ticket] = ButtonStartStop.Model{
				Store:    m.Store,
				Ticket:   taskDef.Ticket,
				IdPrefix: prefixBtnStartStop,
			}
			m.ButtonsFavoriteMap[taskDef.Ticket] = ButtonFavorite.Model{
				Store:    m.Store,
				Ticket:   taskDef.Ticket,
				IdPrefix: prefixBtnFavorite,
			}
		}
	}

	for ticket := range m.ButtonsStartStopMap {
		if !validTickets[ticket] {
			delete(m.ButtonsStartStopMap, ticket)
			delete(m.ButtonsFavoriteMap, ticket)
		}
	}
}
