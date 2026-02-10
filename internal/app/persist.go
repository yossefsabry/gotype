package app

import (
	"fmt"
	"time"

	"github.com/yossefsabry/gotype/internal/storage"
)

type Persister struct {
	path string
	ch   chan storage.Data
	done chan struct{}
}

func NewPersister(path string) *Persister {
	p := &Persister{
		path: path,
		ch:   make(chan storage.Data, 1),
		done: make(chan struct{}),
	}
	go p.loop()
	return p
}

func (p *Persister) loop() {
	for {
		select {
		case data := <-p.ch:
			_ = storage.Save(p.path, data)
		case <-p.done:
			for {
				select {
				case data := <-p.ch:
					_ = storage.Save(p.path, data)
				default:
					return
				}
			}
		}
	}
}

func (p *Persister) Save(data storage.Data) {
	select {
	case p.ch <- data:
	default:
		select {
		case <-p.ch:
		default:
		}
		p.ch <- data
	}
}

func (p *Persister) Close() {
	close(p.done)
}

func loadPersistedData() (string, storage.Data) {
	path, err := storage.DefaultPath()
	if err != nil {
		return "", storage.Data{}
	}
	data, err := storage.Load(path)
	if err != nil {
		return path, storage.Data{}
	}
	if data.BestScores == nil {
		data.BestScores = map[string]storage.BestScore{}
	}
	return path, data
}

func preferencesFromModel(model *Model) storage.Preferences {
	return storage.Preferences{
		ThemeID:         model.ThemeID,
		Mode:            modeToString(model.Options.Mode),
		DurationSeconds: int(model.Options.Duration.Seconds()),
		WordCount:       model.Options.WordCount,
		Punctuation:     model.Options.Punctuation,
		Numbers:         model.Options.Numbers,
	}
}

func applyPreferences(model *Model, prefs storage.Preferences) bool {
	changed := false
	mode := modeFromString(prefs.Mode)
	if model.Options.Mode != mode {
		model.Options.Mode = mode
		changed = true
	}
	if prefs.DurationSeconds > 0 {
		duration := time.Duration(prefs.DurationSeconds) * time.Second
		if model.Options.Duration != duration {
			model.Options.Duration = duration
			changed = true
		}
	}
	if prefs.WordCount > 0 {
		if model.Options.WordCount != prefs.WordCount {
			model.Options.WordCount = prefs.WordCount
			changed = true
		}
	}
	if model.Options.Punctuation != prefs.Punctuation {
		model.Options.Punctuation = prefs.Punctuation
		changed = true
	}
	if model.Options.Numbers != prefs.Numbers {
		model.Options.Numbers = prefs.Numbers
		changed = true
	}
	if prefs.ThemeID != "" {
		theme := ThemeByID(prefs.ThemeID)
		if theme.ID != "" && model.ThemeID != theme.ID {
			model.ThemeID = theme.ID
		}
	}
	return changed
}

func scoreKey(options Options) string {
	if options.Mode == ModeWords {
		return fmt.Sprintf("words:%d|punct=%t|numbers=%t", options.WordCount, options.Punctuation, options.Numbers)
	}
	return fmt.Sprintf("time:%ds|punct=%t|numbers=%t", int(options.Duration.Seconds()), options.Punctuation, options.Numbers)
}

func updateBestScore(data *storage.Data, options Options, stats Stats, now time.Time) bool {
	if data.BestScores == nil {
		data.BestScores = map[string]storage.BestScore{}
	}
	key := scoreKey(options)
	current, ok := data.BestScores[key]
	if ok {
		if stats.WPM < current.WPM {
			return false
		}
		if stats.WPM == current.WPM && stats.Accuracy <= current.Accuracy {
			return false
		}
	}
	data.BestScores[key] = storage.BestScore{
		WPM:       stats.WPM,
		Accuracy:  stats.Accuracy,
		Timestamp: now.Unix(),
	}
	return true
}

func modeToString(mode Mode) string {
	if mode == ModeWords {
		return "words"
	}
	return "time"
}

func modeFromString(value string) Mode {
	if value == "words" {
		return ModeWords
	}
	return ModeTime
}
