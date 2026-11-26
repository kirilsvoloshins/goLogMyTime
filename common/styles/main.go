package styles

import (
	"goLogMyTime/common/colors"
	"goLogMyTime/types"
	"goLogMyTime/utils"
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
)

type SizesType struct {
	ButtonPadding int
	ButtonWidth   int
	SectionWidth  int
}

var Sizes = SizesType{
	ButtonPadding: 1,
	ButtonWidth:   10,
	SectionWidth:  46,
}

type Style struct {
	BgColor                 color.Color
	TextLowContrast         lipgloss.Style
	TextHighContrast        lipgloss.Style
	HintsContextAware       lipgloss.Style
	SectionTitle            lipgloss.Style
	GetSectionStyle         func(focused bool) lipgloss.Style
	GetEmbeddedSectionStyle func(focused bool) lipgloss.Style
	ButtonDisabled          lipgloss.Style
	ButtonInactive          lipgloss.Style
	ButtonActive            lipgloss.Style
	GetStartStopButtonStyle func(active bool) lipgloss.Style
	GetFavoriteButtonStyle  func(isFavorite bool) lipgloss.Style
	GenericButtonStyle      lipgloss.Style
	EmbeddedStatusBar       lipgloss.Style
}

func GetStyle(paletteIndex int) Style {
	palette := colors.Palettes[paletteIndex]
	colors := colors.ColorsHexToColors(palette.ColorsHex)

	sectionTitle := lipgloss.NewStyle().
		Foreground(colors.TitleColor).
		Bold(true).
		PaddingBottom(1)

	buttonDisabled := lipgloss.NewStyle().
		Foreground(colors.DisabledButtonTextColor).
		Background(colors.DisabledButtonColor).
		Padding(0, Sizes.ButtonPadding)

	buttonInactive := lipgloss.NewStyle().
		Foreground(colors.DisabledButtonTextColor).
		Background(colors.DisabledButtonColor).
		Padding(0, Sizes.ButtonPadding)

	buttonActive := lipgloss.NewStyle().
		Foreground(colors.DisabledButtonTextColor).
		Padding(0, Sizes.ButtonPadding).
		Background(colors.FocusedBorderColor).
		Bold(true).
		Width(7)

	buttonFavour := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89")).
		Padding(0, Sizes.ButtonPadding).
		Bold(true)

	buttonUnfavour := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5F87")).
		Padding(0, Sizes.ButtonPadding).
		Bold(true)

	colorHighContrast := lipgloss.Color("#FFF")
	colorLowContrast := lipgloss.Color("#7C7C7C")

	if palette.Mode == "light" {
		colorHighContrast = lipgloss.Color("#000")
		colorLowContrast = lipgloss.Color("#6B6B6B")
	}

	textHighContrast := lipgloss.NewStyle().Foreground(colorHighContrast)
	textLowContrast := lipgloss.NewStyle().Foreground(colorLowContrast)

	return Style{
		BgColor:           colors.BgColor,
		TextLowContrast:   textLowContrast,
		TextHighContrast:  textHighContrast,
		HintsContextAware: textHighContrast,
		SectionTitle:      sectionTitle,
		GetSectionStyle: func(focused bool) lipgloss.Style {
			sectionContainerBlurred := lipgloss.NewStyle().
				Foreground(colors.BaseTextColor).
				BorderForeground(colors.BlurredBorderColor).
				Border(lipgloss.RoundedBorder()).
				Padding(0, 1).
				Align(lipgloss.Right, lipgloss.Top)

			if focused {
				return sectionContainerBlurred.
					BorderForeground(colors.FocusedBorderColor)
			}

			return sectionContainerBlurred
		},
		ButtonDisabled: buttonDisabled,
		ButtonInactive: buttonInactive,
		ButtonActive:   buttonActive,
		GetStartStopButtonStyle: func(active bool) lipgloss.Style {
			if active {
				return buttonActive
			}

			return buttonInactive
		},
		GetFavoriteButtonStyle: func(isFavorite bool) lipgloss.Style {
			if isFavorite {
				return buttonUnfavour
			}

			return buttonFavour
		},
		GenericButtonStyle: buttonInactive,
		GetEmbeddedSectionStyle: func(focused bool) lipgloss.Style {
			embeddedSectionContainer := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colors.BlurredBorderColor)

			if focused {
				embeddedSectionContainer = embeddedSectionContainer.
					BorderForeground(colors.FocusedBorderColor)
			}

			return embeddedSectionContainer.
				Padding(0, 0, 0, 2).
				UnsetAlign()
		},
		EmbeddedStatusBar: lipgloss.NewStyle().Foreground(colors.FadedStatusColor),
	}
}

func GetLongestTicketLen(taskDefs types.TaskDefinitions) int {
	var maxTicketLength int

	for _, taskDef := range taskDefs {
		maxTicketLength = max(maxTicketLength, len(taskDef.Ticket))
	}

	return maxTicketLength
}

func GetTaskInfo(taskDef types.TaskDefinition, padding int) string {
	var sb strings.Builder
	sb.WriteString(taskDef.Info)
	sb.WriteString(" - ")
	sb.WriteString(taskDef.Ticket)
	sb.WriteString(strings.Repeat(" ", padding+1))
	return sb.String()
}

func GetPaddedTimeTableLines(tt *types.TimeTable, taskDefs *types.TaskDefinitionMap) []string {
	var maxTicketLength int

	for _, entry := range *tt {
		maxTicketLength = max(maxTicketLength, len(entry.Ticket))
	}

	var lines []string

	for _, entry := range *tt {
		taskDef, ok := (*taskDefs)[entry.Ticket]
		utils.Assert(ok, "timetable record is missing from task defs")

		padding := maxTicketLength - len(entry.Ticket)

		var sb strings.Builder
		sb.WriteString(taskDef.Info)
		sb.WriteString(" - ")
		sb.WriteString(taskDef.Ticket)
		sb.WriteString(strings.Repeat(" ", padding+1))

		lines = append(lines, sb.String())
	}

	return lines
}
