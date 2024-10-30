package terminal

import (
	"bytes"
	"io"
	"testing"

	"github.com/Alvaroalonsobabbel/wordle/wordle"
	"github.com/stretchr/testify/assert"
)

func TestRoundPrint(t *testing.T) {
	wordle := wordle.NewGame(wordle.WithCustomWord("CHORE"))
	buf := &bytes.Buffer{}
	render := newRender(buf)
	round := newRound(wordle, render)

	t.Run("emtpy rounds", func(t *testing.T) {
		round.print(0)
		render.wg.Wait()
		want := "\x1b[3;9H _  _  _  _  _  "
		assert.Equal(t, want, buf.String())
	})

	t.Run("with one letter", func(t *testing.T) {
		buf.Reset()
		round.add("A")
		render.wg.Wait()
		want := "\x1b[3;9H A  _  _  _  _  "
		assert.Equal(t, want, buf.String())
	})

	t.Run("with two letters and animations space", func(t *testing.T) {
		buf.Reset()
		round.add("A")
		render.wg.Wait()
		want := "\x1b[3;9H A  A  _  _  _  "
		assert.Equal(t, want, buf.String())
	})

	t.Run("after a round exist in wordle status it prints the status with color", func(t *testing.T) {
		buf.Reset()
		wordle.Try("SCORE") //nolint: errcheck
		round.print(0)
		render.wg.Wait()
		want := "\x1b[3;9H\x1b[7m\x1b[90m S \x1b[0m\x1b[7m\x1b[33m C \x1b[0m\x1b[7m\x1b[32m O \x1b[0m\x1b[7m\x1b[32m R \x1b[0m\x1b[7m\x1b[32m E \x1b[0m "
		assert.Equal(t, want, buf.String())

		buf.Reset()
		round.reset()
		round.print(1)
		render.wg.Wait()
		want = "\x1b[4;9H _  _  _  _  _  "
		assert.Equal(t, want, buf.String())
	})
}

func TestAdd(t *testing.T) {
	t.Run("adding one letter", func(t *testing.T) {
		wordle := wordle.NewGame(wordle.WithCustomWord("CHORE"))
		render := newRender(io.Discard)
		round := newRound(wordle, render)

		round.add("A")
		assert.Equal(t, "A", round.status[0])
	})

	t.Run("adding five consecutive letters", func(t *testing.T) {
		wordle := wordle.NewGame(wordle.WithCustomWord("CHORE"))
		render := newRender(io.Discard)
		round := newRound(wordle, render)
		letters := []string{"A", "B", "C", "D", "E"}

		for _, l := range letters {
			round.add(l)
		}

		for i := range letters {
			assert.Equal(t, letters[i], round.status[i])
		}
	})

	t.Run("adding more than 5 letters does not increment the counter nor adds another letter", func(t *testing.T) {
		wordle := wordle.NewGame(wordle.WithCustomWord("CHORE"))
		render := newRender(io.Discard)
		round := newRound(wordle, render)
		letters := []string{"A", "B", "C", "D", "E", "F"}

		for _, l := range letters {
			round.add(l)
		}

		round.add("A")

		assert.Equal(t, "E", round.status[4])
		assert.Equal(t, 5, round.index)
	})
}

func TestBackspace(t *testing.T) {
	t.Run("reverts the counter and replaces the letter with underscore", func(t *testing.T) {
		w := wordle.NewGame(wordle.WithCustomWord("CHORE"))
		r := newRender(io.Discard)
		round := newRound(w, r)
		round.add("A")
		round.add("B")
		round.backspace()

		assert.Equal(t, "_", round.status[1])
		assert.Equal(t, 1, round.index)
	})

	t.Run("when counter is 0, backspace has no effect", func(t *testing.T) {
		w := wordle.NewGame(wordle.WithCustomWord("CHORE"))
		r := newRender(io.Discard)
		round := newRound(w, r)
		round.backspace()

		assert.Equal(t, "_", round.status[0])
		assert.Equal(t, 0, round.index)
	})
}
