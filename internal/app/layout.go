package app

import "time"

type Region struct {
	ID    string
	X     int
	Y     int
	Width int
}

func (r Region) Contains(x, y int) bool {
	return y == r.Y && x >= r.X && x < r.X+r.Width
}

type Layout struct {
	Width       int
	Height      int
	TopY        int
	MenuY       int
	StatsY      int
	TextY       int
	TextWidth   int
	TextX       int
	FooterY     int
	MenuOpen    bool
	Regions     []Region
	MenuRegions []Region
	Separators  []int
}

func (l *Layout) Recalculate(width, height int) {
	l.Width = width
	l.Height = height
	if l.TopY == 0 {
		l.TopY = 1
	}
	menuOffset := 0
	if l.MenuOpen {
		menuOffset = 1
		l.MenuY = l.TopY + 1
	} else {
		l.MenuY = 0
	}
	l.StatsY = l.TopY + 2 + menuOffset
	l.TextY = l.TopY + 4 + menuOffset
	l.FooterY = height - 2

	textWidth := width - 8
	if textWidth > 90 {
		textWidth = 90
	}
	if textWidth < 30 {
		textWidth = width - 4
	}
	if textWidth < 10 {
		textWidth = width
	}
	l.TextWidth = textWidth
	l.TextX = (width - textWidth) / 2

	l.Regions = l.Regions[:0]
	l.MenuRegions = l.MenuRegions[:0]
	l.Separators = l.Separators[:0]

	x := 2
	add := func(id, label string) {
		l.Regions = append(l.Regions, Region{ID: id, X: x, Y: l.TopY, Width: len(label)})
		x += len(label) + 2
	}

	add("opt:punct", regionLabels["opt:punct"])
	add("opt:numbers", regionLabels["opt:numbers"])
	l.Separators = append(l.Separators, x)
	x += 3

	for _, id := range modeOrder {
		add(id, regionLabels[id])
	}
	l.Separators = append(l.Separators, x)
	x += 3

	for _, id := range timeOrder {
		add(id, regionLabels[id])
	}
	l.Separators = append(l.Separators, x)
	x += 3

	add("btn:themes", regionLabels["btn:themes"])

	if l.MenuOpen {
		x := 2
		for _, theme := range ThemeOptions() {
			label := theme.Label
			l.MenuRegions = append(l.MenuRegions, Region{ID: ThemeRegionID(theme.ID), X: x, Y: l.MenuY, Width: len(label)})
			x += len(label) + 2
		}
	}
}

var regionLabels = map[string]string{
	"opt:punct":   "@ punctuation",
	"opt:numbers": "# numbers",
	"mode:time":   "time",
	"mode:words":  "words",
	"mode:quote":  "quote",
	"mode:zen":    "zen",
	"mode:custom": "custom",
	"time:30s":    "30s",
	"time:60s":    "60s",
	"time:10m":    "10m",
	"time:30m":    "30m",
	"btn:themes":  "themes",
}

var modeOrder = []string{
	"mode:time",
	"mode:words",
	"mode:quote",
	"mode:zen",
	"mode:custom",
}

var timeOrder = []string{
	"time:30s",
	"time:60s",
	"time:10m",
	"time:30m",
}

var timeByRegion = map[string]time.Duration{
	"time:30s": 30 * time.Second,
	"time:60s": 60 * time.Second,
	"time:10m": 10 * time.Minute,
	"time:30m": 30 * time.Minute,
}
