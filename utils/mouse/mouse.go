package mouse

import (
	tea "charm.land/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

func IsLeftClickInBoundsOf(buttonId string, msg tea.MouseClickMsg) bool {
	return msg.Button == tea.MouseLeft && zone.Get(buttonId).InBounds(msg)
}
