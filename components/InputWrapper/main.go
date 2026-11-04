package InputWrapper

import (
	"goLogMyTime/common/globals"
	"goLogMyTime/store"
	"goLogMyTime/types"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

type Model struct {
	Store         *store.RootStore
	Zone          types.Zone
	TextInput     textinput.Model
	StringPtr     *string
	hasActiveTask bool
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case types.BlurZoneMsg:
		if types.Zone(msg) == m.Zone {
			m.Blur()
		}
	case tea.KeyPressMsg:
		keyString := msg.String()

		switch keyString {
		case "esc":
			if !m.TextInput.Focused() {
				break
			}

			cmds = append(cmds, globals.FocusZone(globals.InteractiveSections.NowDoing.Zone))
		}
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	cmds = append(cmds, cmd)

	if m.Store.Task.TaskRecordActive != nil {
		m.hasActiveTask = true
		m.Blur()

		switch m.Zone {
		case globals.ZoneTextInputTask:
			*m.StringPtr = m.Store.Task.TaskRecordActive.Ticket
		case globals.ZoneTextInputTitle:
			*m.StringPtr = m.Store.Task.TaskDefinitionActive.Info
		}
	} else {
		if m.hasActiveTask {
			*m.StringPtr = ""
		} else if m.TextInput.Focused() {
			*m.StringPtr = strings.TrimSpace(m.TextInput.Value())
		}

		m.hasActiveTask = false
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) Focus() tea.Cmd {
	return m.TextInput.Focus()
}

func (m *Model) Blur() {
	m.TextInput.Blur()
}

func (m Model) View() tea.View {
	styles := m.TextInput.Styles()
	styles.Blurred.Placeholder = m.Store.Ui.Style.TextLowContrast
	styles.Focused.Placeholder = m.Store.Ui.Style.TextLowContrast
	styles.Focused.Text = m.Store.Ui.Style.TextHighContrast
	styles.Blurred.Text = m.Store.Ui.Style.TextHighContrast
	styles.Focused.Prompt = m.Store.Ui.Style.TextHighContrast
	styles.Blurred.Prompt = m.Store.Ui.Style.TextHighContrast
	m.TextInput.SetStyles(styles)

	return tea.NewView(zone.Mark(
		string(m.Zone),
		m.TextInput.View(),
	))
}
