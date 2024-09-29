package wordle

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"
)

const (
	Absent  = iota // word not found
	Correct        // word found in the correct place
	Present        // word found in an incorrect place

	wordleBaseURL = "https://www.nytimes.com/svc/wordle/v2/"
)

// Allowed list: https://gist.github.com/cfreshman/d5fb56316158a1575898bba1eed3b5da
// Answers list: https://gist.github.com/cfreshman/a7b776506c73284511034e63af1017ee
var (
	//go:embed allowed.txt
	allowedList string
	//go:embed answers.txt
	answersList string

	ordinalNumbers = []string{"1st", "2nd", "3rd", "4th", "5th"}
	finishMessage  = map[int]string{
		1: "Genius",
		2: "Magnificent",
		3: "Impressive",
		4: "Splendid",
		5: "Great",
		6: "Phew!",
	}
)

type Status struct {
	Round        int              `json:"round"`
	PuzzleNumber int              `json:"puzzle_number"`
	Wordle       string           `json:"wordle"`
	HardMode     bool             `json:"hard_mode"`
	Results      [][]map[rune]int `json:"results"`
	Discovered   [5]rune          `json:"discovered"`
	Hints        []rune           `json:"hints"`

	allowedWords []string
}

func NewGame(conf ...ConfigSetter) *Status {
	game := &Status{}

	for _, confSetter := range conf {
		confSetter(game)
	}

	return game
}

func (s *Status) Try(word string) error {
	if err := s.isAllowed(word); err != nil {
		return err
	}

	if s.HardMode {
		if err := s.hardModeCheck(word); err != nil {
			return err
		}
	}
	s.result(word)

	return nil
}

func (s *Status) Finish() (bool, string) {
	if string(s.Discovered[:]) == s.Wordle {
		return true, finishMessage[s.Round]
	}

	if s.Round > 5 {
		return true, s.Wordle
	}

	return false, ""
}

func (s *Status) hardModeCheck(word string) error {
	for i, v := range s.Discovered {
		if v != 0 && v != rune(word[i]) {
			return fmt.Errorf("%s letter must be %c", ordinalNumbers[i], v)
		}
	}

	for _, v := range s.Hints {
		if !strings.Contains(word, string(v)) {
			return fmt.Errorf("Guess must contain %c", v) //nolint: stylecheck
		}
	}

	return nil
}

func (s *Status) result(word string) {
	var (
		hintCounter = maxHints(s.Wordle)
		currentWord []map[rune]int
	)

	for _, v := range word {
		currentWord = append(currentWord, map[rune]int{v: Absent})
	}

	for i, v := range word {
		if v == rune(s.Wordle[i]) {
			currentWord[i][v] = Correct
			s.Discovered[i] = v
			hintCounter[v]--
		}
	}

	for i, v := range word {
		if strings.Contains(s.Wordle, string(v)) {
			s.Hints = append(s.Hints, v)
			if hintCounter[v] > 0 && currentWord[i][v] != Correct {
				currentWord[i][v] = Present
				hintCounter[v]--
			}
		}
	}

	s.Results = append(s.Results, currentWord)
	s.Round++
}

func (s *Status) isAllowed(word string) error {
	if s.allowedWords == nil {
		s.allowedWords = slices.Concat(
			strings.Split(strings.ToUpper(allowedList), "\n"),
			strings.Split(strings.ToUpper(answersList), "\n"),
		)
	}

	if !slices.Contains(s.allowedWords, word) {
		return fmt.Errorf("Not in word list: %s", word) //nolint: stylecheck
	}

	return nil
}

func maxHints(wordle string) map[rune]int {
	hintMap := make(map[rune]int)
	for _, v := range wordle {
		hintMap[v]++
	}

	return hintMap
}

func fetchTodaysWordle(c *http.Client) (string, int) {
	var (
		url = wordleBaseURL + time.Now().Format("2006-01-02") + ".json"
		r   = struct {
			Solution string `json:"solution"`
			Number   int    `json:"days_since_launch"`
		}{}
	)

	resp, err := c.Get(url) //nolint: gosec
	if err != nil {
		log.Fatalf("unable to fetch today's wordle: %v", err)
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Fatalf("unable to decode today's wordle json response: %v", err)
	}

	return strings.ToUpper(r.Solution), r.Number
}
