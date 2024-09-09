package terminal

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	var buffer bytes.Buffer
	terminal := New(&buffer, false, false)

	terminal.Render()
	result := buffer.String()
	expectedResult := "\x1b[H\x1b[2J\x1b[1m6 attempts to find a 5-letter word\n\x1b[0m\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\n\r\n\r   Q W E R T Z U I O P\n\r    A S D F G H J K L\n\r   ← Y X C V B N M ↩︎\n\r\x1b[H\x1b[2J\x1b[1m6 attempts to find a 5-letter word\n\x1b[0m\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\n\r\n\r   Q W E R T Z U I O P\n\r    A S D F G H J K L\n\r   ← Y X C V B N M ↩︎\n\r"
	assert.Equal(t, expectedResult, result)
}

func TestPlayRound(t *testing.T) {
	t.Run("two letters on first round", func(t *testing.T) {
		var buffer bytes.Buffer
		terminal := New(&buffer, false, false)
		terminal.enter(97) // A
		buffer.Reset()
		terminal.Render()
		result := buffer.String()
		expectedResult := "\x1b[H\x1b[2J\x1b[1m6 attempts to find a 5-letter word\n\x1b[0m\n\r\tA _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\n\r\n\r   Q W E R T Z U I O P\n\r    A S D F G H J K L\n\r   ← Y X C V B N M ↩︎\n\r"
		assert.Equal(t, expectedResult, result)

		buffer.Reset()
		terminal.enter(98) // B
		buffer.Reset()
		terminal.Render()
		result = buffer.String()
		expectedResult = "\x1b[H\x1b[2J\x1b[1m6 attempts to find a 5-letter word\n\x1b[0m\n\r\tA B _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\n\r\n\r   Q W E R T Z U I O P\n\r    A S D F G H J K L\n\r   ← Y X C V B N M ↩︎\n\r"
		assert.Equal(t, expectedResult, result)
	})

	t.Run("passing an entire round adds colors to the words", func(t *testing.T) {
		var buffer bytes.Buffer
		terminal := NewTestTerminal(&buffer, false, false, "CHORE")
		terminal.enter(99)  // C
		terminal.enter(104) // H
		terminal.enter(97)  // A
		terminal.enter(114) // R
		terminal.enter(109) // M
		terminal.enter(enter)
		buffer.Reset()
		terminal.Render()
		result := buffer.String()
		expectedResult := "\x1b[H\x1b[2J\x1b[1m6 attempts to find a 5-letter word\n\x1b[0m\n\r\t\x1b[1m\x1b[32mC\x1b[0m \x1b[1m\x1b[32mH\x1b[0m \x1b[1mA\x1b[0m \x1b[1m\x1b[32mR\x1b[0m \x1b[1mM\x1b[0m\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\n\r\n\r   Q W E \x1b[7m\x1b[32mR\x1b[0m T Z U I O P\n\r    \x1b[7m\x1b[90mA\x1b[0m S D F G \x1b[7m\x1b[32mH\x1b[0m J K L\n\r   ← Y X \x1b[7m\x1b[32mC\x1b[0m V B N \x1b[7m\x1b[90mM\x1b[0m ↩︎\n\r"
		assert.Equal(t, expectedResult, result)

		buffer.Reset()
		terminal.enter(99) // C
		buffer.Reset()
		terminal.Render()
		result = buffer.String()
		expectedResult = "\x1b[H\x1b[2J\x1b[1m6 attempts to find a 5-letter word\n\x1b[0m\n\r\t\x1b[1m\x1b[32mC\x1b[0m \x1b[1m\x1b[32mH\x1b[0m \x1b[1mA\x1b[0m \x1b[1m\x1b[32mR\x1b[0m \x1b[1mM\x1b[0m\n\r\tC _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\n\r\n\r   Q W E \x1b[7m\x1b[32mR\x1b[0m T Z U I O P\n\r    \x1b[7m\x1b[90mA\x1b[0m S D F G \x1b[7m\x1b[32mH\x1b[0m J K L\n\r   ← Y X \x1b[7m\x1b[32mC\x1b[0m V B N \x1b[7m\x1b[90mM\x1b[0m ↩︎\n\r"
		assert.Equal(t, expectedResult, result)
	})

	t.Run("backspace returns one space", func(t *testing.T) {
		var buffer bytes.Buffer
		terminal := NewTestTerminal(&buffer, false, false, "CHORE")
		terminal.enter(99)  // C
		terminal.enter(104) // H
		terminal.enter(backspace)
		terminal.enter(97) // A
		buffer.Reset()
		terminal.Render()
		result := buffer.String()
		expectedResult := "\x1b[H\x1b[2J\x1b[1m6 attempts to find a 5-letter word\n\x1b[0m\n\r\tC A _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\n\r\n\r   Q W E R T Z U I O P\n\r    A S D F G H J K L\n\r   ← Y X C V B N M ↩︎\n\r"
		assert.Equal(t, expectedResult, result)
	})

	t.Run("backspace on the first line does nothing", func(t *testing.T) {
		var buffer bytes.Buffer
		terminal := NewTestTerminal(&buffer, false, false, "CHORE")
		terminal.enter(backspace)
		buffer.Reset()
		terminal.Render()
		result := buffer.String()
		expectedResult := "\x1b[H\x1b[2J\x1b[1m6 attempts to find a 5-letter word\n\x1b[0m\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\n\r\n\r   Q W E R T Z U I O P\n\r    A S D F G H J K L\n\r   ← Y X C V B N M ↩︎\n\r"
		assert.Equal(t, expectedResult, result)
	})

	t.Run("writing further than 5 letters does nothing", func(t *testing.T) {
		var buffer bytes.Buffer
		terminal := NewTestTerminal(&buffer, false, false, "CHORE")
		for range 5 {
			terminal.enter(97) // A
		}
		buffer.Reset()
		terminal.Render()
		result := buffer.String()
		expectedResult := "\x1b[H\x1b[2J\x1b[1m6 attempts to find a 5-letter word\n\x1b[0m\n\r\tA A A A A\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\n\r\n\r   Q W E R T Z U I O P\n\r    A S D F G H J K L\n\r   ← Y X C V B N M ↩︎\n\r"
		assert.Equal(t, expectedResult, result)
	})
}
