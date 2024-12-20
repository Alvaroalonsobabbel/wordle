package terminal

import (
	"io"
	"strings"
	"testing"

	"github.com/Alvaroalonsobabbel/wordle/wordle"
	"github.com/stretchr/testify/assert"
)

func TestPrintKey(t *testing.T) {
	tests := []struct {
		initialWord string
		tries       []string
		want        map[string]string
	}{
		{
			initialWord: "ENDOW",
			tries:       []string{"STING", "KNEEL"},
			want: map[string]string{
				"E": "\x1b[11;8H\x1b[7m\x1b[33m E \x1b[0m",
				"T": "\x1b[11;14H\x1b[7m\x1b[90m T \x1b[0m",
				"I": "\x1b[11;23H\x1b[7m\x1b[90m I \x1b[0m",
				"S": "\x1b[12;6H\x1b[7m\x1b[90m S \x1b[0m",
				"G": "\x1b[12;15H\x1b[7m\x1b[90m G \x1b[0m",
				"K": "\x1b[12;24H\x1b[7m\x1b[90m K \x1b[0m",
				"L": "\x1b[12;27H\x1b[7m\x1b[90m L \x1b[0m",
				"N": "\x1b[13;20H\x1b[7m\x1b[32m N \x1b[0m",
				"Z": "\x1b[11;17H Z ",
			},
		},
		{
			initialWord: "HEFTY",
			tries:       []string{"REACT", "DETOX", "TENET"},
			want: map[string]string{
				"E": "\x1b[11;8H\x1b[7m\x1b[32m E \x1b[0m",
				"R": "\x1b[11;11H\x1b[7m\x1b[90m R \x1b[0m",
				"T": "\x1b[11;14H\x1b[7m\x1b[33m T \x1b[0m",
				"O": "\x1b[11;26H\x1b[7m\x1b[90m O \x1b[0m",
				"A": "\x1b[12;3H\x1b[7m\x1b[90m A \x1b[0m",
				"D": "\x1b[12;9H\x1b[7m\x1b[90m D \x1b[0m",
				"X": "\x1b[13;8H\x1b[7m\x1b[90m X \x1b[0m",
				"C": "\x1b[13;11H\x1b[7m\x1b[90m C \x1b[0m",
				"N": "\x1b[13;20H\x1b[7m\x1b[90m N \x1b[0m",
				"M": "\x1b[13;23H M ",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.initialWord+": "+strings.Join(test.tries, " "), func(t *testing.T) {
			w := wordle.NewGame(wordle.WithCustomWord(test.initialWord))
			kb := newKeyboard(w, newRender(io.Discard))

			for _, word := range test.tries {
				assert.NoError(t, w.Try(word))
			}

			kb.print()

			for char, want := range test.want {
				assert.Equal(t, want, kb.keys[char].string())
			}
		})
	}
}
