package app

func (m *Model) ResetReview() {
	m.ReviewStart = 0
}

func (m *Model) InitReviewStart() {
	lines := buildLines(m.Text.Target, m.Layout.TextWidth)
	if len(lines) == 0 {
		m.ReviewStart = 0
		return
	}
	m.ReviewStart = defaultStartLine(lines, len(m.Text.Typed))
}

func (m *Model) ScrollReview(delta int) bool {
	if !m.Timer.Finished {
		return false
	}
	lines := buildLines(m.Text.Target, m.Layout.TextWidth)
	if len(lines) == 0 {
		return false
	}
	maxStart := len(lines) - maxVisibleLines
	if maxStart < 0 {
		maxStart = 0
	}
	start := m.ReviewStart + delta
	if start < 0 {
		start = 0
	}
	if start > maxStart {
		start = maxStart
	}
	if start == m.ReviewStart {
		return false
	}
	m.ReviewStart = start
	return true
}

func (m *Model) ReviewTop() bool {
	if !m.Timer.Finished {
		return false
	}
	if m.ReviewStart == 0 {
		return false
	}
	m.ReviewStart = 0
	return true
}

func (m *Model) ReviewBottom() bool {
	if !m.Timer.Finished {
		return false
	}
	lines := buildLines(m.Text.Target, m.Layout.TextWidth)
	if len(lines) == 0 {
		return false
	}
	maxStart := len(lines) - maxVisibleLines
	if maxStart < 0 {
		maxStart = 0
	}
	if m.ReviewStart == maxStart {
		return false
	}
	m.ReviewStart = maxStart
	return true
}
