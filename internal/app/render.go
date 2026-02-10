package app

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Renderer struct {
	screen     tcell.Screen
	styles     Styles
	themeID    string
	forceClear bool
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
	r.screen.SetStyle(r.styles.Base)
	if r.forceClear {
		r.fillScreen(width, height, r.styles.Base)
		r.forceClear = false
	} else {
		r.screen.Clear()
	}

	r.drawTopBar(model, width)
	r.drawThemeMenu(model, width)
	r.drawStats(model, width)
	r.drawText(model, width, height)
	r.drawFooter(model, width, height)

	r.screen.Show()
}

func (r *Renderer) drawTopBar(model *Model, width int) {
	y := model.Layout.TopY
	r.fillLine(y, width, r.styles.Panel)

	for _, x := range model.Layout.Separators {
		r.drawString(x, y, "|", r.panelStyle(r.styles.Dim))
	}

	for _, region := range model.Layout.Regions {
		label := regionLabels[region.ID]
		style := r.styleForRegion(model, region.ID)
		r.drawString(region.X, region.Y, label, r.panelStyle(style))
	}
}

func (r *Renderer) drawThemeMenu(model *Model, width int) {
	if !model.Layout.MenuOpen {
		return
	}
	y := model.Layout.MenuY
	r.fillLine(y, width, r.styles.Panel)
	for _, region := range model.Layout.MenuRegions {
		themeID, ok := ThemeIDFromRegion(region.ID)
		if !ok {
			continue
		}
		label := ThemeLabel(themeID)
		style := r.styleForRegion(model, region.ID)
		r.drawString(region.X, region.Y, label, r.panelStyle(style))
	}
}

func (r *Renderer) drawStats(model *Model, width int) {
	label := "english"
	timeLeft := formatDuration(model.Options.Duration)
	if model.Timer.Started {
		timeLeft = formatDuration(model.Timer.Remaining)
	}
	stats := fmt.Sprintf("%s  wpm: %d  acc: %d%%  time: %s", label, model.Stats.WPM, model.Stats.Accuracy, timeLeft)
	if model.Timer.Finished {
		stats = fmt.Sprintf("finished  wpm: %d  acc: %d%%", model.Stats.WPM, model.Stats.Accuracy)
	}
	x := (width - len(stats)) / 2
	if x < 0 {
		x = 0
	}
	r.drawString(x, model.Layout.StatsY, stats, r.styles.Dim)
}

func (r *Renderer) drawText(model *Model, width, height int) {
	if width <= 0 || height <= 0 {
		return
	}
	lines := buildLines(model.Text.Target, model.Layout.TextWidth)
	if len(lines) == 0 {
		return
	}
	maxLines := height - model.Layout.TextY - 3
	if maxLines < 1 {
		return
	}
	cursorIndex := len(model.Text.Typed)
	activeLine := lineIndexFor(lines, cursorIndex)
	startLine := 0
	if activeLine > 2 {
		startLine = activeLine - 2
	}
	if startLine+maxLines > len(lines) {
		startLine = int(math.Max(0, float64(len(lines)-maxLines)))
	}
	endLine := startLine + maxLines
	if endLine > len(lines) {
		endLine = len(lines)
	}

	for i := startLine; i < endLine; i++ {
		line := lines[i]
		y := model.Layout.TextY + (i - startLine)
		r.drawLine(model, line, model.Layout.TextX, y)
	}
}

func (r *Renderer) drawFooter(model *Model, width, height int) {
	message := "type to start  r to reset  esc to quit"
	if model.Timer.Finished {
		message = "finished  r to restart  esc to quit"
	}
	if model.UI.Message != "" {
		message = model.UI.Message
	}
	x := (width - len(message)) / 2
	if x < 0 {
		x = 0
	}
	r.drawString(x, model.Layout.FooterY, message, r.styles.Dim)
}

func (r *Renderer) drawLine(model *Model, line Line, x, y int) {
	if y < 0 || y >= model.Layout.Height {
		return
	}
	for i := line.Start; i < line.End; i++ {
		ch := model.Text.Target[i]
		style := r.styles.Dim
		if i < len(model.Text.Typed) {
			if model.Text.Typed[i] == ch {
				style = r.styles.Correct
			} else {
				style = r.styles.Error
			}
		}
		if i == len(model.Text.Typed) && !model.Timer.Finished {
			style = r.styles.Cursor
		}
		r.setContent(x+(i-line.Start), y, ch, style)
	}
}

func (r *Renderer) styleForRegion(model *Model, id string) tcell.Style {
	switch id {
	case "opt:punct":
		if model.Options.Punctuation {
			return r.styles.Accent
		}
		return r.styles.Dim
	case "opt:numbers":
		if model.Options.Numbers {
			return r.styles.Accent
		}
		return r.styles.Dim
	case "btn:themes":
		if model.ThemeMenu {
			return r.styles.Accent
		}
		return r.styles.Dim
	case "mode:time":
		return r.styles.Accent
	case "mode:words", "mode:quote", "mode:zen", "mode:custom":
		return r.styles.Dim
	default:
		if strings.HasPrefix(id, "theme:") {
			themeID, ok := ThemeIDFromRegion(id)
			if ok && model.ThemeID == themeID {
				return r.styles.Accent
			}
			return r.styles.Dim
		}
		if duration, ok := timeByRegion[id]; ok {
			if model.Options.Duration == duration {
				return r.styles.Accent
			}
			return r.styles.Dim
		}
	}
	return r.styles.Dim
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

func formatDuration(duration time.Duration) string {
	if duration < 0 {
		duration = 0
	}
	if duration >= time.Minute {
		minutes := int(duration.Round(time.Second).Minutes())
		return fmt.Sprintf("%dm", minutes)
	}
	seconds := int(duration.Round(time.Second).Seconds())
	return fmt.Sprintf("%ds", seconds)
}
