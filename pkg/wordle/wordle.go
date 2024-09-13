package wordle

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"math/rand/v2"
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
)

type (
	Status int
	Result []Status
	Game   struct {
		Round   int
		Results [6]Result

		wordle       string
		hardMode     bool
		allowedWords []string
		hints        []string
		discovered   [5]string
		hintMap      map[string]int
	}
)

func NewGame(hardMode, offline bool) *Game {
	var wordle string
	switch offline {
	case true:
		wordle = pickRandomWord()
	case false:
		wordle = fetchTodaysWordle()
	}

	return NewTestWordle(hardMode, wordle)
}

func NewTestWordle(hardMode bool, wordle string) *Game {
	return &Game{
		wordle:       wordle,
		allowedWords: generateAllowedWordsList(),
		hardMode:     hardMode,
		hintMap:      calculateMaxHints(wordle),
	}
}

func (g *Game) Try(word string) (*Result, error) {
	if !slices.Contains(g.allowedWords, word) {
		return nil, fmt.Errorf("Not in word list: %s", word) //nolint: stylecheck
	}

	if g.hardMode {
		if err := g.hardModeCheck(word); err != nil {
			return nil, err
		}
	}

	return g.result(word), nil
}

func (g *Game) Finish() (bool, string) {
	if strings.Join(g.discovered[:], "") == g.wordle {
		var msg string
		switch g.Round {
		case 1:
			msg = "Genius"
		case 6:
			msg = "Phew!"
		}

		return true, msg
	}

	if g.Round > 5 {
		return true, g.wordle
	}

	return false, ""
}

func (g *Game) hardModeCheck(word string) error {
	for i, v := range g.discovered {
		if v != "" && v != string(word[i]) {
			return fmt.Errorf("%s letter must be %s", ordinalNumbers[i], v)
		}
	}

	for _, v := range g.hints {
		if !strings.Contains(word, v) {
			return fmt.Errorf("Guess must contain %s", v) //nolint: stylecheck
		}
	}

	return nil
}

func (g *Game) result(word string) *Result {
	var hints = make(map[string]int)
	maps.Copy(hints, g.hintMap)
	g.Results[g.Round] = Result{Absent, Absent, Absent, Absent, Absent}

	for i, v := range g.wordle {
		if word[i] == byte(v) {
			g.Results[g.Round][i] = Correct
			g.discovered[i] = string(word[i])
			hints[string(v)]--
		}
	}

	for i := range g.wordle {
		if strings.Contains(g.wordle, string(word[i])) {
			g.hints = append(g.hints, string(word[i]))
			hintCount := hints[string(word[i])]
			if hintCount > 0 && g.Results[g.Round][i] != Correct {
				g.Results[g.Round][i] = Present
				hints[string(word[i])]--
			}
		}
	}
	defer func() { g.Round++ }()

	return &g.Results[g.Round]
}

func pickRandomWord() string {
	answers := strings.Split(strings.ToUpper(answersList), "\n")

	return strings.ToUpper(answers[rand.IntN(len(answers))]) //nolint: gosec
}

func calculateMaxHints(wordle string) map[string]int {
	hintMap := make(map[string]int)
	for _, v := range wordle {
		hintMap[string(v)]++
	}

	return hintMap
}

func fetchTodaysWordle() string {
	var (
		url = wordleBaseURL + time.Now().Format("2006-01-02") + ".json"
		r   = struct {
			Solution string `json:"solution"`
		}{}
	)

	resp, err := http.Get(url) //nolint: gosec
	if err != nil {
		log.Fatalf("unable to fetch today's wordle: %v", err)
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Fatalf("unable to decode today's wordle json response: %v", err)
	}

	return strings.ToUpper(r.Solution)
}

func generateAllowedWordsList() []string {
	var (
		allowed = strings.Split(strings.ToUpper(allowedList), "\n")
		answers = strings.Split(strings.ToUpper(answersList), "\n")
	)

	return slices.Concat(allowed, answers)
}
