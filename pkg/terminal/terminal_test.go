package terminal

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	var buffer bytes.Buffer
	terminal := New(&buffer, false, false)

	terminal.render()
	result := buffer.String()
	expectedResult := "\x1b[H\x1b[2J\x1b[1m6 attempts to find a 5-letter word\n\x1b[0m\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\n\r\n\r   Q W E R T Z U I O P\n\r    A S D F G H J K L\n\r   ← Y X C V B N M ↩︎\n\r"
	assert.Equal(t, expectedResult, result)
}

func TestEnter(t *testing.T) {
	var buffer = io.Discard
	t.Run("two letters on first round", func(t *testing.T) {
		terminal := New(buffer, false, false)

		terminal.enter(97) // A
		want := "\tA _ _ _ _"
		assert.Equal(t, want, terminal.rounds[0].string())

		terminal.enter(98) // B
		want = "\tA B _ _ _"
		assert.Equal(t, want, terminal.rounds[0].string())
	})

	t.Run("passing an entire round adds colors to the words", func(t *testing.T) {
		terminal := NewTestTerminal(buffer, "CHORE")
		terminal.enter(99)  // C
		terminal.enter(104) // H
		terminal.enter(97)  // A
		terminal.enter(114) // R
		terminal.enter(109) // M
		terminal.enter(enter)
		terminal.enter(99) // C
		want := "\t\x1b[1m\x1b[32mC\x1b[0m \x1b[1m\x1b[32mH\x1b[0m \x1b[1mA\x1b[0m \x1b[1m\x1b[32mR\x1b[0m \x1b[1mM\x1b[0m"
		assert.Equal(t, want, terminal.rounds[0].string())
		want = "\tC _ _ _ _"
		assert.Equal(t, want, terminal.rounds[1].string())
	})

	t.Run("backspace returns one space", func(t *testing.T) {
		terminal := NewTestTerminal(buffer, "CHORE")
		terminal.enter(99)  // C
		terminal.enter(104) // H
		terminal.enter(backspace)
		terminal.enter(97) // A

		want := "\tC A _ _ _"
		assert.Equal(t, want, terminal.rounds[0].string())
	})

	t.Run("backspace on the first line does nothing", func(t *testing.T) {
		terminal := NewTestTerminal(buffer, "CHORE")
		terminal.enter(backspace)
		want := "\t_ _ _ _ _"
		assert.Equal(t, want, terminal.rounds[0].string())
	})

	t.Run("writing further than 5 letters does nothing", func(t *testing.T) {
		terminal := NewTestTerminal(buffer, "CHORE")
		for range 6 {
			terminal.enter(97) // A
		}
		want := "\tA A A A A"
		assert.Equal(t, want, terminal.rounds[0].string())
	})

	t.Run("not allowed letters are ignored", func(t *testing.T) {
		terminal := NewTestTerminal(buffer, "CHORE")
		terminal.enter(49) // 1

		want := "\t_ _ _ _ _"
		assert.Equal(t, want, terminal.rounds[0].string())
	})
}
