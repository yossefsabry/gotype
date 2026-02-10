package app

import "github.com/yossefsabry/gotype/internal/storage"

type ResultsState struct {
	Visible      bool
	NetWPM       int
	RawWPM       int
	Accuracy     int
	Consistency  int
	BestWPM      int
	BestAccuracy int
	HasBaseline  bool
	Improved     bool
	Worse        bool
}

func (m *Model) ResetResults() {
	m.Results = ResultsState{}
}

func (m *Model) FinalizeResults(prevBest storage.BestScore, hasPrev bool) {
	current := ResultsState{
		Visible:     true,
		NetWPM:      m.Stats.WPM,
		RawWPM:      m.Stats.RawWPM,
		Accuracy:    m.Stats.Accuracy,
		Consistency: m.history.StdDev(),
		HasBaseline: hasPrev,
	}
	best := prevBest
	if hasPrev {
		if isBetter(m.Stats, prevBest) {
			current.Improved = true
			best = storage.BestScore{WPM: m.Stats.WPM, Accuracy: m.Stats.Accuracy}
		} else if isWorse(m.Stats, prevBest) {
			current.Worse = true
		}
	} else {
		best = storage.BestScore{WPM: m.Stats.WPM, Accuracy: m.Stats.Accuracy}
	}
	current.BestWPM = best.WPM
	current.BestAccuracy = best.Accuracy
	m.Results = current
}

func isBetter(stats Stats, best storage.BestScore) bool {
	if stats.WPM > best.WPM {
		return true
	}
	if stats.WPM == best.WPM && stats.Accuracy > best.Accuracy {
		return true
	}
	return false
}

func isWorse(stats Stats, best storage.BestScore) bool {
	if stats.WPM < best.WPM {
		return true
	}
	if stats.WPM == best.WPM && stats.Accuracy < best.Accuracy {
		return true
	}
	return false
}
