package app

import (
	"math"
	"time"
)

type StatsHistory struct {
	lastSampleSecond int64
	count            int
	sum              float64
	sumSquares       float64
}

func (h *StatsHistory) Reset() {
	h.lastSampleSecond = -1
	h.count = 0
	h.sum = 0
	h.sumSquares = 0
}

func (h *StatsHistory) Record(elapsed time.Duration, wpm int) {
	if elapsed < 0 {
		return
	}
	second := int64(elapsed / time.Second)
	if second <= h.lastSampleSecond {
		return
	}
	h.lastSampleSecond = second
	value := float64(wpm)
	h.count++
	h.sum += value
	h.sumSquares += value * value
}

func (h *StatsHistory) StdDev() int {
	if h.count == 0 {
		return 0
	}
	mean := h.sum / float64(h.count)
	variance := h.sumSquares/float64(h.count) - mean*mean
	if variance < 0 {
		variance = 0
	}
	return int(math.Round(math.Sqrt(variance)))
}
