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
	newLine       = "\n"
)

func (s *Status) Share() string {
	n := strconv.Itoa(s.Round)
	if string(s.Discovered[:]) != s.Wordle {
		n = "X"
	}

	title := fmt.Sprintf("Wordle %d %s/6*", s.PuzzleNumber, n)

	return title + newLine + s.squaresString()
}

func (s *Status) squaresString() string {
	var finalResult []string

	for _, res := range s.Results {
		var row string
		for _, stat := range res {
			for _, v := range stat {
				switch v {
				case Correct:
					row += correctSquare
				case Present:
					row += presentSquare
				case Absent:
					row += absentSquare
				}
			}
		}
		if row != "" {
			finalResult = append(finalResult, row)
		}
	}

	return strings.Join(finalResult, newLine)
}
