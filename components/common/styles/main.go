package styles

import (
	"goLogMyTime/components/common/colors"
	"goLogMyTime/components/common/globals"
	"image/color"

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
	BgColor                       color.Color
	TextHighContrast              lipgloss.Style
	Tutorial                      lipgloss.Style
	SectionTitle                  lipgloss.Style
	GetSectionStyle               func(focused bool) lipgloss.Style
	GetEmbeddedSectionStyle       func(focused bool) lipgloss.Style
	GetCreateStartStopButtonStyle func(state int) lipgloss.Style
	GetStartStopButtonStyle       func(active bool) lipgloss.Style
	EmbeddedStatusBar             lipgloss.Style
}

func GetStyle(paletteIndex int) Style {
	palette := colors.Palettes[paletteIndex]
	colors := colors.ColorsHexToColors(palette.ColorsHex)

	sectionTitle := lipgloss.NewStyle().
		Foreground(colors.TitleColor).
		Bold(true).
		PaddingBottom(1)

	tutorial := lipgloss.NewStyle().
		Foreground(colors.BaseTextColor).
		Align(lipgloss.Left, lipgloss.Top)

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
		Bold(true)

	colorHighContrast := lipgloss.Color("#FFF")

	if palette.Mode == "light" {
		colorHighContrast = lipgloss.Color("#000")
	}

	textHighContrast := lipgloss.NewStyle().Foreground(colorHighContrast)

	return Style{
		BgColor:          colors.BgColor,
		TextHighContrast: textHighContrast,
		Tutorial:         tutorial,
		SectionTitle:     sectionTitle,
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
		GetCreateStartStopButtonStyle: func(state int) lipgloss.Style {
			buttonStyle := buttonDisabled

			if state == globals.TASK_INPUT_EMPTY {
				buttonStyle = buttonInactive
			}

			if state == globals.TASK_INPUT_VALID_ACTIVE {
				buttonStyle = buttonActive
			}

			return buttonStyle.Width(Sizes.ButtonWidth)
		},
		GetStartStopButtonStyle: func(active bool) lipgloss.Style {
			if active {
				return buttonActive
			}

			return buttonInactive
		},
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
