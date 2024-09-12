package terminal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	round := NewRound()

	t.Run("empty round", func(t *testing.T) {
		expected := "\t_ _ _ _ _"
		assert.Equal(t, expected, round.string())
	})

	t.Run("with one letter", func(t *testing.T) {
		round.add("A")
		expected := "\tA _ _ _ _"
		assert.Equal(t, expected, round.string())
	})

	t.Run("with animation space", func(t *testing.T) {
		round.add("A")
		round.animation = " "
		expected := "\t A A _ _ _"
		assert.Equal(t, expected, round.string())
	})
}

func TestAdd(t *testing.T) {
	t.Run("adding one letter", func(t *testing.T) {
		round := NewRound()
		round.add("A")

		assert.Equal(t, round.status[0], "A")
	})

	t.Run("adding five consecutive letters", func(t *testing.T) {
		round := NewRound()
		letters := []string{"A", "B", "C", "D", "E"}

		for _, l := range letters {
			round.add(l)
		}

		for i := range letters {
			assert.Equal(t, round.status[i], letters[i])
		}
	})

	t.Run("adding more than 5 letters does not increment the counter nor adds another letter", func(t *testing.T) {
		round := NewRound()
		letters := []string{"A", "B", "C", "D", "E", "F"}

		for _, l := range letters {
			round.add(l)
		}

		round.add("A")

		assert.Equal(t, round.status[4], "E")
		assert.Equal(t, 5, round.index)
	})

	t.Run("adding a lower case letter makes it upper case", func(t *testing.T) {
		round := NewRound()
		round.add("a")

		assert.Equal(t, "A", round.status[0])
	})
}

func TestBackspace(t *testing.T) {
	t.Run("reverts the counter and replaces the letter with underscore", func(t *testing.T) {
		round := NewRound()
		round.add("A")
		round.add("B")
		round.backspace()

		assert.Equal(t, "_", round.status[1])
		assert.Equal(t, 1, round.index)
	})

	t.Run("when counter is 0, backspace has no effect", func(t *testing.T) {
		round := NewRound()
		round.backspace()

		assert.Equal(t, "_", round.status[0])
		assert.Equal(t, 0, round.index)
	})
}
