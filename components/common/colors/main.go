package colors

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

type ColorsHex struct {
	Background string
	Text       string
	Border     string
	Accent     string
	Button     string
	ButtonText string
	Stop       string
	Muted      string
	Focus      string
}

type Colors struct {
	BgColor                 color.Color
	TitleColor              color.Color
	BaseTextColor           color.Color
	FadedStatusColor        color.Color
	BlurredBorderColor      color.Color
	FocusedBorderColor      color.Color
	DisabledButtonColor     color.Color
	DisabledButtonTextColor color.Color
}

func ColorsHexToColors(colorsHex ColorsHex) Colors {
	return Colors{
		BgColor:                 lipgloss.Color(colorsHex.Background),
		TitleColor:              lipgloss.Color(colorsHex.Text),
		BaseTextColor:           lipgloss.Color(colorsHex.Accent),
		FadedStatusColor:        lipgloss.Color(colorsHex.Muted),
		BlurredBorderColor:      lipgloss.Color(colorsHex.Border),
		FocusedBorderColor:      lipgloss.Color(colorsHex.Focus),
		DisabledButtonColor:     lipgloss.Color(colorsHex.Button),
		DisabledButtonTextColor: lipgloss.Color(colorsHex.ButtonText),
	}
}

type Palette struct {
	Name      string
	Mode      string
	ColorsHex ColorsHex
}

var Palettes = []Palette{
	{
		Name: "Rosé Pine Dawn",
		Mode: "light",
		ColorsHex: ColorsHex{
			Background: "#faf4ed",
			Text:       "#575279",
			Border:     "#286983",
			Accent:     "#d7827e",
			Button:     "#907aa9",
			ButtonText: "#faf4ed",
			Stop:       "#b4637a",
			Muted:      "#9893a5",
			Focus:      "#56949f",
		},
	},
	{
		Name: "Catppuccin Latte",
		Mode: "light",
		ColorsHex: ColorsHex{
			Background: "#eff1f5",
			Text:       "#4c4f69",
			Border:     "#1e66f5",
			Accent:     "#8839ef",
			Button:     "#1e66f5",
			ButtonText: "#eff1f5",
			Stop:       "#d20f39",
			Muted:      "#9ca0b0",
			Focus:      "#40a02b",
		},
	},
	{
		Name: "Ayu Mirage",
		Mode: "dark",
		ColorsHex: ColorsHex{
			Background: "#1f2430",
			Text:       "#cccac2",
			Border:     "#73d0ff",
			Accent:     "#ffcc66",
			Button:     "#5ccfe6",
			ButtonText: "#1f2430",
			Stop:       "#f28779",
			Muted:      "#707a8c",
			Focus:      "#d5ff80",
		},
	},
	{
		Name: "Mono Gray",
		Mode: "dark",
		ColorsHex: ColorsHex{
			Background: "#1e1e1e",
			Text:       "#d4d4d4",
			Border:     "#6e6e6e",
			Accent:     "#ffffff",
			Button:     "#3a3a3a",
			ButtonText: "#d4d4d4",
			Stop:       "#b04040",
			Muted:      "#6e6e6e",
			Focus:      "#d4d4d4",
		},
	},
	{
		Name: "Srcery",
		Mode: "dark",
		ColorsHex: ColorsHex{
			Background: "#1c1b19",
			Text:       "#fce8c3",
			Border:     "#0aaeb3",
			Accent:     "#fbb829",
			Button:     "#2c78bf",
			ButtonText: "#fce8c3",
			Stop:       "#ef2f27",
			Muted:      "#918175",
			Focus:      "#519f50",
		},
	},
	{
		Name: "Dusk Forest",
		Mode: "dark",
		ColorsHex: ColorsHex{
			Background: "#1a2321",
			Text:       "#d7e3d4",
			Border:     "#5fa782",
			Accent:     "#e5b567",
			Button:     "#3d7a5e",
			ButtonText: "#ffffff",
			Stop:       "#d9685f",
			Muted:      "#6e857b",
			Focus:      "#a1d48c",
		},
	},
	{
		Name: "Synthwave",
		Mode: "dark",
		ColorsHex: ColorsHex{
			Background: "#262335",
			Text:       "#ffffff",
			Border:     "#b893ce",
			Accent:     "#ff7edb",
			Button:     "#36f9f6",
			ButtonText: "#262335",
			Stop:       "#fe4450",
			Muted:      "#848bbd",
			Focus:      "#fede5d",
		},
	},
	{
		Name: "Ayu Dark",
		Mode: "dark",
		ColorsHex: ColorsHex{
			Background: "#0b0e14",
			Text:       "#bfbdb6",
			Border:     "#59c2ff",
			Accent:     "#ffb454",
			Button:     "#39bae6",
			ButtonText: "#0b0e14",
			Stop:       "#f07178",
			Muted:      "#565b66",
			Focus:      "#aad94c",
		},
	},
	{
		Name: "Night Owl",
		Mode: "dark",
		ColorsHex: ColorsHex{
			Background: "#011627",
			Text:       "#d6deeb",
			Border:     "#82aaff",
			Accent:     "#ecc48d",
			Button:     "#7e57c2",
			ButtonText: "#ffffff",
			Stop:       "#ef5350",
			Muted:      "#637777",
			Focus:      "#addb67",
		},
	},
	{
		Name: "Matrix Green",
		Mode: "dark",
		ColorsHex: ColorsHex{
			Background: "#000000",
			Text:       "#b5ffc5",
			Border:     "#008f11",
			Accent:     "#00ff41",
			Button:     "#003b00",
			ButtonText: "#00ff41",
			Stop:       "#ff3b3b",
			Muted:      "#1f7a35",
			Focus:      "#00ff41",
		},
	},
}
