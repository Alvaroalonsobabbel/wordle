package terminal

import (
	"fmt"
	"strings"
	"time"

	"github.com/Alvaroalonsobabbel/wordle/wordle"
)

const flash = "\x1b[30m\x1b[47m %s \x1b[0m"

type key struct {
	value  string
	column int
	row    int
}

func (k *key) string() string { //nolint: revive
	return fmt.Sprintf("\x1b[%d;%dH%s", k.row, k.column, k.value)
}

type keyboard struct {
	keys   map[string]*key
	wordle *wordle.Status
	render *render
}

func newKeyboard(w *wordle.Status, r *render) *keyboard { //nolint: revive
	return &keyboard{
		wordle: w,
		render: r,
		keys: map[string]*key{
			"Q":  {" Q ", 2, 11},
			"W":  {" W ", 5, 11},
			"E":  {" E ", 8, 11},
			"R":  {" R ", 11, 11},
			"T":  {" T ", 14, 11},
			"Z":  {" Z ", 17, 11},
			"U":  {" U ", 20, 11},
			"I":  {" I ", 23, 11},
			"O":  {" O ", 26, 11},
			"P":  {" P ", 29, 11},
			"A":  {" A ", 3, 12},
			"S":  {" S ", 6, 12},
			"D":  {" D ", 9, 12},
			"F":  {" F ", 12, 12},
			"G":  {" G ", 15, 12},
			"H":  {" H ", 18, 12},
			"J":  {" J ", 21, 12},
			"K":  {" K ", 24, 12},
			"L":  {" L ", 27, 12},
			"↩︎": {" ↩︎ ", 2, 13},
			"Y":  {" Y ", 5, 13},
			"X":  {" X ", 8, 13},
			"C":  {" C ", 11, 13},
			"V":  {" V ", 14, 13},
			"B":  {" B ", 17, 13},
			"N":  {" N ", 20, 13},
			"M":  {" M ", 23, 13},
			"←":  {" ← ", 26, 13},
		},
	}
}

func (kb *keyboard) print() {
	for v, k := range kb.keys {
		if strings.Contains(string(kb.wordle.Used), v) {
			k.value = fmt.Sprintf(greyBackground, v)
		}
		if strings.Contains(string(kb.wordle.Hints), v) {
			k.value = fmt.Sprintf(yellowBackground, v)
		}
		if strings.Contains(string(kb.wordle.Discovered[:]), v) {
			k.value = fmt.Sprintf(greenBackground, v)
		}

		kb.render.string(k.string())
	}
}

func (kb *keyboard) flash(l byte) {
	var char string
	switch l {
	case backspace:
		char = "←"
	case enter:
		char = "↩︎"
	default:
		char = strings.ToUpper(string(l))
	}

	key, ok := kb.keys[char]
	if !ok {
		return
	}

	oldValue := key.value
	key.value = fmt.Sprintf(flash, char)
	kb.render.string(key.string())
	time.AfterFunc(25*time.Millisecond, func() {
		key.value = oldValue
		kb.render.string(key.string())
	})
}
