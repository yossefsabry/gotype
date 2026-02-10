package app

import "fmt"

func (r *Renderer) drawResults(model *Model, width, height, keyboardStartY int) {
	resultsTop := model.Layout.FooterY - 2
	resultsBottom := model.Layout.FooterY - 1
	if resultsTop < 0 || resultsBottom < 0 || resultsBottom >= height {
		return
	}
	keyboardBottom := keyboardStartY - 1
	if keyboardStartY < height {
		if kbdHeight := keyboardHeight(); kbdHeight > 0 {
			keyboardBottom = keyboardStartY + kbdHeight - 1
		}
	}
	if resultsTop <= keyboardBottom {
		return
	}
	r.fillLine(resultsTop, width, r.styles.Base)
	r.fillLine(resultsBottom, width, r.styles.Base)
	if !model.Results.Visible || !model.Timer.Finished {
		return
	}
	prefix := "final  net: "
	netValue := fmt.Sprintf("%d", model.Results.NetWPM)
	rest := fmt.Sprintf("  raw: %d  acc: %d%%  cons: %d", model.Results.RawWPM, model.Results.Accuracy, model.Results.Consistency)
	lineLen := len(prefix) + len(netValue) + len(rest)
	startX := (width - lineLen) / 2
	if startX < 0 {
		startX = 0
	}
	netStyle := r.styles.Dim
	if model.Results.Improved || !model.Results.HasBaseline {
		netStyle = r.styles.Accent
	} else if model.Results.Worse {
		netStyle = r.styles.Error
	}
	r.drawString(startX, resultsTop, prefix, r.styles.Dim)
	r.drawString(startX+len(prefix), resultsTop, netValue, netStyle)
	r.drawString(startX+len(prefix)+len(netValue), resultsTop, rest, r.styles.Dim)

	bestLine := fmt.Sprintf("best   wpm: %d  acc: %d%%", model.Results.BestWPM, model.Results.BestAccuracy)
	indicator := ""
	indicatorStyle := r.styles.Dim
	newBest := ""
	if model.Results.HasBaseline {
		if model.Results.Improved {
			indicator = " ^"
			indicatorStyle = r.styles.Accent
			newBest = "  NEW BEST!"
		} else if model.Results.Worse {
			indicator = " v"
			indicatorStyle = r.styles.Error
		} else {
			indicator = " ="
		}
	}
	lineLen = len(bestLine) + len(indicator) + len(newBest)
	startX = (width - lineLen) / 2
	if startX < 0 {
		startX = 0
	}
	r.drawString(startX, resultsBottom, bestLine, r.styles.Dim)
	if indicator != "" {
		r.drawString(startX+len(bestLine), resultsBottom, indicator, indicatorStyle)
	}
	if newBest != "" {
		r.drawString(startX+len(bestLine)+len(indicator), resultsBottom, newBest, r.styles.Accent)
	}
}
