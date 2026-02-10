package app

type Line struct {
	Start int
	End   int
}

func buildLines(target []rune, width int) []Line {
	if width <= 0 || len(target) == 0 {
		return nil
	}
	lines := make([]Line, 0, len(target)/width+1)
	start := 0
	for start < len(target) {
		end, next := nextLineEnd(target, start, width)
		if end <= start {
			break
		}
		lines = append(lines, Line{Start: start, End: end})
		start = next
	}
	return lines
}

func nextLineEnd(target []rune, start, width int) (int, int) {
	max := start + width
	if max > len(target) {
		max = len(target)
	}
	if max <= start {
		return start, start
	}
	breakAt := -1
	for i := start; i < max; i++ {
		if target[i] == ' ' {
			breakAt = i
		}
	}
	if breakAt != -1 && breakAt+1 <= max {
		end := breakAt + 1
		return end, end
	}
	return max, max
}

func lineIndexFor(lines []Line, index int) int {
	if len(lines) == 0 {
		return 0
	}
	for i, line := range lines {
		if index >= line.Start && index < line.End {
			return i
		}
	}
	return len(lines) - 1
}
