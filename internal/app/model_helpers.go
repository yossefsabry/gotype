package app

import "time"


// calc the time for stats
func (m *Model) elapsedForStats(now time.Time) time.Duration {
	elapsed := now.Sub(m.Timer.Start)
	if m.Timer.Finished {
		if m.Options.Mode == ModeTime {
			elapsed = m.Options.Duration
		} else if !m.Timer.End.IsZero() {
			elapsed = m.Timer.End.Sub(m.Timer.Start)
		}
	}
	if elapsed < time.Second {
		elapsed = time.Second
	}
	return elapsed
}

// check if the layout should be focused based on the timer state
func (m *Model) focusActive() bool {
	return m.Timer.Started && !m.Timer.Finished
}

// sync the layout focus with the timer state and recalculate if needed
func (m *Model) syncLayoutFocus() bool {
	focus := m.focusActive()
	if m.Layout.Focus == focus {
		return false
	}
	m.Layout.Focus = focus
	if m.Layout.Width <= 0 || m.Layout.Height <= 0 {
		return false
	}
	m.Layout.Recalculate(m.Layout.Width, m.Layout.Height, m.Options.Mode, focus)
	return true
}

// reset the mistakes map for the model
func (m *Model) resetMistakes() {
	if m.Mistakes == nil {
		m.Mistakes = make(map[rune]int, 32)
		return
	}
	for key := range m.Mistakes {
		delete(m.Mistakes, key)
	}
}

// record a mistake for a given rune in the model's mistakes map
func (m *Model) recordMistake(r rune) {
	if m.Mistakes == nil {
		m.Mistakes = make(map[rune]int, 32)
	}
	m.Mistakes[r]++
}

// recalculate the current streak of correct characters in the model's stats
func (m *Model) recalculateStreak() {
	streak := 0
	for i := len(m.Text.Typed) - 1; i >= 0; i-- {
		if m.Text.Typed[i] != m.Text.Target[i] {
			break
		}
		streak++
	}
	m.Stats.Streak = streak
}
