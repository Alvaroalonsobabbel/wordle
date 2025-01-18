package input_test

import (
	"bytes"
	"testing"

	"github.com/Alvaroalonsobabbel/wordle/input"
	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantRune rune
	}{
		{"io.Reader pass a rune to SendRune()", "a", 'a'},
		{"when multiple runes are sent, SendRune() receives only the first", "hello", 'h'},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockScreen := &spyScreen{cont: true}
			input.New(bytes.NewBufferString(test.input), mockScreen).Start()

			assert.Equal(t, test.wantRune, mockScreen.gotRune)
		})
	}
}

type spyScreen struct {
	gotRune rune
	cont    bool
}

func (m *spyScreen) SendRune(r rune) {
	m.gotRune = r
	m.cont = false
}

func (m *spyScreen) Continue() bool {
	return m.cont
}
