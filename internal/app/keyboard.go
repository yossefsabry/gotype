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
		newKey("space", ' ', 20),
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
	if width < len(label)+4 {
		width = len(label) + 4
	}
	if width < 5 {
		width = 5
	}
	return Key{Label: label, Rune: key, Width: width}
}
