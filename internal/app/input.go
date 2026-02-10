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
	case tcell.KeyCtrlR:
		m.Reset()
		return true, false
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if m.Timer.Finished {
			return false, false
		}
		if m.Backspace(now) {
			return true, false
		}
		return false, false
	case tcell.KeyRune:
		r := event.Rune()
		if r == '\n' || r == '\t' {
			return false, false
		}
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

func (m *Model) HandleClick(x, y int, now time.Time) bool {
	for _, region := range m.Layout.Regions {
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
			return true
		default:
			m.SetMessage("only time mode is ready", now, 2*time.Second)
			return true
		}
	case strings.HasPrefix(id, "time:"):
		if duration, ok := timeByRegion[id]; ok {
			m.Options.Duration = duration
			m.Reset()
			return true
		}
	}
	return false
}
