package wordle

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"maps"
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

type Status int

type Result []Status

type Game struct {
	wordle      string // unexport
	hardMode    bool
	guessesList []string
	hints       []string
	discovered  [5]string
	roundCount  int
	hintMap     map[string]int
}

func NewGame(hardMode bool) *Game {
	game := &Game{
		wordle:      fetchTodaysWordle(),
		guessesList: generateWordsList(),
		hardMode:    hardMode,
		hintMap:     make(map[string]int),
	}
	game.calculateMaxHints()

	return game
}

func NewTestWordle(hardMode bool, wordle string) *Game {
	game := &Game{
		wordle:      wordle,
		guessesList: generateWordsList(),
		hardMode:    hardMode,
		hintMap:     make(map[string]int),
	}
	game.calculateMaxHints()

	return game
}

func (g *Game) Try(word string) (Result, error) {
	if !slices.Contains(g.guessesList, word) {
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
		switch g.roundCount {
		case 1:
			msg = "Genius"
		case 6:
			msg = "Phew!"
		}

		return true, msg
	}

	if g.roundCount > 5 {
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

func (g *Game) result(word string) Result {
	var (
		res   = Result{Absent, Absent, Absent, Absent, Absent}
		hints = make(map[string]int)
	)
	maps.Copy(hints, g.hintMap)

	for i, v := range g.wordle {
		if word[i] == byte(v) {
			res[i] = Correct
			g.discovered[i] = string(word[i])
			hints[string(v)]--
		}
	}

	for i := range g.wordle {
		if strings.Contains(g.wordle, string(word[i])) {
			g.hints = append(g.hints, string(word[i]))
			hintCount := hints[string(word[i])]
			if hintCount > 0 && res[i] != Correct {
				res[i] = Present
				hints[string(word[i])]--
			}
		}
	}
	g.roundCount++

	return res
}

// func (g *Game) pickRandomWord() {
//
// 	g.Wordle = strings.ToUpper(g.wordList[rand.IntN(len(g.wordList))])
// }

func (g *Game) calculateMaxHints() {
	for _, v := range g.wordle {
		g.hintMap[string(v)]++
	}
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
		log.Fatalf("unable to get toda's wordle: %v", err)
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Fatalf("unable to decode today's wordle json response: %v", err)
	}

	return strings.ToUpper(r.Solution)
}

func generateWordsList() []string {
	var (
		allowed   = strings.Split(strings.ToUpper(allowedList), "\n")
		responses = strings.Split(strings.ToUpper(answersList), "\n")
	)

	return slices.Concat(allowed, responses)
}
