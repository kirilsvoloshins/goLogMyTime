package NowDoing

import (
	"goLogMyTime/common/globals"
	ButtonCreateStartStop "goLogMyTime/components/ButtonCreateStartStopTask"
	"goLogMyTime/components/InputWrapper"
	"goLogMyTime/store"
	"goLogMyTime/types"
	"goLogMyTime/utils"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

type Model struct {
	Store                 *store.RootStore
	Width                 int
	IsFocused             bool
	Shortcut              types.SectionShortcut
	Zone                  types.Zone
	TextInputTask         *InputWrapper.Model
	TextInputTitle        *InputWrapper.Model
	ButtonCreateStartStop *ButtonCreateStartStop.Model
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case types.BlurZoneMsg:
		zoneToBlur := types.Zone(msg)

		if zoneToBlur == m.Zone {
			m.IsFocused = false
		}

	case types.FocusZoneMsg:
		zoneToFocus := types.Zone(msg)

		switch zoneToFocus {
		case m.Zone:
			m.IsFocused = true
		case globals.ZoneTextInputTask:
			m.IsFocused = true
			m.TextInputTitle.TextInput.Blur()
			cmds = append(cmds, m.TextInputTask.Focus())
		case globals.ZoneTextInputTitle:
			m.IsFocused = true
			m.TextInputTask.TextInput.Blur()
			cmds = append(cmds, m.TextInputTitle.Focus())
		}

	case tea.KeyPressMsg:
		if !m.IsFocused {
			break
		}

		if (utils.IsTabKeyMsg(msg) || utils.IsShiftTabKeyMsg(msg)) && !m.hasActiveTask() {
			if m.TextInputTitle.TextInput.Focused() {
				cmds = append(cmds, globals.FocusZone(globals.ZoneTextInputTask))
			} else if m.TextInputTask.TextInput.Focused() {
				cmds = append(cmds, globals.FocusZone(globals.ZoneTextInputTitle))
			}
			break
		}

		if utils.IsEnterKeyMsg(msg) {
			taskStateChange := m.Store.Task.TaskToggleIfValid(m.Store.Ui.InputTicket, m.Store.Ui.InputTitle)
			if taskStateChange == globals.TASK_GOT_STOPPED {
				cmds = append(cmds, globals.FocusZone(globals.ZoneTextInputTask))
			}
		}
	}

	m.TextInputTask, cmd = m.TextInputTask.Update(msg)
	m.TextInputTask.TextInput.SetValue(m.Store.Ui.InputTicket)
	cmds = append(cmds, cmd)

	m.TextInputTitle, cmd = m.TextInputTitle.Update(msg)
	m.TextInputTitle.TextInput.SetValue(m.Store.Ui.InputTitle)
	cmds = append(cmds, cmd)

	m.ButtonCreateStartStop, cmd = m.ButtonCreateStartStop.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() tea.View {
	style := m.Store.Ui.Style

	var title string
	var content string

	if m.hasActiveTask() {
		title = lipgloss.Sprintf("Now doing [%s]", m.Shortcut)
	} else {
		title = lipgloss.Sprintf("Start doing [%s]", m.Shortcut)
	}

	content = lipgloss.JoinVertical(
		lipgloss.Right,
		m.TextInputTitle.View().Content,
		"",
		m.TextInputTask.View().Content,
		"",
		m.ButtonCreateStartStop.View().Content,
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

func (m *Model) hasActiveTask() bool {
	return m.Store.Task.TaskRecordActive != nil
}
