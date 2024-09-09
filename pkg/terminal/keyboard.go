package terminal

import (
	"fmt"
	"strings"

	"github.com/Alvaroalonsobabbel/wordle/pkg/wordle"
)

type keyboard struct {
	alphabet map[string]string
}

func NewKB() *keyboard { //nolint: revive
	return &keyboard{alphabet: newAlphabetMap()}
}

func (k *keyboard) update(res wordle.Result, word string) {
	status := strings.Split(word, "")
	for i, v := range status {
		switch res[i] {
		case wordle.Present:
			k.alphabet[v] = fmt.Sprintf(yellowBackground, v)
		case wordle.Correct:
			k.alphabet[v] = fmt.Sprintf(greenBackground, v)
		case wordle.Absent:
			k.alphabet[v] = fmt.Sprintf(greyBackground, v)
		}
	}
}

func (k *keyboard) string() string {
	var (
		firstRow  = []string{"Q", "W", "E", "R", "T", "Z", "U", "I", "O", "P"}
		secondRow = []string{"A", "S", "D", "F", "G", "H", "J", "K", "L"}
		thirdRow  = []string{"←", "Y", "X", "C", "V", "B", "N", "M", "↩︎"}
	)

	return "   " + k.renderRow(firstRow) + newLine +
		"    " + k.renderRow(secondRow) + newLine +
		"   " + k.renderRow(thirdRow) + newLine
}

func (k *keyboard) renderRow(row []string) string {
	for i, v := range row {
		row[i] = k.alphabet[v]
	}

	return strings.Join(row, " ")
}

func newAlphabetMap() map[string]string {
	var alphabetMap = make(map[string]string)
	for i := range 26 {
		letter := string(rune('A' + i))
		alphabetMap[letter] = letter
	}
	alphabetMap["←"] = "←"
	alphabetMap["↩︎"] = "↩︎"

	return alphabetMap
}
