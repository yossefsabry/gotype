package app

import "github.com/gdamore/tcell/v2"

type Styles struct {
	Base    tcell.Style
	Panel   tcell.Style
	Dim     tcell.Style
	Accent  tcell.Style
	Correct tcell.Style
	Error   tcell.Style
	Cursor  tcell.Style
	PanelBg tcell.Color
}

func NewStyles() Styles {
	bg := hexColor(0x2c2e31)
	panel := hexColor(0x323437)
	text := hexColor(0xd1d0c5)
	dim := hexColor(0x646669)
	accent := hexColor(0xe2b714)
	errorColor := hexColor(0xca4754)

	base := tcell.StyleDefault.Background(bg).Foreground(text)
	return Styles{
		Base:    base,
		Panel:   tcell.StyleDefault.Background(panel).Foreground(text),
		Dim:     base.Foreground(dim),
		Accent:  base.Foreground(accent),
		Correct: base.Foreground(text),
		Error:   base.Foreground(errorColor),
		Cursor:  base.Background(accent).Foreground(bg),
		PanelBg: panel,
	}
}

func hexColor(value int32) tcell.Color {
	return tcell.NewHexColor(value)
}
