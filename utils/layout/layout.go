package layout

import commonStyles "goLogMyTime/components/common/styles"

func GetSectionWidth(windowWidth int) int {
	return commonStyles.PanelWidth
}

func GetEditorWidth(windowWidth int) int {
	return windowWidth - GetSectionWidth(windowWidth)
}
