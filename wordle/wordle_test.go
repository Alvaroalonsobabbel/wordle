package wordle

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTry(t *testing.T) {
	t.Run("correct result assignments", func(t *testing.T) {
		tests := []struct {
			word           string
			inputWord      string
			expectedResult []map[rune]int
		}{
			{
				"CLEAR", "CEDAR", []map[rune]int{{'C': Correct}, {'E': Present}, {'D': Absent}, {'A': Correct}, {'R': Correct}},
			},
			{
				"CHARM", "BLAST", []map[rune]int{{'B': Absent}, {'L': Absent}, {'A': Correct}, {'S': Absent}, {'T': Absent}},
			},
			{
				"TIGHT", "FIGHT", []map[rune]int{{'F': Absent}, {'I': Correct}, {'G': Correct}, {'H': Correct}, {'T': Correct}},
			},
			{
				"CRACK", "OPIUM", []map[rune]int{{'O': Absent}, {'P': Absent}, {'I': Absent}, {'U': Absent}, {'M': Absent}},
			},
			{
				"CHORE", "ROACH", []map[rune]int{{'R': Present}, {'O': Present}, {'A': Absent}, {'C': Present}, {'H': Present}},
			},
			{
				// Second to last L should be absent as the L has already been discovered.
				"SPOIL", "QUILL", []map[rune]int{{'Q': Absent}, {'U': Absent}, {'I': Present}, {'L': Absent}, {'L': Correct}},
			},
		}

		for _, test := range tests {
			t.Run(test.word+test.inputWord, func(t *testing.T) {
				wordle := &Status{Wordle: test.word}

				assert.NoError(t, wordle.Try(test.inputWord))
				assert.Equal(t, test.expectedResult, wordle.Results[0])
			})
		}
	})

	t.Run("consecutive hints", func(t *testing.T) {
		wordle := &Status{Wordle: "STILL"}
		assert.NoError(t, wordle.Try("LOVER"))

		expectedResult := []map[rune]int{{'L': Present}, {'O': Absent}, {'V': Absent}, {'E': Absent}, {'R': Absent}}
		assert.Equal(t, expectedResult, wordle.Results[wordle.Round-1])

		assert.NoError(t, wordle.Try("ALLOW"))
		expectedResult = []map[rune]int{{'A': Absent}, {'L': Present}, {'L': Present}, {'O': Absent}, {'W': Absent}}
		assert.Equal(t, expectedResult, wordle.Results[wordle.Round-1])

		assert.NoError(t, wordle.Try("LEVEL"))
		expectedResult = []map[rune]int{{'L': Present}, {'E': Absent}, {'V': Absent}, {'E': Absent}, {'L': Correct}}
		assert.Equal(t, expectedResult, wordle.Results[wordle.Round-1])
	})

	t.Run("a word not in the list returns an error", func(t *testing.T) {
		wordle := &Status{Wordle: "WHELP"}
		badWord := "AAAAA"
		err := wordle.Try(badWord)
		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("Not in word list: %s", badWord), err)
	})

	t.Run("hard mode: hints must be used", func(t *testing.T) {
		wordle := &Status{Wordle: "WORLD", HardMode: true}
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
		wordle := &Status{Wordle: "WORLD", HardMode: true}
		assert.NoError(t, wordle.Try("WEARY"))

		err := wordle.Try("OPIUM")
		assert.Error(t, err)
		expectedError := fmt.Errorf("1st letter must be W")
		assert.Equal(t, expectedError, err)

		assert.NoError(t, wordle.Try("WARTY"))

		err = wordle.Try("WHERE")
		assert.Error(t, err)
		expectedError = fmt.Errorf("3rd letter must be R")
		assert.Equal(t, expectedError, err)
	})
}

func TestFinish(t *testing.T) {
	t.Run("finish returns false while game is running", func(t *testing.T) {
		wordle := &Status{Wordle: "HELLO"}
		assert.NoError(t, wordle.Try("WORLD"))
		assert.False(t, wordle.Finish())
	})

	t.Run("finish returns true if game ends due to win", func(t *testing.T) {
		wordle := &Status{Wordle: "HELLO"}
		assert.NoError(t, wordle.Try("WORLD"))
		assert.False(t, wordle.Finish())
		assert.NoError(t, wordle.Try("HELLO"))
		assert.True(t, wordle.Finish())
	})

	t.Run("finish returns true if game ends due to lose", func(t *testing.T) {
		wordle := &Status{Wordle: "HELLO"}
		for range 6 {
			assert.NoError(t, wordle.Try("WORLD"))
		}
		assert.True(t, wordle.Finish())
	})
}

func TestIsAllowed(t *testing.T) {
	wordle := &Status{Wordle: "HELLO"}
	assert.Nil(t, wordle.allowedWords)

	assert.NoError(t, wordle.isAllowed("CHORE"))
	assert.Error(t, wordle.isAllowed("AAAAA"))
	wordle.allowedWords = nil
	assert.NoError(t, wordle.isAllowed("CHORE"))
}

func TestNewGame(t *testing.T) {
	tests := []struct {
		name       string
		hardMode   bool
		settings   ConfigSetter
		wantWordle *Status
	}{
		{
			name:       "with no config settings",
			wantWordle: &Status{Wordle: "HELLO", PuzzleNumber: 123},
		},
		{
			name:       "WithCustomWord",
			settings:   WithCustomWord("WORLD"),
			wantWordle: &Status{Wordle: "WORLD"},
		},
		{
			name:       "WithCustomWord and hard mode",
			hardMode:   true,
			settings:   WithCustomWord("WORLD"),
			wantWordle: &Status{Wordle: "WORLD", HardMode: true},
		},
		{
			name:       "WithSavedWordle with today's game returns saved wordle",
			settings:   WithSavedWordle(&Status{Wordle: "HELLO", PuzzleNumber: 123, HardMode: true}),
			wantWordle: &Status{Wordle: "HELLO", PuzzleNumber: 123, HardMode: true},
		},
		{
			name:       "WithSavedWordle with yesterday's game returns today's game",
			settings:   WithSavedWordle(&Status{Wordle: "WORLD", PuzzleNumber: 122, HardMode: true}),
			wantWordle: &Status{Wordle: "HELLO", PuzzleNumber: 123},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := &http.Client{Transport: &mockNYTAPI{resp: &http.Response{
				Body:       io.NopCloser(strings.NewReader(`{"solution": "hello", "days_since_launch": 123}`)),
				StatusCode: http.StatusOK,
			}}}

			var got *Status
			if tt.settings == nil {
				got = newCustomClientGame(tt.hardMode, client)
			} else {
				got = newCustomClientGame(tt.hardMode, client, tt.settings)
			}
			assert.Equal(t, tt.wantWordle, got)
		})
	}
}

type mockNYTAPI struct {
	resp *http.Response
}

func (m *mockNYTAPI) RoundTrip(*http.Request) (*http.Response, error) {
	return m.resp, nil
}

func TestFetchDailyWordle(t *testing.T) {
	tests := []struct {
		name     string
		mockResp *http.Response
		wantErr  bool
	}{
		{
			name: "happy path",
			mockResp: &http.Response{
				Body:       io.NopCloser(strings.NewReader(`{"solution": "hello", "days_since_launch": 123}`)),
				StatusCode: http.StatusOK,
			},
		},
		{
			name:     "NYT API not responding with 200",
			mockResp: &http.Response{StatusCode: http.StatusNotFound},
			wantErr:  true,
		},
		{
			name: "unable to decode json body",
			mockResp: &http.Response{
				Body:       io.NopCloser(strings.NewReader(`"solution": "hello", "days_since_launch": 123}`)),
				StatusCode: http.StatusOK,
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := &http.Client{Transport: &mockNYTAPI{resp: test.mockResp}}
			word, day, err := fetchTodaysWordle(client)

			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, "HELLO", word)
			assert.Equal(t, 123, day)
		})
	}
}
