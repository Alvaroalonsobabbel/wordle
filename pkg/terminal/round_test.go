package terminal

import (
	"testing"

	"github.com/Alvaroalonsobabbel/wordle/pkg/wordle"
	"github.com/stretchr/testify/assert"
)

func TestRoundsString(t *testing.T) {
	wordle := wordle.NewGame(wordle.WithCustomWord("CHORE"))
	rounds := newRounds(wordle)

	t.Run("emtpy rounds", func(t *testing.T) {
		want := "\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r"
		assert.Equal(t, want, rounds.string())
	})

	t.Run("with one letter", func(t *testing.T) {
		rounds.add("A")
		expected := "\tA _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r"
		assert.Equal(t, expected, rounds.string())
	})

	t.Run("with two letters and animations space", func(t *testing.T) {
		rounds.add("A")
		rounds.all[0].animation = " "
		expected := "\t A A _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r"
		assert.Equal(t, expected, rounds.string())
	})

	t.Run("after a round exist in wordle status it prints the status with color", func(t *testing.T) {
		wordle.Try("SCORE")
		expected := "\t\033[1mS\033[0m \x1b[1m\x1b[33mC\x1b[0m \x1b[1m\x1b[32mO\x1b[0m \x1b[1m\x1b[32mR\x1b[0m \x1b[1m\x1b[32mE\x1b[0m\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r"
		assert.Equal(t, expected, rounds.string())
	})

	t.Run("adding a futher letter prints the solved round plus the new one", func(t *testing.T) {
		rounds.add("A")
		expected := "\t\033[1mS\033[0m \x1b[1m\x1b[33mC\x1b[0m \x1b[1m\x1b[32mO\x1b[0m \x1b[1m\x1b[32mR\x1b[0m \x1b[1m\x1b[32mE\x1b[0m\n\r\tA _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r\t_ _ _ _ _\n\r"
		assert.Equal(t, expected, rounds.string())
	})
}

func TestAdd(t *testing.T) {
	t.Run("adding one letter", func(t *testing.T) {
		wordle := wordle.NewGame(wordle.WithCustomWord("CHORE"))
		rounds := newRounds(wordle)

		rounds.add("A")

		assert.Equal(t, "A", rounds.all[0].status[0])
	})

	t.Run("adding five consecutive letters", func(t *testing.T) {
		wordle := wordle.NewGame(wordle.WithCustomWord("CHORE"))
		rounds := newRounds(wordle)
		letters := []string{"A", "B", "C", "D", "E"}

		for _, l := range letters {
			rounds.add(l)
		}

		for i := range letters {
			assert.Equal(t, letters[i], rounds.all[0].status[i])
		}
	})

	t.Run("adding more than 5 letters does not increment the counter nor adds another letter", func(t *testing.T) {
		wordle := wordle.NewGame(wordle.WithCustomWord("CHORE"))
		rounds := newRounds(wordle)
		letters := []string{"A", "B", "C", "D", "E", "F"}

		for _, l := range letters {
			rounds.add(l)
		}

		rounds.add("A")

		assert.Equal(t, "E", rounds.all[0].status[4])
		assert.Equal(t, 5, rounds.all[0].index)
	})

	t.Run("adding a lower case letter makes it upper case", func(t *testing.T) {
		wordle := wordle.NewGame(wordle.WithCustomWord("CHORE"))
		rounds := newRounds(wordle)
		rounds.add("a")

		assert.Equal(t, "A", rounds.all[0].status[0])
	})
}

func TestBackspace(t *testing.T) {
	t.Run("reverts the counter and replaces the letter with underscore", func(t *testing.T) {
		wordle := wordle.NewGame(wordle.WithCustomWord("CHORE"))
		rounds := newRounds(wordle)
		rounds.add("A")
		rounds.add("B")
		rounds.backspace()

		assert.Equal(t, "_", rounds.all[0].status[1])
		assert.Equal(t, 1, rounds.all[0].index)
	})

	t.Run("when counter is 0, backspace has no effect", func(t *testing.T) {
		wordle := wordle.NewGame(wordle.WithCustomWord("CHORE"))
		rounds := newRounds(wordle)
		rounds.backspace()

		assert.Equal(t, "_", rounds.all[0].status[0])
		assert.Equal(t, 0, rounds.all[0].index)
	})
}
