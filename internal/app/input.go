package app

import (
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

func (m *Model) HandleKey(event *tcell.EventKey, now time.Time) (bool, bool) {
	switch event.Key() {
	case tcell.KeyCtrlC, tcell.KeyEsc:
		return false, true
	case tcell.KeyTab:
		m.Reset()
		return true, false
	case tcell.KeyCtrlW:
		if m.Timer.Finished {
			return false, false
		}
		if m.BackspaceWord(now) {
			return true, false
		}
		return false, false
	case tcell.KeyUp:
		if m.ScrollReview(-1) {
			return true, false
		}
		return false, false
	case tcell.KeyDown:
		if m.ScrollReview(1) {
			return true, false
		}
		return false, false
	case tcell.KeyPgUp:
		if m.ScrollReview(-maxVisibleLines) {
			return true, false
		}
		return false, false
	case tcell.KeyPgDn:
		if m.ScrollReview(maxVisibleLines) {
			return true, false
		}
		return false, false
	case tcell.KeyHome:
		if m.ReviewTop() {
			return true, false
		}
		return false, false
	case tcell.KeyEnd:
		if m.ReviewBottom() {
			return true, false
		}
		return false, false
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if m.Timer.Finished {
			return false, false
		}
		if event.Modifiers()&tcell.ModCtrl != 0 {
			if m.BackspaceWord(now) {
				return true, false
			}
			return false, false
		}
		if m.Backspace(now) {
			return true, false
		}
		return false, false
	case tcell.KeyRune:
		r := event.Rune()
		if r == '\n' {
			return false, false
		}
		m.registerKey(r, now)
		if m.Timer.Finished {
			if r == 'r' || r == 'R' {
				m.Reset()
				return true, false
			}
			return false, false
		}
		if !m.Timer.Started {
			m.StartTimer(now)
		}
		m.AddRune(r, now)
		return true, false
	}
	return false, false
}

func (m *Model) registerKey(r rune, now time.Time) {
	r = normalizeRune(r)
	m.LastKey = r
	m.LastKeyAt = now
}

// HandleClick checks if the click at (x, y) is within any interactive 
// region and applies the corresponding action.
func (m *Model) HandleClick(x, y int, now time.Time) bool {
	for _, region := range m.Layout.Regions {
		if region.Contains(x, y) {
			return m.applyRegion(region.ID, now)
		}
	}
	for _, region := range m.Layout.MenuRegions {
		if region.Contains(x, y) {
			return m.applyRegion(region.ID, now)
		}
	}
	return false
}

func (m *Model) applyRegion(id string, now time.Time) bool {
	switch {
	case id == "opt:punct":
		m.Options.Punctuation = !m.Options.Punctuation
		m.Reset()
		return true
	case id == "opt:numbers":
		m.Options.Numbers = !m.Options.Numbers
		m.Reset()
		return true
	case strings.HasPrefix(id, "mode:"):
		switch id {
		case "mode:time":
			m.Options.Mode = ModeTime
			m.Reset()
			m.Layout.Recalculate(m.Layout.Width, m.Layout.Height, m.Options.Mode, m.focusActive())
			return true
		case "mode:words":
			m.Options.Mode = ModeWords
			m.Reset()
			m.Layout.Recalculate(m.Layout.Width, m.Layout.Height, m.Options.Mode, m.focusActive())
			return true
		}
		return false
	case strings.HasPrefix(id, "sel:"):
		option, ok := selectorByID(id)
		if !ok {
			return false
		}
		if m.Options.Mode == ModeWords {
			m.Options.WordCount = option.WordCount
			m.Reset()
			return true
		}
		m.Options.Duration = option.Duration
		m.Reset()
		return true
	case id == "btn:themes":
		m.ThemeMenu = !m.ThemeMenu
		m.Layout.MenuOpen = m.ThemeMenu
		m.Layout.Recalculate(m.Layout.Width, m.Layout.Height, m.Options.Mode, m.focusActive())
		return true
	case strings.HasPrefix(id, "theme:"):
		themeID, ok := ThemeIDFromRegion(id)
		if !ok {
			return false
		}
		_ = m.SetTheme(themeID)
		m.ThemeMenu = false
		m.Layout.MenuOpen = false
		m.Layout.Recalculate(m.Layout.Width, m.Layout.Height, m.Options.Mode, m.focusActive())
		return true
	}
	return false
}
