package terminal

import (
	"bytes"
	"testing"

	"github.com/Alvaroalonsobabbel/wordle/pkg/wordle"
	"github.com/stretchr/testify/assert"
)

func TestRenderAll(t *testing.T) {
	buffer, screen := newScreenTest()

	want := "\x1b[H\x1b[2J\x1b[1m6 attempts to find a 5-letter word\n\x1b[0m\n\r\t _  _  _  _  _ \n\r\t _  _  _  _  _ \n\r\t _  _  _  _  _ \n\r\t _  _  _  _  _ \n\r\t _  _  _  _  _ \n\r\t _  _  _  _  _ \n\r\n\r\n\r  Q  W  E  R  T  Z  U  I  O  P \n\r   A  S  D  F  G  H  J  K  L \n\r  ↩︎  Y  X  C  V  B  N  M  ← \n\r\n\r\n\r"
	screen.renderAll()
	assert.Equal(t, want, buffer.String())
}

func TestRenderRound(t *testing.T) {
	buffer, screen := newScreenTest()

	want := "\x1b[3;0H\t A  _  _  _  _ "
	screen.add("A")
	screen.renderRound()
	assert.Equal(t, want, buffer.String())
}

func TestRenderKB(t *testing.T) {
	buffer, screen := newScreenTest()

	want := "\x1b[11;0H  Q  W  E  R  T  Z  U  I  O  P \n\r   A  S  D  F  G  H  J  K  L \n\r  ↩︎  Y  X  C  V  B  N  M  ← \n\r"
	screen.renderKB()
	assert.Equal(t, want, buffer.String())
}

func TestRenderMsg(t *testing.T) {
	buffer, screen := newScreenTest()

	want := "\033[10;0HMESSAGE"
	screen.msg = "MESSAGE"
	screen.renderMsg()
	assert.Equal(t, want, buffer.String())
}

func TestRenderPostGame(t *testing.T) {
	buffer, screen := newScreenTest()

	t.Run("when postGame is empty", func(t *testing.T) {
		buffer.Reset()
		want := "\x1b[15;0H(s)hare (e)xit"
		screen.renderPostGame()
		assert.Equal(t, want, buffer.String())
	})

	t.Run("when postGame has something to display", func(t *testing.T) {
		buffer.Reset()
		want := "\x1b[15;0H\x1b[JMESSAGE\n\r(s)hare (e)xit"
		screen.postGame = "MESSAGE"
		screen.renderPostGame()
		assert.Equal(t, want, buffer.String())
	})
}

func newScreenTest() (*bytes.Buffer, *screen) {
	b := &bytes.Buffer{}

	return b, newTestScreen(b, wordle.NewGame(wordle.WithCustomWord("CHORE")))
}
