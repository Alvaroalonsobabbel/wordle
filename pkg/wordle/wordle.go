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
	Absent  = Status(iota) // word not found
	Correct                // word found in the correct place
	Present                // word found in an incorrect place

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

type (
	Status int
	Result []Status

	Game struct {
		Round   int
		Results [6]Result

		conf         *config
		allowedWords []string
		hints        []rune
		discovered   [5]rune
	}
)

func NewGame(conf ...ConfigSetter) *Game {
	config := &config{}

	for _, confSetter := range conf {
		confSetter(config)
	}

	return &Game{
		conf:         config,
		allowedWords: allowedWords(),
	}
}

func (g *Game) Try(word string) error {
	if !slices.Contains(g.allowedWords, word) {
		return fmt.Errorf("Not in word list: %s", word) //nolint: stylecheck
	}

	if g.conf.hardMode {
		if err := g.hardModeCheck(word); err != nil {
			return err
		}
	}
	g.result(word)

	return nil
}

func (g *Game) Finish() (bool, string) {
	if string(g.discovered[:]) == g.conf.wordle {
		return true, finishMessage[g.Round]
	}

	if g.Round > 5 {
		return true, g.conf.wordle
	}

	return false, ""
}

func (g *Game) hardModeCheck(word string) error {
	for i, v := range g.discovered {
		if v != 0 && v != rune(word[i]) {
			return fmt.Errorf("%s letter must be %c", ordinalNumbers[i], v)
		}
	}

	for _, v := range g.hints {
		if !strings.Contains(word, string(v)) {
			return fmt.Errorf("Guess must contain %c", v) //nolint: stylecheck
		}
	}

	return nil
}

func (g *Game) result(word string) {
	hintCounter := maxHints(g.conf.wordle)
	g.Results[g.Round] = Result{Absent, Absent, Absent, Absent, Absent}

	for i, v := range g.conf.wordle {
		if rune(word[i]) == v {
			g.Results[g.Round][i] = Correct
			g.discovered[i] = rune(word[i])
			hintCounter[v]--
		}
	}

	for i := range g.conf.wordle {
		if strings.Contains(g.conf.wordle, string(word[i])) {
			g.hints = append(g.hints, rune(word[i]))
			if hintCounter[rune(word[i])] > 0 && g.Results[g.Round][i] != Correct {
				g.Results[g.Round][i] = Present
				hintCounter[rune(word[i])]--
			}
		}
	}

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
