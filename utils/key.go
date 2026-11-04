package utils

import tea "charm.land/bubbletea/v2"

func IsEnterKeyMsg(msg tea.KeyPressMsg) bool {
	return msg.String() == "enter"
}

func IsSpaceKeyMsg(msg tea.KeyPressMsg) bool {
	return msg.String() == "space"
}

func IsTabKeyMsg(msg tea.KeyPressMsg) bool {
	return msg.String() == "tab"
}

func IsShiftTabKeyMsg(msg tea.KeyPressMsg) bool {
	return msg.String() == "shift+tab"
}

func IsUpMovementKeyMsg(msg tea.KeyPressMsg) bool {
	str := msg.String()
	return str == "k" || str == "up"
}

func IsRightMovementKeyMsg(msg tea.KeyPressMsg) bool {
	str := msg.String()
	return str == "l"
}

func IsDownMovementKeyMsg(msg tea.KeyPressMsg) bool {
	str := msg.String()
	return str == "j" || str == "down"
}
