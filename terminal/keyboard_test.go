package terminal

import (
	"fmt"
	"testing"

	"github.com/Alvaroalonsobabbel/wordle/wordle"
	"github.com/stretchr/testify/assert"
)

func TestKeyboardUpdate(t *testing.T) {
	result := wordle.Result{wordle.Absent, wordle.Correct, wordle.Correct, wordle.Absent, wordle.Present}

	kb := NewKB()
	kb.update(result, "HELLO")

	assert.Equal(t, fmt.Sprintf(greyBackground, "H"), kb.alphabet["H"])
	assert.Equal(t, fmt.Sprintf(greenBackground, "E"), kb.alphabet["E"])
	assert.Equal(t, fmt.Sprintf(yellowBackground, "O"), kb.alphabet["O"])
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
}
