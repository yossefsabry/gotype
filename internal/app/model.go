package app

import (
	"math"
	"math/rand"
	"time"
)

type Mode int

// for the two modes
const (
	ModeTime Mode = iota
	ModeWords
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
	RawWPM    int
	Accuracy  int
	Streak    int
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
	Options           Options
	Timer             Timer
	Stats             Stats
	Text              Text
	Generator         *Generator
	Layout            Layout
	UI                UIState
	ThemeID           string
	ThemeMenu         bool
	LastKey           rune
	LastKeyAt         time.Time
	Results           ResultsState
	ReviewStart       int
	Mistakes          map[rune]int
	history           StatsHistory
	lineCache         LineCache
	targetVersion     int
	lastDerivedSecond int64
}

// this value for the inital word count
const (
	initialWordCount     = 220
	extendWordCount      = 80
	keyHighlightDuration = 450 * time.Millisecond
)

// creating the model
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

// updating the model to the initial state
func (m *Model) Reset() {
	if m.Options.Mode == ModeWords {
		m.Text.Target = m.Generator.Build(m.Options.WordCount, m.Options)
	} else {
		m.Text.Target = m.Generator.Build(initialWordCount, m.Options)
	}
	m.bumpTargetVersion()
	m.Text.Typed = m.Text.Typed[:0]
	if m.Options.Mode == ModeWords {
		m.Timer = Timer{}
	} else {
		m.Timer = Timer{Remaining: m.Options.Duration}
	}
	m.Stats = Stats{}
	m.ResetResults()
	m.ResetReview()
	m.resetMistakes()
	m.history.Reset()
	m.lastDerivedSecond = -1
	m.LastKey = 0
	m.UpdateDerived(time.Now())
	m.syncLayoutFocus()
}

// when start typing the timer starts
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
	m.syncLayoutFocus()
}

// update the mode data every second and when needed
func (m *Model) Update(now time.Time) bool {
	changed := false
	if m.Options.Mode == ModeTime && m.Timer.Started && m.Timer.Running {
		remaining := m.Timer.End.Sub(now)
		if remaining <= 0 {
			m.Timer.Remaining = 0
			m.Timer.Running = false
			m.Timer.Finished = true
			changed = true
		} else {
			remaining = remaining.Truncate(time.Second)
			if remaining != m.Timer.Remaining {
				m.Timer.Remaining = remaining
				changed = true
			}
		}
	}
	if m.Timer.Started {
		currentSecond := int64(m.elapsedForStats(now) / time.Second)
		if currentSecond != m.lastDerivedSecond {
			if m.UpdateDerived(now) {
				changed = true
			}
		}
	}
	if m.UI.Message != "" && now.After(m.UI.MessageUntil) {
		m.UI.Message = ""
		changed = true
	}
	if m.LastKey != 0 && now.Sub(m.LastKeyAt) > keyHighlightDuration {
		m.LastKey = 0
		changed = true
	}
	if m.syncLayoutFocus() {
		changed = true
	}
	return changed
}

// when user type char add it to the text and update the prograph view
func (m *Model) AddRune(r rune, now time.Time) {
	index := len(m.Text.Typed)
	m.ensureTarget(index + 1)
	expected := m.Text.Target[index]
	m.Text.Typed = append(m.Text.Typed, r)
	if r == expected {
		m.Stats.Correct++
		m.Stats.Streak++
	} else {
		m.Stats.Incorrect++
		m.Stats.Streak = 0
		m.recordMistake(normalizeRune(r))
	}
	if m.Options.Mode == ModeWords && len(m.Text.Typed) >= len(m.Text.Target) {
		m.Timer.Finished = true
		m.Timer.Running = false
		m.Timer.End = now
	}
	m.UpdateDerived(now)
	m.syncLayoutFocus()
}

// remove the last char and update the prograph view
func (m *Model) Backspace(now time.Time) bool {
	if len(m.Text.Typed) == 0 {
		return false
	}
	index := len(m.Text.Typed) - 1
	m.removeTypedRange(index, index+1)
	m.UpdateDerived(now)
	return true
}

// remove the last word and update the prograph view
func (m *Model) BackspaceWord(now time.Time) bool {
	if len(m.Text.Typed) == 0 {
		return false
	}
	end := len(m.Text.Typed)
	start := end
	for start > 0 && m.Text.Typed[start-1] == ' ' {
		start--
	}
	for start > 0 && m.Text.Typed[start-1] != ' ' {
		start--
	}
	if start < end {
		m.removeTypedRange(start, end)
		m.UpdateDerived(now)
		return true
	}
	return false
}

// calc the WPM and accuracy every second and when needed
func (m *Model) UpdateDerived(now time.Time) bool {
	newAccuracy := 100
	if total := m.Stats.Correct + m.Stats.Incorrect; total > 0 {
		newAccuracy = int(math.Round(float64(m.Stats.Correct) / float64(total) * 100))
	}
	newWPM := 0
	newRawWPM := 0
	if m.Timer.Started {
		elapsed := m.elapsedForStats(now)
		m.lastDerivedSecond = int64(elapsed / time.Second)
		minutes := elapsed.Minutes()
		newWPM = int(float64(m.Stats.Correct)/5.0/minutes + 0.5)
		newRawWPM = int(float64(m.Stats.Correct+m.Stats.Incorrect)/5.0/minutes + 0.5)
		m.history.Record(elapsed, newWPM)
	} else {
		m.lastDerivedSecond = -1
	}
	changed := newAccuracy != m.Stats.Accuracy || newWPM != m.Stats.WPM || newRawWPM != m.Stats.RawWPM
	m.Stats.Accuracy = newAccuracy
	m.Stats.WPM = newWPM
	m.Stats.RawWPM = newRawWPM
	return changed
}

// set the message to show in the UI for a certain duration
func (m *Model) SetMessage(text string, now time.Time, duration time.Duration) {
	m.UI.Message = text
	m.UI.MessageUntil = now.Add(duration)
}

// ensure the target text is long enough to type, if not extend it
func (m *Model) ensureTarget(minLength int) {
	if len(m.Text.Target) >= minLength {
		return
	}
	if m.Options.Mode == ModeWords {
		return
	}
	m.Text.Target = m.Generator.Extend(m.Text.Target, extendWordCount, m.Options)
	m.bumpTargetVersion()
}

// remove the typed chars in the given range
func (m *Model) removeTypedRange(start, end int) {
	if start < 0 {
		start = 0
	}
	if end > len(m.Text.Typed) {
		end = len(m.Text.Typed)
	}
	if start >= end {
		return
	}
	for i := start; i < end; i++ {
		expected := m.Text.Target[i]
		if m.Text.Typed[i] == expected {
			m.Stats.Correct--
		} else {
			m.Stats.Incorrect--
		}
	}
	copy(m.Text.Typed[start:], m.Text.Typed[end:])
	m.Text.Typed = m.Text.Typed[:len(m.Text.Typed)-(end-start)]
	m.recalculateStreak()
}

// recalculate the current streak after removing chars
// this is needed to update the streak after backspacing
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

// set the theme and return true if it changed
func (m *Model) SetTheme(id string) bool {
	if m.ThemeID == id {
		return false
	}
	m.ThemeID = id
	return true
}
