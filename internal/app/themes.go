package app

import "github.com/gdamore/tcell/v2"

type Theme struct {
	ID         string
	Label      string
	Background tcell.Color
	Panel      tcell.Color
	Text       tcell.Color
	Dim        tcell.Color
	Accent     tcell.Color
	Error      tcell.Color
	CursorText tcell.Color
}

var themeOptions = []Theme{
	{
		ID:         "rose-pine",
		Label:      "rose-pine",
		Background: hexColor(0x191724),
		Panel:      hexColor(0x1f1d2e),
		Text:       hexColor(0xe0def4),
		Dim:        hexColor(0x6e6a86),
		Accent:     hexColor(0xebbcba),
		Error:      hexColor(0xeb6f92),
		CursorText: hexColor(0x191724),
	},
	{
		ID:         "dark",
		Label:      "dark",
		Background: hexColor(0x2c2e31),
		Panel:      hexColor(0x323437),
		Text:       hexColor(0xd1d0c5),
		Dim:        hexColor(0x646669),
		Accent:     hexColor(0xe2b714),
		Error:      hexColor(0xca4754),
		CursorText: hexColor(0x2c2e31),
	},
	{
		ID:         "light",
		Label:      "light",
		Background: hexColor(0xf5f5f5),
		Panel:      hexColor(0xe9e9e9),
		Text:       hexColor(0x2b2b2b),
		Dim:        hexColor(0x6f6f6f),
		Accent:     hexColor(0x1f6feb),
		Error:      hexColor(0xd73a49),
		CursorText: hexColor(0xf5f5f5),
	},
	{
		ID:         "caption",
		Label:      "captsion",
		Background: hexColor(0x232323),
		Panel:      hexColor(0x2b2b2b),
		Text:       hexColor(0xdadada),
		Dim:        hexColor(0x8a8a8a),
		Accent:     hexColor(0x8ec07c),
		Error:      hexColor(0xfb4934),
		CursorText: hexColor(0x232323),
	},
	{
		ID:         "transparent",
		Label:      "transparent",
		Background: tcell.ColorDefault,
		Panel:      tcell.ColorDefault,
		Text:       hexColor(0xe0def4),
		Dim:        hexColor(0x6e6a86),
		Accent:     hexColor(0xebbcba),
		Error:      hexColor(0xeb6f92),
		CursorText: tcell.ColorBlack,
	},
	{
		ID:         "forest",
		Label:      "forest",
		Background: hexColor(0x0f1f1b),
		Panel:      hexColor(0x142822),
		Text:       hexColor(0xd2e6d6),
		Dim:        hexColor(0x6b8b7c),
		Accent:     hexColor(0x7fd1b9),
		Error:      hexColor(0xe26d5c),
		CursorText: hexColor(0x0f1f1b),
	},
}

const themeRegionPrefix = "theme:"

func ThemeOptions() []Theme {
	return themeOptions
}

func DefaultThemeID() string {
	if len(themeOptions) == 0 {
		return "dark"
	}
	return themeOptions[0].ID
}

func ThemeByID(id string) Theme {
	for _, theme := range themeOptions {
		if theme.ID == id {
			return theme
		}
	}
	if len(themeOptions) == 0 {
		return Theme{}
	}
	return themeOptions[0]
}

func ThemeLabel(id string) string {
	for _, theme := range themeOptions {
		if theme.ID == id {
			return theme.Label
		}
	}
	return id
}

func ThemeRegionID(id string) string {
	return themeRegionPrefix + id
}

func ThemeIDFromRegion(region string) (string, bool) {
	if len(region) <= len(themeRegionPrefix) || region[:len(themeRegionPrefix)] != themeRegionPrefix {
		return "", false
	}
	return region[len(themeRegionPrefix):], true
}
