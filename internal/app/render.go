package app

import "github.com/gdamore/tcell/v2"

// handles all rendering logic 
type Renderer struct {
	screen     tcell.Screen
	styles     Styles
	themeID    string
	forceClear bool
	lastWidth  int
	lastHeight int
}

func NewRenderer(screen tcell.Screen) *Renderer {
	// get default theme
	defaultTheme := DefaultThemeID()
	// passing the style for detault style too render structure
	return &Renderer{
		screen:  screen,
		styles:  NewStyles(ThemeByID(defaultTheme)),
		themeID: defaultTheme,
	}
}

func (r *Renderer) Render(model *Model) {
	// check for updates every render and update the theme if needed
	r.syncTheme(model)

	width, height := r.screen.Size()
	// if the size of the terminal has changed since the last render, 
	// 	we need to clear the screen to avoid render shits
	if width != r.lastWidth || height != r.lastHeight {
		r.forceClear = true
		r.lastWidth = width
		r.lastHeight = height
	}

	// the base style is used to fill the screen with the 
	// background color of the theme,
	r.screen.SetStyle(r.styles.Base)
	// clear and render (refresh)
	if r.forceClear {
		r.fillScreen(width, height, r.styles.Base)
		r.forceClear = false
	}

	// so if the timer is not active then we need to render the top bar and theme menu
	focus := model.focusActive()
	if !focus {
		r.fillLine(0, width, r.styles.Base)
		r.drawTopBar(model, width)
		r.drawThemeMenu(model, width)
	}

	// if the timer is active then we need to render the stats, text, keyboard and results
	r.drawStats(model, width)
	keyboardStartY := r.keyboardStartY(model, height)
	r.drawText(model, width, height, keyboardStartY)
	r.drawKeyboard(model, width, height, keyboardStartY)
	r.drawResults(model, width, height, keyboardStartY)
	r.drawFooter(model, width, height)

	r.screen.Show()
}

func (r *Renderer) syncTheme(model *Model) {

	if model.ThemeID == "" { return }
	if r.themeID == model.ThemeID { return }

	// get the id for the theme and update the renderer
	r.themeID = model.ThemeID
	// get the sytles for theme
	r.styles = NewStyles(ThemeByID(model.ThemeID)) 

	// force clear the screen to apply new theme
	r.forceClear = true
}

// fill the screen with the given style, used to clear the screen when the 
// size changes or theme changes
func (r *Renderer) fillScreen(width, height int, style tcell.Style) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r.screen.SetContent(x, y, ' ', nil, style)
		}
	}
}

// fill a line with the given style, used to clear the top bar when 
// the timer is not active
func (r *Renderer) fillLine(y, width int, style tcell.Style) {
	for x := 0; x < width; x++ {
		r.setContent(x, y, ' ', style)
	}
}

// get the background style for the panels based on the theme, used to render the
func (r *Renderer) panelStyle(style tcell.Style) tcell.Style {
	return style.Background(r.styles.PanelBg)
}

// drawing string at the given x,y values
func (r *Renderer) drawString(x, y int, text string, style tcell.Style) {
	for i, ch := range text {
		r.setContent(x+i, y, ch, style)
	}
}

// set content at the given x,y values with the given character and style,
func (r *Renderer) setContent(x, y int, ch rune, style tcell.Style) {
	width, height := r.screen.Size()
	if x < 0 || y < 0 || x >= width || y >= height {
		return
	}
	r.screen.SetContent(x, y, ch, nil, style)
}
