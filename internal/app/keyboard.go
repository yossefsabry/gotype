package app

type Key struct {
	Label string
	Rune  rune
	Width int
}

const (
	keyboardKeyWidth  = 3
	keyboardSpaceWide = 18
)

// keyboardRows defines the layout of the on-screen keyboard. Each row is a slice of Key structs,
var keyboardRows = [][]Key{
	newKeyRow("qwertyuiop[]"),
	newKeyRow("asdfghjkl;'"),
	newKeyRow("zxcvbnm,./"),
	{
		newKey("space", ' ', keyboardSpaceWide),
	},
}

// newKeyRow creates a slice of Key structs for a given string of characters.
// Each character is converted into a Key with the specified width.
func newKeyRow(chars string) []Key {
	row := make([]Key, 0, len(chars))
	for _, ch := range chars {
		row = append(row, newKey(string(ch), ch, keyboardKeyWidth))
	}
	return row
}

// newKey creates a Key struct with the given label, rune,
// and width. It ensures that the width is at least
func newKey(label string, key rune, width int) Key {
	if width < len(label)+2 {
		width = len(label) + 2
	}
	if width < 3 {
		width = 3
	}
	return Key{Label: label, Rune: key, Width: width}
}
