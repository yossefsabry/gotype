package app

import (
	"fmt"
	"time"

	"github.com/yossefsabry/gotype/internal/storage"
)

// Persister handles saving data to disk in a non-blocking way
// saving data like best scores and preferences without blocking the main thread
type Persister struct {
	path string
	ch   chan storage.Data
	done chan struct{}
}

// NewPersister creates a new Persister with the given file path and 
// starts the background loop
func NewPersister(path string) *Persister {
	p := &Persister{
		path: path,
		ch:   make(chan storage.Data, 1),
		done: make(chan struct{}),
	}
	go p.loop()
	return p
}

// loop runs in the background and listens for data to save or a signal to stop
func (p *Persister) loop() {
	for {
		select {
		case data := <-p.ch:
			_ = storage.Save(p.path, data)
		case <-p.done:
			for {
				select {
				case data := <-p.ch:
					// save on close to ensure we don't lose any pending data
					_ = storage.Save(p.path, data)
				default:
					return
				}
			}
		}
	}
}

// Save sends data to be saved in the background,
// if the channel is full it will drop the oldest data
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

// close signals the background loop to stop and ensures any pending data is saved
func (p *Persister) Close() {
	close(p.done)
}

// loadPersistedData loads the persisted data from disk and returns the 
// file path and data
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

// saving perferences from the model to the storage format
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

// loading preferences from storage and applying them to the model,
// returns true if any changes were made
func applyPreferences(model *Model, prefs storage.Preferences) bool {
	changed := false
	mode := modeFromString(prefs.Mode)
	if model.Options.Mode != mode {
		model.Options.Mode = mode
		changed = true
	}
	// only apply duration if it's greater than 0 to avoid overriding defaults 
	// with invalid values
	if prefs.DurationSeconds > 0 {
		duration := time.Duration(prefs.DurationSeconds) * time.Second
		if model.Options.Duration != duration {
			model.Options.Duration = duration
			changed = true
		}
	}
	// count is applied if it's only greater then 0 so (when user type)
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

// generte a uniqe key for best score based on options, options -> for 
//  different modes
func scoreKey(options Options) string {
	if options.Mode == ModeWords {
		return fmt.Sprintf("words:%d|punct=%t|numbers=%t", options.WordCount,
			options.Punctuation, options.Numbers)
	}
	return fmt.Sprintf("time:%ds|punct=%t|numbers=%t",
		int(options.Duration.Seconds()), options.Punctuation,
		options.Numbers)
}

// udpate the best score in the data if new stats are better than the current
// stats
func updateBestScore(data *storage.Data, options Options, stats Stats, now time.Time) bool {
	// no score yet
	if data.BestScores == nil {
		data.BestScores = map[string]storage.BestScore{}
	}
	key := scoreKey(options)
	current, ok := data.BestScores[key]
	// compare between the new stats and old stats that is store
	if ok {
		// if one off this is not true so the old score is better than the 
		// new score
		if stats.WPM < current.WPM {
			return false
		}
		if stats.WPM == current.WPM && stats.Accuracy <= current.Accuracy {
			return false
		}
	}
	// update too the new score
	data.BestScores[key] = storage.BestScore{
		WPM:       stats.WPM,
		Accuracy:  stats.Accuracy,
		Timestamp: now.Unix(),
	}
	return true
}

// convert the mode enum to a string for storage
func modeToString(mode Mode) string {
	switch mode {
	case ModeWords:
		return "words"
	default:
		return "time"
	}
}

// convert the mode string from storage back to the enum
func modeFromString(value string) Mode {
	switch value {
	case "words":
		return ModeWords
	case "zen":
		return ModeTime
	default:
		return ModeTime
	}
}
