package terminal

import (
	"testing"

	"github.com/Alvaroalonsobabbel/wordle/pkg/wordle"
	"github.com/stretchr/testify/assert"
)

func TestRoundString(t *testing.T) {
	wordle := wordle.NewGame(wordle.WithCustomWord("CHORE"))
	rounds := newRounds(wordle)

	t.Run("emtpy rounds", func(t *testing.T) {
		want := "\t _  _  _  _  _ "
		assert.Equal(t, want, rounds.string(0))
	})

	t.Run("with one letter", func(t *testing.T) {
		rounds.add("A")
		want := "\t A  _  _  _  _ "
		assert.Equal(t, want, rounds.string(0))
	})

	t.Run("with two letters and animations space", func(t *testing.T) {
		rounds.add("A")
		rounds.all[0].animation = " "
		want := "\t  A  A  _  _  _ "
		assert.Equal(t, want, rounds.string(0))
	})

	t.Run("after a round exist in wordle status it prints the status with color", func(t *testing.T) {
		wordle.Try("SCORE") //nolint: errcheck
		want := "\t\x1b[7m\x1b[90m S \x1b[0m\x1b[7m\x1b[33m C \x1b[0m\x1b[7m\x1b[32m O \x1b[0m\x1b[7m\x1b[32m R \x1b[0m\x1b[7m\x1b[32m E \x1b[0m"
		assert.Equal(t, want, rounds.string(0))
		want = "\t _  _  _  _  _ "
		assert.Equal(t, want, rounds.string(1))
	})
}

func TestAdd(t *testing.T) {
	t.Run("adding one letter", func(t *testing.T) {
		rounds := newRounds(wordle.NewGame(wordle.WithCustomWord("CHORE")))

		rounds.add("A")

		assert.Equal(t, "A", rounds.all[0].status[0])
	})

	t.Run("adding five consecutive letters", func(t *testing.T) {
		rounds := newRounds(wordle.NewGame(wordle.WithCustomWord("CHORE")))
		letters := []string{"A", "B", "C", "D", "E"}

		for _, l := range letters {
			rounds.add(l)
		}

		for i := range letters {
			assert.Equal(t, letters[i], rounds.all[0].status[i])
		}
	})

	t.Run("adding more than 5 letters does not increment the counter nor adds another letter", func(t *testing.T) {
		rounds := newRounds(wordle.NewGame(wordle.WithCustomWord("CHORE")))
		letters := []string{"A", "B", "C", "D", "E", "F"}

		for _, l := range letters {
			rounds.add(l)
		}

		rounds.add("A")

		assert.Equal(t, "E", rounds.all[0].status[4])
		assert.Equal(t, 5, rounds.all[0].index)
	})

	t.Run("adding a lower case letter makes it upper case", func(t *testing.T) {
		rounds := newRounds(wordle.NewGame(wordle.WithCustomWord("CHORE")))
		rounds.add("a")

		assert.Equal(t, "A", rounds.all[0].status[0])
	})
}

func TestBackspace(t *testing.T) {
	t.Run("reverts the counter and replaces the letter with underscore", func(t *testing.T) {
		rounds := newRounds(wordle.NewGame(wordle.WithCustomWord("CHORE")))
		rounds.add("A")
		rounds.add("B")
		rounds.backspace()

		assert.Equal(t, "_", rounds.all[0].status[1])
		assert.Equal(t, 1, rounds.all[0].index)
	})

	t.Run("when counter is 0, backspace has no effect", func(t *testing.T) {
		rounds := newRounds(wordle.NewGame(wordle.WithCustomWord("CHORE")))
		rounds.backspace()

		assert.Equal(t, "_", rounds.all[0].status[0])
		assert.Equal(t, 0, rounds.all[0].index)
	})
}
