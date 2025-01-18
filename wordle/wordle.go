package wordle

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"time"
)

const (
	Absent  = iota // word not found
	Correct        // word found in the correct place
	Present        // word found in an incorrect place

	wordleBaseURL = "https://www.nytimes.com/svc/wordle/v2/%s.json"
)

// Allowed list: https://gist.github.com/cfreshman/d5fb56316158a1575898bba1eed3b5da
// Answers list: https://gist.github.com/cfreshman/a7b776506c73284511034e63af1017ee
var (
	//go:embed allowed.txt
	allowedList string
	//go:embed answers.txt
	answersList string

	ordinalNumbers = []string{"1st", "2nd", "3rd", "4th", "5th"}
)

type Status struct {
	Round        int              `json:"round"`
	PuzzleNumber int              `json:"puzzle_number"`
	Wordle       string           `json:"wordle"`
	HardMode     bool             `json:"hard_mode"`
	Results      [][]map[rune]int `json:"results"`
	Discovered   [5]rune          `json:"discovered"`
	Hints        []rune           `json:"hints"`
	Used         []rune           `json:"used"`

	allowedWords []string
}

func NewGame(conf ...ConfigSetter) *Status {
	s := &Status{}
	for _, confSetter := range conf {
		confSetter(s)
	}

	return s
}

func (s *Status) Try(word string) error {
	if err := s.checkWord(word); err != nil {
		return err
	}
	if err := s.isAllowed(word); err != nil {
		return err
	}
	if err := s.hardModeCheck(word); err != nil {
		return err
	}
	s.result(word)

	return nil
}

func (s *Status) Finish() bool {
	return string(s.Discovered[:]) == s.Wordle || s.Round > 5
}

func (s *Status) checkWord(word string) error {
	if len(word) < 5 {
		return fmt.Errorf("Not enough letters") //nolint: stylecheck
	}
	if len(word) > 5 {
		return fmt.Errorf("Too many letters") //nolint: stylecheck
	}
	if !regexp.MustCompile(`^[A-Z]*$`).MatchString(word) {
		return fmt.Errorf("Only letters are allowed") //nolint: stylecheck
	}
	return nil
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

func (s *Status) hardModeCheck(word string) error {
	if !s.HardMode {
		return nil
	}
	for i, v := range s.Discovered {
		if v != 0 && v != rune(word[i]) {
			return fmt.Errorf("%s letter must be %c", ordinalNumbers[i], v)
		}
	}
	for _, v := range s.Hints {
		if !strings.ContainsRune(word, v) {
			return fmt.Errorf("Guess must contain %c", v) //nolint: stylecheck
		}
	}

	return nil
}

func (s *Status) result(word string) {
	var (
		currentWord []map[rune]int
		hintCounter = make(map[rune]int)
	)

	for _, v := range s.Wordle {
		hintCounter[v]++
	}

	for i, v := range word {
		currentWord = append(currentWord, map[rune]int{v: Absent})
		s.Used = append(s.Used, v)

		if v == rune(s.Wordle[i]) {
			currentWord[i][v] = Correct
			s.Discovered[i] = v
			hintCounter[v]--
		}
	}

	for i, v := range word {
		if strings.ContainsRune(s.Wordle, v) {
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

func fetchTodaysWordle(c *http.Client) (string, int, error) {
	resp, err := c.Get(fmt.Sprintf(wordleBaseURL, time.Now().Format(time.DateOnly)))
	if err != nil {
		return "", 0, fmt.Errorf("unable to fetch today's wordle: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("NYT API returned a non-200 status: %v", resp.StatusCode)
	}
	r := struct {
		Solution string `json:"solution"`
		Number   int    `json:"days_since_launch"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", 0, fmt.Errorf("unable to decode today's wordle json response: %v", err)
	}
	return strings.ToUpper(r.Solution), r.Number, nil
}
