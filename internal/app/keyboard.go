package app

type Key struct {
	Label string
	Rune  rune
	Width int
}

var keyboardRows = [][]Key{
	newKeyRow("qwertyuiop[]"),
	newKeyRow("asdfghjkl;'"),
	newKeyRow("zxcvbnm,./"),
	{
		newKey("space", ' ', 12),
	},
}

func newKeyRow(chars string) []Key {
	row := make([]Key, 0, len(chars))
	for _, ch := range chars {
		row = append(row, newKey(string(ch), ch, 0))
	}
	return row
}

func newKey(label string, key rune, width int) Key {
	if width < len(label)+2 {
		width = len(label) + 2
	}
	if width < 3 {
		width = 3
	}
	return Key{Label: label, Rune: key, Width: width}
}
