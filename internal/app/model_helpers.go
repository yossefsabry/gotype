package app

import "time"

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

func (m *Model) resetMistakes() {
	if m.Mistakes == nil {
		m.Mistakes = make(map[rune]int, 32)
		return
	}
	for key := range m.Mistakes {
		delete(m.Mistakes, key)
	}
}

func (m *Model) recordMistake(r rune) {
	if m.Mistakes == nil {
		m.Mistakes = make(map[rune]int, 32)
	}
	m.Mistakes[r]++
}

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
