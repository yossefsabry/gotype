package app

type LineCache struct {
	width   int
	version int
	lines   []Line
}

// bumpTargetVersion increments the target version to indicate that the 
// text has changed and the line cache should be invalidated.
func (m *Model) bumpTargetVersion() {
	m.targetVersion++
}

// linesForWidth returns the cached lines for the given 
// width if the cache is valid.
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
