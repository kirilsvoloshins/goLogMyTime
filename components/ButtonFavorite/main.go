package ButtonFavorite

import (
	"goLogMyTime/store"
	"goLogMyTime/types"
	"goLogMyTime/utils"

	tea "charm.land/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

type Model struct {
	Store    *store.RootStore
	Ticket   string
	IdPrefix string
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.MouseClickMsg:
		if utils.IsLeftClickInBoundsOf(m.getZone(), msg) {
			m.Store.Task.TaskFavoriteToggle(m.Ticket)
		}
	}

	return nil
}

func (m *Model) View() tea.View {
	isFavorite := m.Store.Task.IsTaskFavorite(m.Ticket)
	buttonLabel := m.getLabel(isFavorite)
	buttonStyle := m.Store.Ui.Style.GetFavoriteButtonStyle(isFavorite)

	return tea.NewView(zone.Mark(
		string(m.getZone()),
		buttonStyle.Render(buttonLabel),
	))
}

func (m *Model) getZone() types.Zone {
	return types.Zone(m.IdPrefix + m.Ticket)
}

func (m *Model) getLabel(isFavorite bool) string {
	if isFavorite {
		return "💔"
	}

	return "❤️ "
}
