package app

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

func (l *Layout) Recalculate(width, height int, mode Mode) {
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
	add := func(id string) {
		label := labelForRegion(id, mode)
		l.Regions = append(l.Regions, Region{ID: id, X: x, Y: l.TopY, Width: len(label)})
		x += len(label) + 2
	}

	add("opt:punct")
	add("opt:numbers")
	l.Separators = append(l.Separators, x)
	x += 3

	for _, id := range modeOrder {
		add(id)
	}
	l.Separators = append(l.Separators, x)
	x += 3

	for _, id := range selectorOrder {
		add(id)
	}
	l.Separators = append(l.Separators, x)
	x += 3

	add("btn:themes")

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
	"btn:themes":  "themes",
}

var modeOrder = []string{
	"mode:time",
	"mode:words",
}

func labelForRegion(id string, mode Mode) string {
	if label, ok := selectorLabel(id, mode); ok {
		return label
	}
	if label, ok := regionLabels[id]; ok {
		return label
	}
	return id
}
