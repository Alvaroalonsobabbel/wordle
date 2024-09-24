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
		wordle.Try("HELLO")
		expected := "   Q W \x1b[7m\x1b[90mE\x1b[0m R T Z U I \x1b[7m\x1b[90mO\x1b[0m P\n\r    A S D F G \x1b[7m\x1b[33mH\x1b[0m J K \x1b[7m\x1b[90mL\x1b[0m\n\r   ↩︎ Y X C V B N M ←\n\r"
		assert.Equal(t, expected, kb.string())
	})
}

func TestMapRunes(t *testing.T) {
	w := wordle.NewGame(wordle.WithCustomWord("CHAIR"))
	kb := NewKB(w)

	word := "HELLO"
	wantStatus := []int{wordle.Present, wordle.Absent, wordle.Absent, wordle.Absent, wordle.Absent}
	w.Try(word)

	kb.mapRunes()

	for i, v := range word {
		rune, ok := kb.used[v]
		assert.True(t, ok)
		assert.Equal(t, wantStatus[i], rune)
	}
}
