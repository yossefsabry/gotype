package app

import "github.com/gdamore/tcell/v2"

const (
	keyboardKeyGap        = 1
	keyboardMinKeyGap     = 0
	keyboardRowGap        = 0
	keyboardFooterGap     = 2
	keyboardMinKeyWidth   = 2
	keyboardMinKeyPadding = 1
)

func (r *Renderer) drawKeyboard(model *Model, width, height, keyboardStartY int) {
	if len(keyboardRows) == 0 || width <= 0 || height <= 0 {
		return
	}
	startY := keyboardStartY
	for rowIndex, row := range keyboardRows {
		rowWidths, keyGap, rowWidth := keyboardRowLayout(row, width)
		startX := (width - rowWidth) / 2
		if startX < 0 {
			startX = 0
		}
		x := startX
		y := startY + rowIndex*(1+keyboardRowGap)
		if y < 0 || y >= height {
			continue
		}
		r.fillLine(y, width, r.styles.Base)
		for i, key := range row {
			style := r.styles.Key
			if key.Rune != 0 {
				if key.Rune == model.LastKey {
					style = r.styles.KeyActive
				} else if model.Mistakes != nil && model.Mistakes[key.Rune] > 0 {
					style = r.styles.KeyError
				}
			}
			keyWidth := key.Width
			if i < len(rowWidths) {
				keyWidth = rowWidths[i]
			}
			r.drawKey(x, y, key, keyWidth, style)
			x += keyWidth + keyGap
		}
	}
}

func keyboardHeight() int {
	if len(keyboardRows) == 0 {
		return 0
	}
	return len(keyboardRows) + (len(keyboardRows)-1)*keyboardRowGap
}

func keyboardRowLayout(row []Key, availableWidth int) ([]int, int, int) {
	if len(row) == 0 {
		return nil, keyboardKeyGap, 0
	}
	widths := make([]int, len(row))
	baseSum := 0
	baseWidth := row[0].Width
	sameWidth := true
	for i, key := range row {
		widths[i] = key.Width
		baseSum += key.Width
		if key.Width != baseWidth {
			sameWidth = false
		}
	}
	keyGap := keyboardKeyGap
	rowWidth := baseSum + keyGap*(len(row)-1)
	if rowWidth <= availableWidth {
		return widths, keyGap, rowWidth
	}
	if keyGap > keyboardMinKeyGap {
		keyGap = keyboardMinKeyGap
		rowWidth = baseSum + keyGap*(len(row)-1)
		if rowWidth <= availableWidth {
			return widths, keyGap, rowWidth
		}
	}
	if sameWidth {
		availableForKeys := availableWidth - keyGap*(len(row)-1)
		if availableForKeys <= 0 {
			return widths, keyGap, rowWidth
		}
		minWidth := minKeyWidth(row[0])
		keyWidth := availableForKeys / len(row)
		if keyWidth < minWidth {
			keyWidth = minWidth
		}
		if keyWidth > baseWidth {
			keyWidth = baseWidth
		}
		for i := range widths {
			widths[i] = keyWidth
		}
		rowWidth = keyWidth*len(row) + keyGap*(len(row)-1)
		return widths, keyGap, rowWidth
	}
	minWidths := make([]int, len(row))
	minSum := 0
	for i, key := range row {
		minWidth := minKeyWidth(key)
		minWidths[i] = minWidth
		minSum += minWidth
	}
	availableForKeys := availableWidth - keyGap*(len(row)-1)
	if availableForKeys <= minSum {
		rowWidth = minSum + keyGap*(len(row)-1)
		return minWidths, keyGap, rowWidth
	}
	shrink := baseSum - availableForKeys
	for shrink > 0 {
		reduced := false
		for i := range widths {
			if widths[i] > minWidths[i] {
				widths[i]--
				shrink--
				reduced = true
				if shrink == 0 {
					break
				}
			}
		}
		if !reduced {
			break
		}
	}
	rowWidth = sumKeyWidths(widths) + keyGap*(len(row)-1)
	return widths, keyGap, rowWidth
}

func minKeyWidth(key Key) int {
	minWidth := len(key.Label) + keyboardMinKeyPadding
	if minWidth < keyboardMinKeyWidth {
		minWidth = keyboardMinKeyWidth
	}
	if minWidth > key.Width {
		return key.Width
	}
	return minWidth
}

func sumKeyWidths(widths []int) int {
	if len(widths) == 0 {
		return 0
	}
	sum := 0
	for _, width := range widths {
		sum += width
	}
	return sum
}

func (r *Renderer) drawKey(x, y int, key Key, width int, style tcell.Style) {
	for i := 0; i < width; i++ {
		r.setContent(x+i, y, ' ', style)
	}
	labelX := x + (width-len(key.Label))/2
	r.drawString(labelX, y, key.Label, style)
}

func (r *Renderer) keyboardStartY(model *Model, height int) int {
	keyboardHeight := keyboardHeight()
	if keyboardHeight == 0 {
		return model.Layout.FooterY
	}
	gap := keyboardFooterGap
	minY := model.Layout.StatsY + 2
	available := model.Layout.FooterY - gap - minY
	if available < keyboardHeight {
		return height
	}
	startY := model.Layout.FooterY - keyboardHeight - gap
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
