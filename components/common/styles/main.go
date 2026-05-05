package commonStyles

import (
	"charm.land/lipgloss/v2"
)

var buttonPadding = 1
var ButtonWidth = 10
var PanelWidth = 45

var Tutorial = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FF00FF")).
	Padding(0, 0).
	Align(lipgloss.Left, lipgloss.Top)

var EmbeddedSectionContainerBlurred = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("69")).
	Border(lipgloss.RoundedBorder())

var EmbeddedSectionContainerFocused = EmbeddedSectionContainerBlurred.
	BorderForeground(lipgloss.Color("205"))

var SectionContainerBlurred = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FF00FF")).
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("69")).
	Border(lipgloss.RoundedBorder()).
	Padding(0, 1).
	Align(lipgloss.Right, lipgloss.Top)

var SectionContainerFocused = SectionContainerBlurred.
	BorderForeground(lipgloss.Color("205"))

var SectionTitle = lipgloss.NewStyle().
	Bold(true).
	PaddingBottom(1).
	Foreground(lipgloss.Color("#FFFFFF"))

var BorderedCentered = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FF00FF")).
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("69")).
	Border(lipgloss.RoundedBorder()).
	Padding(0, 1).
	Align(lipgloss.Right, lipgloss.Right)

var ButtonStyleDisabled = lipgloss.NewStyle().
	Foreground(lipgloss.Color("255")).
	Background(lipgloss.Color("62")).
	Padding(0, buttonPadding)

var ButtonStyleNotStarted = lipgloss.NewStyle().
	Foreground(lipgloss.Color("255")).
	Background(lipgloss.Color("62")).
	Padding(0, buttonPadding)

var ButtonStyleStarted = lipgloss.NewStyle().
	Foreground(lipgloss.Color("255")).
	Background(lipgloss.Color("62")).
	Padding(0, buttonPadding).
	Background(lipgloss.Color("205")).
	Bold(true)

func GetSectionStyle(isFocused bool) lipgloss.Style {
	if isFocused {
		return SectionContainerFocused
	}

	return SectionContainerBlurred
}
