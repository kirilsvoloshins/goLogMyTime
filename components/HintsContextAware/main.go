package HintsContextAware

import (
	"goLogMyTime/common/globals"
	"goLogMyTime/store"
	"goLogMyTime/types"
	"strings"

	tea "charm.land/bubbletea/v2"
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
	content := m.getHintsOfActiveSection(m.Store.Ui.FocusedZone)

	return tea.NewView(
		style.HintsContextAware.Width(m.Width).Render(
			content,
		),
	)
}

const (
	hintExit                = "ctrl+c = exit"
	hintUnfocusInput        = "esc = focus now doing"
	hintChangeColorTheme    = "t/T = theme"
	hintFocusShortcuts      = "1,2,3,4,5 = focus [section]"
	hintGoUp                = "k,up = ↑ up"
	hintGoDown              = "j,down = ↓ down"
	hintSwitchBetweenInputs = "tab/shift+tab = ← → jump between inputs"
	hintStartTaskFromInput  = "Enter = start/stop from valid input"
	hintToggleSelected      = "Enter/space = toggle task"
	hintFavoriteSelected    = "f = toggle favorite"
	hintNextDate            = "tab = → next date"
	hintPrevDate            = "shift+tab = ← prev date"
	hintFocusTaskInput      = "/ = focus task"
	hintFocusNotes          = "i = focus notes"
)

func (m *Model) getHintsOfActiveSection(zoneActive types.Zone) string {
	hints := []string{}

	switch zoneActive {
	case globals.ZoneUnfocused:
		hints = append(hints, hintFocusTaskInput, hintFocusNotes, hintFocusShortcuts, hintChangeColorTheme)
	case globals.InteractiveSections.NowDoing.Zone:
		hints = append(hints, hintFocusTaskInput, hintFocusNotes, hintGoDown, hintFocusShortcuts, hintChangeColorTheme)
	case globals.ZoneTextInputTitle, globals.ZoneTextInputTask:
		hints = append(hints, hintUnfocusInput, hintSwitchBetweenInputs, hintStartTaskFromInput)
	case globals.InteractiveSections.TasksRecent.Zone, globals.InteractiveSections.TasksFavorite.Zone:
		hints = append(hints, hintFocusTaskInput, hintFocusNotes, hintToggleSelected, hintFavoriteSelected, hintGoUp, hintGoDown, hintFocusShortcuts, hintChangeColorTheme)
	case globals.InteractiveSections.TimeTable.Zone:
		hints = append(hints, hintFocusTaskInput, hintFocusNotes, hintPrevDate, hintNextDate, hintGoUp, hintFocusShortcuts, hintChangeColorTheme)
	case globals.InteractiveSections.EmbeddedFile.Zone:
		hints = append(hints, hintUnfocusInput)
	}

	hints = append(hints, hintExit)

	return " " + strings.Join(hints, " • ")
}
