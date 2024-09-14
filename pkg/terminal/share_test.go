package terminal

import (
	"strings"
	"testing"

	"github.com/Alvaroalonsobabbel/wordle/pkg/wordle"
	"github.com/stretchr/testify/assert"
)

func TestShareString(t *testing.T) {
	t.Run("win in two tries", func(t *testing.T) {
		wordle := wordle.NewTestWordle(false, "HELLO")
		wordle.Try("CELLO") //nolint: errcheck
		wordle.Try("HELLO") //nolint: errcheck

		share := newShare(&wordle.Results)
		got := share.string()
		want := "Wordle 2/6*" + newLine +
			absentSquare + strings.Repeat(correctSquare, 4) +
			newLine + strings.Repeat(correctSquare, 5)

		assert.Equal(t, want, got)
	})

	t.Run("win in 6 tries", func(t *testing.T) {
		wordle := wordle.NewTestWordle(false, "LIGHT")
		tries := []string{"SCARF", "MIGHT", "FIGHT", "TIGHT", "RIGHT", "LIGHT"}
		for _, word := range tries {
			err := wordle.Try(word)
			assert.NoError(t, err)
		}

		share := newShare(&wordle.Results)
		got := share.string()
		want := "Wordle 6/6*" + newLine +
			strings.Repeat(absentSquare, 5) + newLine +
			absentSquare + strings.Repeat(correctSquare, 4) + newLine +
			absentSquare + strings.Repeat(correctSquare, 4) + newLine +
			absentSquare + strings.Repeat(correctSquare, 4) + newLine +
			absentSquare + strings.Repeat(correctSquare, 4) + newLine +
			strings.Repeat(correctSquare, 5)

		assert.Equal(t, want, got)
	})
}

func TestGenerateEmojiString(t *testing.T) {
	tests := []struct {
		name    string
		results *[6]wordle.Result
		want    string
	}{
		{
			"all present",
			&[6]wordle.Result{
				{wordle.Present, wordle.Present, wordle.Present, wordle.Present, wordle.Present},
			},
			strings.Repeat(presentSquare, 5),
		},
		{
			"all correct",
			&[6]wordle.Result{
				{wordle.Correct, wordle.Correct, wordle.Correct, wordle.Correct, wordle.Correct},
			},
			strings.Repeat(correctSquare, 5),
		},
		{
			"all absent",
			&[6]wordle.Result{
				{wordle.Absent, wordle.Absent, wordle.Absent, wordle.Absent, wordle.Absent},
			},
			strings.Repeat(absentSquare, 5),
		},
		{
			"misc one row",
			&[6]wordle.Result{
				{wordle.Absent, wordle.Present, wordle.Correct, wordle.Present, wordle.Absent},
			},
			absentSquare + presentSquare + correctSquare + presentSquare + absentSquare,
		},
		{
			"misc two rows",
			&[6]wordle.Result{
				{wordle.Absent, wordle.Present, wordle.Correct, wordle.Present, wordle.Absent},
				{wordle.Correct, wordle.Correct, wordle.Correct, wordle.Correct, wordle.Correct},
			},
			absentSquare + presentSquare + correctSquare + presentSquare + absentSquare + newLine +
				strings.Repeat(correctSquare, 5),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			share := newShare(test.results)
			got := share.generateEmojiString()
			assert.Equal(t, test.want, got)
		})
	}
}
