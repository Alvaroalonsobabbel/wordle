package terminal

import (
	"testing"

	"github.com/Alvaroalonsobabbel/wordle/pkg/wordle"
	"github.com/stretchr/testify/assert"
)

func TestKeyboardString(t *testing.T) {
	wordle := wordle.NewGame(wordle.WithCustomWord("CHAIR"))
	kb := NewKB(wordle)

	t.Run("game with no rounds return a basic keyboard", func(t *testing.T) {
		expected := "   Q W E R T Z U I O P\n\r    A S D F G H J K L\n\r   ↩︎ Y X C V B N M ←\n\r"
		assert.Equal(t, expected, kb.string())
	})

	t.Run("game with 1 round highlight the used letters", func(t *testing.T) {
		wordle.Try("HELLO") //nolint: errcheck
		expected := "   Q W \x1b[7m\x1b[90mE\x1b[0m R T Z U I \x1b[7m\x1b[90mO\x1b[0m P\n\r    A S D F G \x1b[7m\x1b[33mH\x1b[0m J K \x1b[7m\x1b[90mL\x1b[0m\n\r   ↩︎ Y X C V B N M ←\n\r"
		assert.Equal(t, expected, kb.string())
	})
}

func TestMapRunes(t *testing.T) {
	var (
		w     = wordle.NewGame(wordle.WithCustomWord("CHAIR"))
		kb    = NewKB(w)
		tests = []struct {
			word string
			want map[rune]int
		}{
			{
				"HELLO", map[rune]int{'H': wordle.Present, 'E': wordle.Absent, 'L': wordle.Absent, 'O': wordle.Absent},
			},
			{
				"CELLO", map[rune]int{'C': wordle.Correct},
			},
			{
				"OPIUM", map[rune]int{'O': wordle.Absent, 'P': wordle.Absent, 'I': wordle.Present, 'U': wordle.Absent, 'M': wordle.Absent},
			},
			{
				"CHAIR", map[rune]int{'C': wordle.Correct, 'H': wordle.Correct, 'A': wordle.Correct, 'I': wordle.Correct, 'R': wordle.Correct},
			},
		}
		assertMapRunes = func(t testing.TB, want map[rune]int, kb *keyboard) {
			t.Helper()

			for k, v := range want {
				got := kb.used[k]
				assert.Equal(t, v, got)
			}
		}
	)

	for _, test := range tests {
		w.Try(test.word) //nolint: errcheck
		kb.mapRunes()
		assertMapRunes(t, test.want, kb)
	}
}
