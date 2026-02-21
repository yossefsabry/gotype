package app

type Line struct {
	Start int
	End   int
}

const (
	maxWordsPerLine = 10
	maxVisibleLines = 3
)

func buildLines(target []rune, width int) []Line {
	if width <= 0 || len(target) == 0 {
		return nil
	}
	words := collectWords(target)
	if len(words) == 0 {
		return nil
	}
	lines := make([]Line, 0, len(words)/maxWordsPerLine+1)
	index := 0
	for index < len(words) {
		lineStart := words[index].Start
		lineEnd := words[index].End
		lineWords := 0
		lineLen := 0
		for index < len(words) && lineWords < maxWordsPerLine {
			word := words[index]
			wordLen := word.End - word.Start
			addLen := wordLen
			if lineWords > 0 {
				addLen++
			}
			if lineWords > 0 && lineLen+addLen > width {
				break
			}
			if lineWords > 0 {
				lineLen++
			}
			lineLen += wordLen
			lineEnd = word.End
			if word.End < len(target) && target[word.End] == ' ' {
				lineEnd = word.End + 1
				lineLen++
			}
			lineWords++
			index++
		}
		lines = append(lines, Line{Start: lineStart, End: lineEnd})
	}
	return lines
}

type wordRange struct {
	Start int
	End   int
}

func collectWords(target []rune) []wordRange {
	words := make([]wordRange, 0, 64)
	start := -1
	for i, ch := range target {
		if ch != ' ' {
			if start == -1 {
				start = i
			}
			continue
		}
		if start != -1 {
			words = append(words, wordRange{Start: start, End: i})
			start = -1
		}
	}
	if start != -1 {
		words = append(words, wordRange{Start: start, End: len(target)})
	}
	return words
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

func defaultStartLine(lines []Line, cursorIndex int) int {
	if len(lines) == 0 {
		return 0
	}
	activeLine := lineIndexFor(lines, cursorIndex)
	startLine := 0
	if activeLine > 1 {
		startLine = activeLine - 1
	}
	maxStart := len(lines) - maxVisibleLines
	if maxStart < 0 {
		maxStart = 0
	}
	if startLine > maxStart {
		startLine = maxStart
	}
	return startLine
}
