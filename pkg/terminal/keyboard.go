package terminal

import (
	"fmt"
	"strings"

	"github.com/Alvaroalonsobabbel/wordle/pkg/wordle"
)

type keyboard struct {
	w     *wordle.Status
	used  map[rune]int
	flash string
}

func NewKB(w *wordle.Status) *keyboard { //nolint: revive
	return &keyboard{
		w:    w,
		used: make(map[rune]int),
	}
}

func (kb *keyboard) string() string {
	kb.mapRunes()

	var (
		firstRow  = []string{"Q", "W", "E", "R", "T", "Z", "U", "I", "O", "P"}
		secondRow = []string{"A", "S", "D", "F", "G", "H", "J", "K", "L"}
		thirdRow  = []string{"↩︎", "Y", "X", "C", "V", "B", "N", "M", "←"}
	)

	return "   " + kb.renderRow(firstRow) + newLine +
		"    " + kb.renderRow(secondRow) + newLine +
		"   " + kb.renderRow(thirdRow) + newLine
}

func (kb *keyboard) renderRow(row []string) string {
	for i, v := range row {
		if v == kb.flash {
			row[i] = fmt.Sprintf(flash, v)
			continue
		}
		stat, ok := kb.used[[]rune(v)[0]]
		if ok {
			var color string
			switch stat {
			case wordle.Correct:
				color = greenBackground
			case wordle.Present:
				color = yellowBackground
			case wordle.Absent:
				color = greyBackground
			}
			row[i] = fmt.Sprintf(color, v)
		}
	}

	return strings.Join(row, " ")
}

func (kb *keyboard) mapRunes() {
	for _, all := range kb.w.Results {
		for _, round := range all {
			for k, v := range round {
				if kb.used[k] != wordle.Correct {
					kb.used[k] = v
				}
			}
		}
	}
}
