package terminal

import (
	"bytes"
	"io"
	"testing"

	"github.com/Alvaroalonsobabbel/wordle/pkg/wordle"
	"github.com/stretchr/testify/assert"
)

func TestKeyboardString(t *testing.T) {
	w := wordle.NewGame(wordle.WithCustomWord("CHAIR"))
	buf := &bytes.Buffer{}
	r, _ := newRender(buf)
	kb := newKeyboard(w, r)

	t.Run("game with no rounds return a basic keyboard", func(t *testing.T) {
		kb.print()
		r.wg.Wait()

		want := "\x1b[11;2H Q  W  E  R  T  Z  U  I  O  P \x1b[12;3H A  S  D  F  G  H  J  K  L \x1b[13;2H ↩︎  Y  X  C  V  B  N  M  ← "
		assert.Equal(t, want, buf.String())
	})

	t.Run("game with 1 round highlight the used letters", func(t *testing.T) {
		w.Try("HELLO") //nolint: errcheck
		buf.Reset()
		kb.print()
		r.wg.Wait()

		want := "\x1b[11;2H Q  W \x1b[7m\x1b[90m E \x1b[0m R  T  Z  U  I \x1b[7m\x1b[90m O \x1b[0m P \x1b[12;3H A  S  D  F  G \x1b[7m\x1b[33m H \x1b[0m J  K \x1b[7m\x1b[90m L \x1b[0m\x1b[13;2H ↩︎  Y  X  C  V  B  N  M  ← "
		assert.Equal(t, want, buf.String())
	})
}

func TestMapRunes(t *testing.T) {
	var (
		w     = wordle.NewGame(wordle.WithCustomWord("CHAIR"))
		r, _  = newRender(io.Discard)
		kb    = newKeyboard(w, r)
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
