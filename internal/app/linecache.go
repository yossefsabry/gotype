package app

type LineCache struct {
	width   int
	version int
	lines   []Line
}

func (m *Model) bumpTargetVersion() {
	m.targetVersion++
}

func (m *Model) linesForWidth(width int) []Line {
	if width <= 0 {
		return nil
	}
	if m.lineCache.width == width && m.lineCache.version == m.targetVersion {
		return m.lineCache.lines
	}
	lines := buildLines(m.Text.Target, width)
	m.lineCache.width = width
	m.lineCache.version = m.targetVersion
	m.lineCache.lines = lines
	return lines
}
