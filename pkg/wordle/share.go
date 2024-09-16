package wordle

import (
	"fmt"
	"strings"
)

const (
	absentSquare  = "⬜️"
	correctSquare = "🟩"
	presentSquare = "🟨"
	newLine       = "\n\r"
)

func (g *Game) Share() string {
	squares := g.generateEmojiString()
	title := fmt.Sprintf("Wordle %d/6*", g.Round)

	return title + newLine + squares + newLine
}

func (g *Game) generateEmojiString() string {
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
