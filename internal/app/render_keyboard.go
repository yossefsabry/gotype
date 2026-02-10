package app

import "github.com/gdamore/tcell/v2"

func (r *Renderer) drawKeyboard(model *Model, width, height, keyboardStartY int) {
	if len(keyboardRows) == 0 {
		return
	}
	startY := keyboardStartY
	keyGap := 2
	for rowIndex, row := range keyboardRows {
		rowWidth := keyboardRowWidth(row, keyGap)
		startX := (width - rowWidth) / 2
		if startX < 0 {
			startX = 0
		}
		x := startX
		y := startY + rowIndex
		if y < 0 || y >= height {
			continue
		}
		r.fillLine(y, width, r.styles.Base)
		for _, key := range row {
			style := r.styles.Key
			if key.Rune != 0 {
				if key.Rune == model.LastKey {
					style = r.styles.KeyActive
				} else if model.Mistakes != nil && model.Mistakes[key.Rune] > 0 {
					style = r.styles.KeyError
				}
			}
			r.drawKey(x, y, key, style)
			x += key.Width + keyGap
		}
	}
}

func keyboardRowWidth(row []Key, keyGap int) int {
	if len(row) == 0 {
		return 0
	}
	width := 0
	for i, key := range row {
		if i > 0 {
			width += keyGap
		}
		width += key.Width
	}
	return width
}

func (r *Renderer) drawKey(x, y int, key Key, style tcell.Style) {
	for i := 0; i < key.Width; i++ {
		r.setContent(x+i, y, ' ', style)
	}
	labelX := x + (key.Width-len(key.Label))/2
	r.drawString(labelX, y, key.Label, style)
}

func (r *Renderer) keyboardStartY(model *Model, height int) int {
	keyboardHeight := len(keyboardRows)
	if keyboardHeight == 0 {
		return model.Layout.FooterY
	}
	gap := 2
	startY := model.Layout.FooterY - keyboardHeight - gap
	minY := model.Layout.StatsY + 2
	if startY < minY {
		startY = minY
	}
	if startY+keyboardHeight >= height {
		startY = height - keyboardHeight
	}
	if startY < 0 {
		startY = 0
	}
	return startY
}
