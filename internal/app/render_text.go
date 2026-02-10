package app

func (r *Renderer) drawText(model *Model, width, height, keyboardStartY int) {
	if width <= 0 || height <= 0 {
		return
	}
	lines := model.linesForWidth(model.Layout.TextWidth)
	if len(lines) == 0 {
		return
	}
	lineSpacing := 1
	maxLines := maxVisibleLines
	areaTop := model.Layout.StatsY + 2
	areaBottom := keyboardStartY - 2
	if areaBottom < areaTop {
		areaBottom = areaTop
	}
	textBlockHeight := (maxLines-1)*lineSpacing + 1
	textStartY := areaTop
	if areaBottom-areaTop+1 >= textBlockHeight {
		textStartY = areaTop + (areaBottom-areaTop+1-textBlockHeight)/2
	}
	cursorIndex := len(model.Text.Typed)
	startLine := defaultStartLine(lines, cursorIndex)
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
		r.drawLine(model, line, lineX, y)
	}
}

func (r *Renderer) drawLine(model *Model, line Line, x, y int) {
	if y < 0 || y >= model.Layout.Height {
		return
	}
	for i := line.Start; i < line.End; i++ {
		ch := model.Text.Target[i]
		style := r.styles.Dim
		if i < len(model.Text.Typed) {
			if model.Text.Typed[i] == ch {
				style = r.styles.Correct
			} else {
				style = r.styles.Error
			}
		}
		if i == len(model.Text.Typed) && !model.Timer.Finished {
			style = r.styles.Cursor
		}
		style = style.Bold(true)
		r.setContent(x+(i-line.Start), y, ch, style)
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
	return end - line.Start
}
