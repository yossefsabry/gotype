package app

import "github.com/gdamore/tcell/v2"

type Renderer struct {
	screen     tcell.Screen
	styles     Styles
	themeID    string
	forceClear bool
	lastWidth  int
	lastHeight int
}

func NewRenderer(screen tcell.Screen) *Renderer {
	defaultTheme := DefaultThemeID()
	return &Renderer{
		screen:  screen,
		styles:  NewStyles(ThemeByID(defaultTheme)),
		themeID: defaultTheme,
	}
}

func (r *Renderer) Render(model *Model) {
	r.syncTheme(model)
	width, height := r.screen.Size()
	if width != r.lastWidth || height != r.lastHeight {
		r.forceClear = true
		r.lastWidth = width
		r.lastHeight = height
	}
	r.screen.SetStyle(r.styles.Base)
	if r.forceClear {
		r.fillScreen(width, height, r.styles.Base)
		r.forceClear = false
	}

	r.drawTopBar(model, width)
	r.drawThemeMenu(model, width)
	r.drawStats(model, width)
	keyboardStartY := r.keyboardStartY(model, height)
	r.drawText(model, width, height, keyboardStartY)
	r.drawKeyboard(model, width, height, keyboardStartY)
	r.drawResults(model, width, height, keyboardStartY)
	r.drawFooter(model, width, height)

	r.screen.Show()
}

func (r *Renderer) syncTheme(model *Model) {
	if model.ThemeID == "" {
		return
	}
	if r.themeID == model.ThemeID {
		return
	}
	r.themeID = model.ThemeID
	r.styles = NewStyles(ThemeByID(model.ThemeID))
	r.forceClear = true
}

func (r *Renderer) fillScreen(width, height int, style tcell.Style) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r.screen.SetContent(x, y, ' ', nil, style)
		}
	}
}

func (r *Renderer) fillLine(y, width int, style tcell.Style) {
	for x := 0; x < width; x++ {
		r.setContent(x, y, ' ', style)
	}
}

func (r *Renderer) panelStyle(style tcell.Style) tcell.Style {
	return style.Background(r.styles.PanelBg)
}

func (r *Renderer) drawString(x, y int, text string, style tcell.Style) {
	for i, ch := range text {
		r.setContent(x+i, y, ch, style)
	}
}

func (r *Renderer) setContent(x, y int, ch rune, style tcell.Style) {
	width, height := r.screen.Size()
	if x < 0 || y < 0 || x >= width || y >= height {
		return
	}
	r.screen.SetContent(x, y, ch, nil, style)
}
