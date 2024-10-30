package terminal

import (
	"fmt"
	"strings"
	"time"

	"github.com/Alvaroalonsobabbel/wordle/wordle"
)

const flash = "\x1b[30m\x1b[47m %s \x1b[0m"

type keyboard struct {
	wordle    *wordle.Status
	render    *render
	used      map[rune]int
	flashChar string
}

func newKeyboard(w *wordle.Status, r *render) *keyboard { //nolint: revive
	return &keyboard{
		wordle: w,
		render: r,
		used:   make(map[rune]int),
	}
}

func (kb *keyboard) print() {
	kb.mapRunes()

	var (
		firstRow  = []string{"Q", "W", "E", "R", "T", "Z", "U", "I", "O", "P"}
		secondRow = []string{"A", "S", "D", "F", "G", "H", "J", "K", "L"}
		thirdRow  = []string{"↩︎", "Y", "X", "C", "V", "B", "N", "M", "←"}
	)

	kb.render.string("\033[11;2H" + kb.renderRow(firstRow))
	kb.render.string("\033[12;3H" + kb.renderRow(secondRow))
	kb.render.string("\033[13;2H" + kb.renderRow(thirdRow))
}

func (kb *keyboard) renderRow(row []string) string {
	for i, v := range row {
		if v == kb.flashChar {
			row[i] = fmt.Sprintf(flash, v)
		} else {
			row[i] = fmt.Sprintf(emptyChar, v)
		}

		stat, ok := kb.used[[]rune(v)[0]]
		if ok {
			var c string
			switch stat {
			case wordle.Correct:
				c = greenBackground
			case wordle.Present:
				c = yellowBackground
			case wordle.Absent:
				c = greyBackground
			}
			row[i] = fmt.Sprintf(c, v)
		}
	}

	return strings.Join(row, "")
}

func (kb *keyboard) mapRunes() {
	for _, all := range kb.wordle.Results {
		for _, round := range all {
			for k, v := range round {
				if kb.used[k] != wordle.Correct {
					kb.used[k] = v
				}
			}
		}
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

	kb.flashChar = char
	kb.print()
	time.Sleep(25 * time.Millisecond)
	kb.flashChar = ""
	kb.print()
}
