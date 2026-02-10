package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

func (r *Renderer) drawTopBar(model *Model, width int) {
	y := model.Layout.TopY
	r.fillLine(y, width, r.styles.Panel)

	for _, x := range model.Layout.Separators {
		r.drawString(x, y, "|", r.panelStyle(r.styles.Dim))
	}

	for _, region := range model.Layout.Regions {
		label := labelForRegion(region.ID, model.Options.Mode)
		style := r.styleForRegion(model, region.ID)
		r.drawString(region.X, region.Y, label, r.panelStyle(style))
	}
}

func (r *Renderer) drawThemeMenu(model *Model, width int) {
	if !model.Layout.MenuOpen {
		r.fillLine(model.Layout.TopY+1, width, r.styles.Base)
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
	status := "time: " + formatDuration(model.Options.Duration)
	if model.Options.Mode == ModeWords {
		status = fmt.Sprintf("words: %d", model.WordsLeft())
	} else if model.Timer.Started {
		status = "time: " + formatDuration(model.Timer.Remaining)
	}
	chars := len(model.Text.Typed)
	stats := fmt.Sprintf("%s  wpm: %d  acc: %d%%  ch: %d  streak: %d  %s", label, model.Stats.WPM, model.Stats.Accuracy, chars, model.Stats.Streak, status)
	if model.Timer.Finished {
		stats = fmt.Sprintf("finished  wpm: %d  acc: %d%%  ch: %d", model.Stats.WPM, model.Stats.Accuracy, chars)
	}
	r.fillLine(model.Layout.StatsY, width, r.styles.Base)
	x := (width - len(stats)) / 2
	if x < 0 {
		x = 0
	}
	r.drawString(x, model.Layout.StatsY, stats, r.styles.Dim)
}

func (r *Renderer) drawFooter(model *Model, width, height int) {
	message := " type to start <tab> reset  <ctrl+w> del word  <esc> quit "
	if model.Timer.Finished {
		message = " finished <tab> restart  <ctrl+w> del word  <esc> quit  up/down review "
	}
	if model.UI.Message != "" {
		message = model.UI.Message
	}
	r.fillLine(model.Layout.FooterY, width, r.styles.Base)
	x := (width - len(message)) / 2
	if x < 0 {
		x = 0
	}
	r.drawString(x, model.Layout.FooterY, message, r.styles.Dim)
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
		if model.Options.Mode == ModeTime {
			return r.styles.Accent
		}
		return r.styles.Dim
	case "mode:words":
		if model.Options.Mode == ModeWords {
			return r.styles.Accent
		}
		return r.styles.Dim
	default:
		if strings.HasPrefix(id, "theme:") {
			themeID, ok := ThemeIDFromRegion(id)
			if ok && model.ThemeID == themeID {
				return r.styles.Accent
			}
			return r.styles.Dim
		}
		if option, ok := selectorByID(id); ok {
			if model.Options.Mode == ModeWords {
				if model.Options.WordCount == option.WordCount {
					return r.styles.Accent
				}
				return r.styles.Dim
			}
			if model.Options.Duration == option.Duration {
				return r.styles.Accent
			}
			return r.styles.Dim
		}
	}
	return r.styles.Dim
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
