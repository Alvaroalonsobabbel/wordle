package terminal

import (
	"fmt"
	"io"
	"testing"

	"github.com/Alvaroalonsobabbel/wordle/wordle"
	"github.com/stretchr/testify/assert"
)

type mockReader struct {
	data []byte
}

func (m *mockReader) Read(p []byte) (int, error) {
	n := copy(p, m.data)
	m.data = m.data[n:]
	return len(p), nil
}

func TestRead(t *testing.T) {
	mockReader := &mockReader{}
	terminal := newTestTerminal(io.Discard, mockReader)

	tests := []struct {
		name     string
		input    []byte
		wantExit bool
	}{
		{"ctrl-c exits the game", []byte{ctrlC}, true},
		{"'A' does not exits the game", []byte{'A'}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockReader.data = test.input
			_, exit := terminal.read()
			if test.wantExit {
				assert.True(t, exit)
			} else {
				assert.False(t, exit)
			}
		})
	}
}

func newTestTerminal(w io.Writer, r io.Reader) *terminal { //nolint: revive
	render := newRender(w)
	wordle := &wordle.Status{Wordle: "CHORE"}
	return &terminal{
		reader:   r,
		render:   render,
		wordle:   wordle,
		keyboard: newKeyboard(wordle, render),
		round:    newRound(wordle, render),
	}
}

func TestFinishingMessage(t *testing.T) {
	tests := []struct {
		misses      int
		expectedMsg string
	}{
		{0, "Genius"},
		{1, "Magnificent"},
		{2, "Impressive"},
		{3, "Splendid"},
		{4, "Great"},
		{5, "Phew!"},
		{6, "HELLO"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("guessing in %d attempts returns %s", test.misses+1, test.expectedMsg), func(t *testing.T) {
			wordle := &wordle.Status{Wordle: "HELLO"}
			terminal := New(wordle)
			for range test.misses {
				assert.NoError(t, wordle.Try("CHAIR"))
			}
			if test.misses < 6 {
				assert.NoError(t, wordle.Try("HELLO"))
			}
			assert.Equal(t, test.expectedMsg, terminal.finishingMsg())
		})
	}
}
