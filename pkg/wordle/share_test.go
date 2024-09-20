package wordle

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShareString(t *testing.T) {
	t.Run("win in two tries", func(t *testing.T) {
		wordle := NewGame(WithCustomWord("HELLO"))
		wordle.Try("CELLO") //nolint: errcheck
		wordle.Try("HELLO") //nolint: errcheck

		got := wordle.Share()
		want := "Wordle 0 2/6*" + newLine +
			absentSquare + strings.Repeat(correctSquare, 4) +
			newLine + strings.Repeat(correctSquare, 5) + newLine

		assert.Equal(t, want, got)
	})

	t.Run("win in 6 tries", func(t *testing.T) {
		wordle := NewGame(WithCustomWord("LIGHT"))
		tries := []string{"SCARF", "MIGHT", "FIGHT", "TIGHT", "RIGHT", "LIGHT"}
		for _, word := range tries {
			err := wordle.Try(word)
			assert.NoError(t, err)
		}

		got := wordle.Share()
		want := "Wordle 0 6/6*" + newLine +
			strings.Repeat(absentSquare, 5) + newLine +
			strings.Repeat(absentSquare+strings.Repeat(correctSquare, 4)+newLine, 4) +
			strings.Repeat(correctSquare, 5) + newLine

		assert.Equal(t, want, got)
	})

	t.Run("lose", func(t *testing.T) {
		wordle := NewGame(WithCustomWord("LIGHT"))
		word := "SCARF"
		for range 6 {
			err := wordle.Try(word)
			assert.NoError(t, err)
		}

		got := wordle.Share()
		want := "Wordle 0 X/6*" + newLine + strings.Repeat(strings.Repeat(absentSquare, 5)+newLine, 6)

		assert.Equal(t, want, got)
	})
}

func TestGenerateEmojiString(t *testing.T) {
	tests := []struct {
		name    string
		results [6]Result
		want    string
	}{
		{
			"all present",
			[6]Result{
				{Present, Present, Present, Present, Present},
			},
			strings.Repeat(presentSquare, 5),
		},
		{
			"all correct",
			[6]Result{
				{Correct, Correct, Correct, Correct, Correct},
			},
			strings.Repeat(correctSquare, 5),
		},
		{
			"all absent",
			[6]Result{
				{Absent, Absent, Absent, Absent, Absent},
			},
			strings.Repeat(absentSquare, 5),
		},
		{
			"misc one row",
			[6]Result{
				{Absent, Present, Correct, Present, Absent},
			},
			absentSquare + presentSquare + correctSquare + presentSquare + absentSquare,
		},
		{
			"misc two rows",
			[6]Result{
				{Absent, Present, Correct, Present, Absent},
				{Correct, Correct, Correct, Correct, Correct},
			},
			absentSquare + presentSquare + correctSquare + presentSquare + absentSquare + newLine +
				strings.Repeat(correctSquare, 5),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wordle := &Game{Results: test.results}
			got := wordle.generateEmojiString()
			assert.Equal(t, test.want, got)
		})
	}
}
