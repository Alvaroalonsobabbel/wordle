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
	game.allowedWords = allowedWords()

	for _, confSetter := range conf {
		confSetter(game)
	}

	return game
}

func (g *Status) Try(word string) error {
	if !slices.Contains(g.allowedWords, word) {
		return fmt.Errorf("Not in word list: %s", word) //nolint: stylecheck
	}

	if g.HardMode {
		if err := g.hardModeCheck(word); err != nil {
			return err
		}
	}
	g.result(word)

	return nil
}

func (g *Status) Finish() (bool, string) {
	if string(g.Discovered[:]) == g.Wordle {
		return true, finishMessage[g.Round]
	}

	if g.Round > 5 {
		return true, g.Wordle
	}

	return false, ""
}

func (g *Status) hardModeCheck(word string) error {
	for i, v := range g.Discovered {
		if v != 0 && v != rune(word[i]) {
			return fmt.Errorf("%s letter must be %c", ordinalNumbers[i], v)
		}
	}

	for _, v := range g.Hints {
		if !strings.Contains(word, string(v)) {
			return fmt.Errorf("Guess must contain %c", v) //nolint: stylecheck
		}
	}

	return nil
}

func (g *Status) result(word string) {
	var (
		hintCounter = maxHints(g.Wordle)
		currentWord []map[rune]int
	)

	for _, v := range word {
		currentWord = append(currentWord, map[rune]int{v: Absent})
	}

	for i, v := range word {
		if v == rune(g.Wordle[i]) {
			currentWord[i][v] = Correct
			g.Discovered[i] = v
			hintCounter[v]--
		}
	}

	for i, v := range word {
		if strings.Contains(g.Wordle, string(v)) {
			g.Hints = append(g.Hints, v)
			if hintCounter[v] > 0 && currentWord[i][v] != Correct {
				currentWord[i][v] = Present
				hintCounter[v]--
			}
		}
	}

	g.Results = append(g.Results, currentWord)
	g.Round++
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

func allowedWords() []string {
	return slices.Concat(
		strings.Split(strings.ToUpper(allowedList), "\n"),
		strings.Split(strings.ToUpper(answersList), "\n"),
	)
}
