package wordle

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTry(t *testing.T) {
	t.Run("correct result assignments", func(t *testing.T) {
		tests := []struct {
			word           string
			inputWord      string
			expectedStatus *Result
		}{
			{
				"CLEAR", "CEDAR", &Result{Correct, Present, Absent, Correct, Correct},
			},
			{
				"CHARM", "BLAST", &Result{Absent, Absent, Correct, Absent, Absent},
			},
			{
				"TIGHT", "FIGHT", &Result{Absent, Correct, Correct, Correct, Correct},
			},
			{
				"CRACK", "OPIUM", &Result{Absent, Absent, Absent, Absent, Absent},
			},
			{
				"CHORE", "ROACH", &Result{Present, Present, Absent, Present, Present},
			},
			{
				// Second to last L should be absent as one L has been provided and
				// it was already been discovered.
				"SPOIL", "QUILL", &Result{Absent, Absent, Present, Absent, Correct},
			},
		}

		for _, test := range tests {
			t.Run(test.word+test.inputWord, func(t *testing.T) {
				wordle := NewTestWordle(false, test.word)
				err := wordle.Try(test.inputWord)
				assert.NoError(t, err)

				for i := range len(test.word) {
					assert.Equal(t, (*test.expectedStatus)[i], wordle.Results[0][i], fmt.Sprintf("error in char %d", i))
				}
			})
		}
	})

	t.Run("consecutive hints", func(t *testing.T) {
		wordle := NewTestWordle(false, "STILL")
		err := wordle.Try("LOVER")
		assert.NoError(t, err)
		expectedResult := Result{Present, Absent, Absent, Absent, Absent}
		assert.Equal(t, expectedResult, wordle.Results[wordle.Round-1])

		wordle.Try("ALLOW") //nolint: errcheck
		expectedResult = Result{Absent, Present, Present, Absent, Absent}
		assert.Equal(t, expectedResult, wordle.Results[wordle.Round-1])

		wordle.Try("LEVEL") //nolint: errcheck
		expectedResult = Result{Present, Absent, Absent, Absent, Correct}
		assert.Equal(t, expectedResult, wordle.Results[wordle.Round-1])
	})

	t.Run("a word not in the list returns an error", func(t *testing.T) {
		badWord := "AAAAA"
		expectedError := fmt.Errorf("Not in word list: %s", badWord)
		wordle := NewGame(false, true)

		err := wordle.Try(badWord)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("hard mode: hints must be used", func(t *testing.T) {
		testWord := "WORLD"
		wordle := &Game{
			hardMode:     true,
			wordle:       testWord,
			allowedWords: allowedWords(),
		}

		err := wordle.Try("DIARY")
		assert.NoError(t, err)

		err = wordle.Try("OPIUM")
		assert.Error(t, err)
		expectedError := fmt.Errorf("Guess must contain D")
		assert.Equal(t, expectedError, err)

		err = wordle.Try("ADOBE")
		assert.Error(t, err)
		expectedError = fmt.Errorf("Guess must contain R")
		assert.Equal(t, expectedError, err)
	})

	t.Run("hard mode: discovered words must be used in the correct place", func(t *testing.T) {
		testWord := "WORLD"
		wordle := &Game{
			hardMode:     true,
			wordle:       testWord,
			allowedWords: allowedWords(),
		}

		err := wordle.Try("WEARY")
		assert.NoError(t, err)

		err = wordle.Try("OPIUM")
		assert.Error(t, err)
		expectedError := fmt.Errorf("1st letter must be W")
		assert.Equal(t, expectedError, err)

		err = wordle.Try("WARTY")
		assert.NoError(t, err)

		err = wordle.Try("WHERE")
		assert.Error(t, err)
		expectedError = fmt.Errorf("3rd letter must be R")
		assert.Equal(t, expectedError, err)
	})
}

func TestFinish(t *testing.T) {
	t.Run("finish returns false while game is running", func(t *testing.T) {
		wordle := NewTestWordle(false, "HELLO")
		err := wordle.Try("WORLD")
		assert.NoError(t, err)

		ok, msg := wordle.Finish()
		assert.False(t, ok)
		assert.Empty(t, msg)
	})

	t.Run("finish returns true if game ends due to win", func(t *testing.T) {
		wordle := NewTestWordle(false, "HELLO")
		err := wordle.Try("WORLD")
		assert.NoError(t, err)

		ok, msg := wordle.Finish()
		assert.False(t, ok)
		assert.Empty(t, msg)

		err = wordle.Try("HELLO")
		assert.NoError(t, err)

		ok, msg = wordle.Finish()
		assert.True(t, ok)
		assert.Empty(t, msg)
	})

	t.Run("wining first try return 'Genius'", func(t *testing.T) {
		wordle := NewTestWordle(false, "HELLO")
		err := wordle.Try("HELLO")
		assert.NoError(t, err)

		ok, msg := wordle.Finish()
		assert.True(t, ok)
		assert.Equal(t, msg, "Genius")
	})

	t.Run("wining last try return 'Phew!'", func(t *testing.T) {
		wordle := NewTestWordle(false, "HELLO")
		for range 5 {
			wordle.Try("WORLD") //nolint: errcheck
		}
		wordle.Try("HELLO") //nolint: errcheck

		ok, msg := wordle.Finish()
		assert.True(t, ok)
		assert.Equal(t, msg, "Phew!")
	})

	t.Run("finish returns true if game ends due to lose and err msg is current wordle", func(t *testing.T) {
		wordle := NewTestWordle(false, "HELLO")

		for range 6 {
			wordle.Try("WORLD") //nolint: errcheck
		}

		ok, msg := wordle.Finish()
		assert.True(t, ok)
		assert.Equal(t, "HELLO", msg)
	})
}

func TestFetchTodaysWordle(t *testing.T) {
	assert.Equal(t, 5, len(fetchTodaysWordle()))
}
