package terminal

import (
	"fmt"
	"testing"

	"github.com/Alvaroalonsobabbel/wordle/pkg/wordle"
	"github.com/stretchr/testify/assert"
)

func TestKeyboardUpdate(t *testing.T) {
	result := wordle.Result{wordle.Absent, wordle.Correct, wordle.Correct, wordle.Absent, wordle.Present}

	kb := NewKB()
	kb.update(result, "HELLO")

	assert.Equal(t, fmt.Sprintf(greyBackground, "H"), kb.am["H"])
	assert.Equal(t, fmt.Sprintf(greenBackground, "E"), kb.am["E"])
	assert.Equal(t, fmt.Sprintf(yellowBackground, "O"), kb.am["O"])

	t.Run("correct letters get permanently marked as greenn", func(t *testing.T) {
		result := wordle.Result{wordle.Absent, wordle.Absent, wordle.Present, wordle.Absent, wordle.Absent}
		kb.update(result, "STEAM")

		assert.Equal(t, fmt.Sprintf(greenBackground, "E"), kb.am["E"])
	})
}

func TestKeyboardString(t *testing.T) {
	kb := NewKB()
	expected := "   Q W E R T Z U I O P\n\r    A S D F G H J K L\n\r   ← Y X C V B N M ↩︎\n\r"
	assert.Equal(t, expected, kb.string())
}

func TestGenerateAlphabet(t *testing.T) {
	amap := newAlphabetMap()

	assert.Equal(t, "A", amap["A"])
	assert.Equal(t, "Z", amap["Z"])

	_, ok := amap["Ä"]
	assert.False(t, ok)
}
