package terminal

import (
	"fmt"
	"strings"

	"github.com/Alvaroalonsobabbel/wordle/pkg/wordle"
)

const (
	absentSquare  = "⬜️"
	correctSquare = "🟩"
	presentSquare = "🟨"
)

type share struct {
	results *[6]wordle.Result
	tries   int
}

func newShare(r *[6]wordle.Result) *share {
	return &share{
		results: r,
	}
}

func (s *share) string() string {
	squares := s.generateEmojiString()
	title := fmt.Sprintf("Wordle %d/6*", s.tries)

	return title + newLine + squares + newLine
}

func (s *share) generateEmojiString() string {
	var finalResult []string

	for _, res := range *s.results {
		var row string
		for _, stat := range res {
			switch stat {
			case wordle.Correct:
				row += correctSquare
			case wordle.Present:
				row += presentSquare
			case wordle.Absent:
				row += absentSquare
			}
		}
		if row != "" {
			finalResult = append(finalResult, row)
			s.tries++
		}
	}

	return strings.Join(finalResult, newLine)
}
