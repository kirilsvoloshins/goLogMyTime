package TextInputWrapper

import (
	"goLogMyTime/domain/taskStore"
	"goLogMyTime/utils/assert"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

type Model struct {
	DataSource *taskStore.TaskStore
	Zone       string
	TextInput  textinput.Model
	StringPtr  *string
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyPressMsg:
		if !m.TextInput.Focused() {
			break
		}

		switch msg.String() {

		case "esc":
			m.Blur()

		}
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	cmds = append(cmds, cmd)

	if m.TextInput.Focused() {
		assert.Assert(m.StringPtr != nil, "another nil bites the dust")
		*m.StringPtr = strings.TrimSpace(m.TextInput.Value())
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
	return tea.NewView(zone.Mark(
		m.Zone,
		m.TextInput.View(),
	))
}
