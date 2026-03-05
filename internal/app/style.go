package app

import "github.com/gdamore/tcell/v2"

// styles for different UI elements based on theme
type Styles struct {
	Base      tcell.Style
	Panel     tcell.Style
	Dim       tcell.Style
	Accent    tcell.Style
	Correct   tcell.Style
	Error     tcell.Style
	Cursor    tcell.Style
	Key       tcell.Style
	KeyActive tcell.Style
	KeyError  tcell.Style
	PanelBg   tcell.Color
}

// create styles from theme colors
func NewStyles(theme Theme) Styles {
	base := tcell.StyleDefault.Background(theme.Background).Foreground(theme.Text)
	panel := tcell.StyleDefault.Background(theme.Panel).Foreground(theme.Text)
	return Styles{
		Base:      base,
		Panel:     panel,
		Dim:       base.Foreground(theme.Dim),
		Accent:    base.Foreground(theme.Accent),
		Correct:   base.Foreground(theme.Text),
		Error:     base.Foreground(theme.Error),
		Cursor:    base.Background(theme.Accent).Foreground(theme.CursorText),
		Key:       tcell.StyleDefault.Background(theme.KeyBackground).Foreground(theme.KeyText),
		KeyActive: tcell.StyleDefault.Background(theme.KeyActiveBg).Foreground(theme.KeyActiveText),
		KeyError:  tcell.StyleDefault.Background(theme.KeyBackground).Foreground(theme.Error),
		PanelBg:   theme.Panel,
	}
}

// convert hex too tcell color
// for remapping the colors too the nearest colors that is supported by the terminal
func hexColor(value int32) tcell.Color {
	return tcell.NewHexColor(value)
}
