package app

import (
	"math"
	"math/rand"
	"time"
)

type Mode int

const (
	ModeTime Mode = iota
	ModeWords
	ModeQuote
	ModeZen
	ModeCustom
)

type Options struct {
	Punctuation bool
	Numbers     bool
	Mode        Mode
	Duration    time.Duration
	WordCount   int
}

type Timer struct {
	Started   bool
	Running   bool
	Finished  bool
	Start     time.Time
	End       time.Time
	Remaining time.Duration
}

type Stats struct {
	Correct   int
	Incorrect int
	WPM       int
	Accuracy  int
}

type Text struct {
	Target []rune
	Typed  []rune
}

type UIState struct {
	Message      string
	MessageUntil time.Time
}

type Model struct {
	Options   Options
	Timer     Timer
	Stats     Stats
	Text      Text
	Generator *Generator
	Layout    Layout
	UI        UIState
	ThemeID   string
	ThemeMenu bool
}

const (
	initialWordCount = 220
	extendWordCount  = 80
)

func NewModel() *Model {
	model := &Model{
		Options: Options{
			Mode:      ModeTime,
			Duration:  60 * time.Second,
			WordCount: 50,
		},
		Generator: NewGenerator(rand.NewSource(time.Now().UnixNano())),
		ThemeID:   DefaultThemeID(),
	}
	model.Reset()
	return model
}

func (m *Model) Reset() {
	if m.Options.Mode == ModeWords {
		m.Text.Target = m.Generator.Build(m.Options.WordCount, m.Options)
	} else {
		m.Text.Target = m.Generator.Build(initialWordCount, m.Options)
	}
	m.Text.Typed = m.Text.Typed[:0]
	if m.Options.Mode == ModeWords {
		m.Timer = Timer{}
	} else {
		m.Timer = Timer{Remaining: m.Options.Duration}
	}
	m.Stats = Stats{}
	m.UpdateDerived(time.Now())
}

func (m *Model) StartTimer(now time.Time) {
	if m.Timer.Started {
		return
	}
	m.Timer.Started = true
	m.Timer.Running = true
	m.Timer.Finished = false
	m.Timer.Start = now
	if m.Options.Mode == ModeTime {
		m.Timer.End = now.Add(m.Options.Duration)
		m.Timer.Remaining = m.Options.Duration
	}
}

func (m *Model) Update(now time.Time) bool {
	changed := false
	if m.Options.Mode == ModeTime && m.Timer.Started && m.Timer.Running {
		remaining := m.Timer.End.Sub(now)
		if remaining <= 0 {
			m.Timer.Remaining = 0
			m.Timer.Running = false
			m.Timer.Finished = true
			changed = true
		} else if remaining != m.Timer.Remaining {
			m.Timer.Remaining = remaining
			changed = true
		}
	}
	if m.Timer.Started {
		if m.UpdateDerived(now) {
			changed = true
		}
	}
	if m.UI.Message != "" && now.After(m.UI.MessageUntil) {
		m.UI.Message = ""
		changed = true
	}
	return changed
}

func (m *Model) AddRune(r rune, now time.Time) {
	index := len(m.Text.Typed)
	m.ensureTarget(index + 1)
	expected := m.Text.Target[index]
	m.Text.Typed = append(m.Text.Typed, r)
	if r == expected {
		m.Stats.Correct++
	} else {
		m.Stats.Incorrect++
	}
	if m.Options.Mode == ModeWords && len(m.Text.Typed) >= len(m.Text.Target) {
		m.Timer.Finished = true
		m.Timer.Running = false
		m.Timer.End = now
	}
	m.UpdateDerived(now)
}

func (m *Model) Backspace(now time.Time) bool {
	if len(m.Text.Typed) == 0 {
		return false
	}
	index := len(m.Text.Typed) - 1
	expected := m.Text.Target[index]
	if m.Text.Typed[index] == expected {
		m.Stats.Correct--
	} else {
		m.Stats.Incorrect--
	}
	m.Text.Typed = m.Text.Typed[:index]
	m.UpdateDerived(now)
	return true
}

func (m *Model) UpdateDerived(now time.Time) bool {
	newAccuracy := 100
	if total := m.Stats.Correct + m.Stats.Incorrect; total > 0 {
		newAccuracy = int(math.Round(float64(m.Stats.Correct) / float64(total) * 100))
	}
	newWPM := 0
	if m.Timer.Started {
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
		minutes := elapsed.Minutes()
		newWPM = int(float64(m.Stats.Correct)/5.0/minutes + 0.5)
	}
	changed := newAccuracy != m.Stats.Accuracy || newWPM != m.Stats.WPM
	m.Stats.Accuracy = newAccuracy
	m.Stats.WPM = newWPM
	return changed
}

func (m *Model) SetMessage(text string, now time.Time, duration time.Duration) {
	m.UI.Message = text
	m.UI.MessageUntil = now.Add(duration)
}

func (m *Model) ensureTarget(minLength int) {
	if len(m.Text.Target) >= minLength {
		return
	}
	if m.Options.Mode == ModeWords {
		return
	}
	m.Text.Target = m.Generator.Extend(m.Text.Target, extendWordCount, m.Options)
}

func (m *Model) WordsLeft() int {
	if m.Options.Mode != ModeWords {
		return 0
	}
	index := len(m.Text.Typed)
	if index >= len(m.Text.Target) {
		return 0
	}
	words := 0
	inWord := false
	for i := index; i < len(m.Text.Target); i++ {
		if m.Text.Target[i] != ' ' {
			if !inWord {
				words++
				inWord = true
			}
		} else {
			inWord = false
		}
	}
	return words
}

func (m *Model) SetTheme(id string) bool {
	if m.ThemeID == id {
		return false
	}
	m.ThemeID = id
	return true
}
