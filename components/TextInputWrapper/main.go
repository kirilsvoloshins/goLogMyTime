package TextInputWrapper

import (
	"goLogMyTime/utils/assert"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

type Model struct {
	Zone      string
	TextInput textinput.Model
	StringPtr *string
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			m.TextInput.Blur()
			return m, nil
		}
	default:
		// TODO: implement enabled/disabled states - those must ignore the input events etc...
		if m.TextInput.Focused() {
			assert.Assert(m.StringPtr != nil, "another nil bites the dust")
			*m.StringPtr = strings.TrimSpace(m.TextInput.Value())
		}
	}

	m.TextInput, cmd = m.TextInput.Update(msg)

	return m, cmd
}

func (m Model) View() tea.View {
	return tea.NewView(
		zone.Mark(m.Zone, m.TextInput.View()),
	)
}
