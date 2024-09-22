package wordle

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	absentSquare  = "⬜️"
	correctSquare = "🟩"
	presentSquare = "🟨"
	newLine       = "\n\r"
)

func (g *Game) Share() string {
	n := strconv.Itoa(g.Round)
	if string(g.discovered[:]) != g.wordle {
		n = "X"
	}

	title := fmt.Sprintf("Wordle %d %s/6*", g.wordleNumber, n)

	return title + newLine + g.squaresString()
}

func (g *Game) squaresString() string {
	var finalResult []string

	for _, res := range g.Results {
		var row string
		for _, stat := range res {
			switch stat {
			case Correct:
				row += correctSquare
			case Present:
				row += presentSquare
			case Absent:
				row += absentSquare
			}
		}
		if row != "" {
			finalResult = append(finalResult, row)
		}
	}

	return strings.Join(finalResult, newLine)
}
