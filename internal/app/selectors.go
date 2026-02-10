package app

import "time"

type SelectorOption struct {
	ID        string
	Duration  time.Duration
	WordCount int
	LabelTime string
	LabelWord string
}

var selectorOptions = []SelectorOption{
	{
		ID:        "sel:30s",
		Duration:  30 * time.Second,
		WordCount: 10,
		LabelTime: "30s",
		LabelWord: "10",
	},
	{
		ID:        "sel:60s",
		Duration:  60 * time.Second,
		WordCount: 25,
		LabelTime: "60s",
		LabelWord: "25",
	},
	{
		ID:        "sel:10m",
		Duration:  10 * time.Minute,
		WordCount: 50,
		LabelTime: "10m",
		LabelWord: "50",
	},
	{
		ID:        "sel:30m",
		Duration:  30 * time.Minute,
		WordCount: 100,
		LabelTime: "30m",
		LabelWord: "100",
	},
}

var selectorOrder = []string{
	"sel:30s",
	"sel:60s",
	"sel:10m",
	"sel:30m",
}

func selectorByID(id string) (SelectorOption, bool) {
	for _, option := range selectorOptions {
		if option.ID == id {
			return option, true
		}
	}
	return SelectorOption{}, false
}

func selectorLabel(id string, mode Mode) (string, bool) {
	option, ok := selectorByID(id)
	if !ok {
		return "", false
	}
	if mode == ModeWords {
		return option.LabelWord, true
	}
	return option.LabelTime, true
}
