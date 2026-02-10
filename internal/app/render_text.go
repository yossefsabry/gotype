package app

import "github.com/gdamore/tcell/v2"

func (r *Renderer) drawText(model *Model, width, height, keyboardStartY int) {
	if width <= 0 || height <= 0 {
		return
	}
	textWidth := model.Layout.TextWidth / textScale
	if textWidth < 1 {
		textWidth = 1
	}
	lines := model.linesForWidth(textWidth)
	if len(lines) == 0 {
		return
	}
	lineHeight := textScale
	lineSpacing := textScale
	areaTop := model.Layout.TextY
	areaBottom := keyboardStartY - 2
	if areaBottom < areaTop {
		areaBottom = areaTop
	}
	availableHeight := areaBottom - areaTop + 1
	maxLines := maxVisibleLines
	textBlockHeight := (maxLines-1)*lineSpacing + lineHeight
	if availableHeight < textBlockHeight {
		maxLines = availableHeight / lineHeight
		if maxLines < 1 {
			maxLines = 1
		}
		textBlockHeight = (maxLines-1)*lineSpacing + lineHeight
	}
	textStartY := areaTop
	if availableHeight >= textBlockHeight {
		textStartY = areaTop + (availableHeight-textBlockHeight)/2
	}
	cursorIndex := len(model.Text.Typed)
	startLine := defaultStartLine(lines, cursorIndex)
	if !model.Timer.Finished {
		maxStart := len(lines) - maxLines
		if maxStart < 0 {
			maxStart = 0
		}
		if startLine > maxStart {
			startLine = maxStart
		}
	}
	if model.Timer.Finished {
		startLine = model.ReviewStart
		maxStart := len(lines) - maxLines
		if maxStart < 0 {
			maxStart = 0
		}
		if startLine < 0 {
			startLine = 0
		}
		if startLine > maxStart {
			startLine = maxStart
		}
	}
	endLine := startLine + maxLines
	if endLine > len(lines) {
		endLine = len(lines)
	}

	clearTop := model.Layout.StatsY + 1
	if clearTop > areaTop {
		clearTop = areaTop
	}
	for y := clearTop; y <= areaBottom; y++ {
		r.fillLine(y, width, r.styles.Base)
	}

	for i := startLine; i < endLine; i++ {
		line := lines[i]
		y := textStartY + (i-startLine)*lineSpacing
		lineX := r.centeredLineX(model, line)
		r.drawLine(model, line, lineX, y, lineHeight)
	}
}

func (r *Renderer) drawLine(model *Model, line Line, x, y, scale int) {
	if y < 0 || y >= model.Layout.Height {
		return
	}
	for i := line.Start; i < line.End; i++ {
		target := model.Text.Target[i]
		renderCh := target
		style := r.styles.Dim
		if i < len(model.Text.Typed) {
			typed := model.Text.Typed[i]
			if typed == target {
				style = r.styles.Correct
			} else {
				style = r.styles.Error
				if target == ' ' {
					renderCh = '_'
				}
			}
		}
		if i == len(model.Text.Typed) && !model.Timer.Finished {
			style = r.styles.Cursor
		}
		style = style.Bold(true)
		r.drawRuneBlock(x+(i-line.Start)*scale, y, renderCh, style, scale)
	}
}

func (r *Renderer) centeredLineX(model *Model, line Line) int {
	lineLen := lineVisualWidth(model.Text.Target, line)
	if model.Layout.TextWidth <= lineLen {
		return model.Layout.TextX
	}
	return model.Layout.TextX + (model.Layout.TextWidth-lineLen)/2
}

func lineVisualWidth(target []rune, line Line) int {
	end := line.End
	for end > line.Start && target[end-1] == ' ' {
		end--
	}
	if end < line.Start {
		return 0
	}
	return (end - line.Start) * textScale
}

func (r *Renderer) drawRuneBlock(x, y int, ch rune, style tcell.Style, scale int) {
	if scale != textScale {
		for dy := 0; dy < scale; dy++ {
			rowY := y + dy
			for dx := 0; dx < scale; dx++ {
				r.setContent(x+dx, rowY, ch, style)
			}
		}
		return
	}
	glyph, ok := glyphForRune(ch)
	if !ok {
		for dy := 0; dy < scale; dy++ {
			rowY := y + dy
			for dx := 0; dx < scale; dx++ {
				r.setContent(x+dx, rowY, ch, style)
			}
		}
		return
	}
	for dy := 0; dy < scale; dy++ {
		rowY := y + dy
		rowBits := glyph[dy]
		if rowBits == 0 {
			continue
		}
		for dx := 0; dx < scale; dx++ {
			bit := uint8(1 << (scale - 1 - dx))
			if rowBits&bit == 0 {
				continue
			}
			r.setContent(x+dx, rowY, ch, style)
		}
	}
}
