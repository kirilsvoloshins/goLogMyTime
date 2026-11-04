package utils

import (
	"goLogMyTime/types"

	tea "charm.land/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

func IsLeftClickInBoundsOf(zoneBeingChecked types.Zone, msg tea.MouseClickMsg) bool {
	return msg.Button == tea.MouseLeft && zone.Get(string(zoneBeingChecked)).InBounds(msg)
}
