package terminal

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Alvaroalonsobabbel/wordle/pkg/wordle"
)

const (
	title = "\033[1m6 attempts to find a 5-letter word\n\033[0m"

	// Byte relevant characters.
	enter     = 13
	backspace = 127
	ctrlC     = 3
	esc       = 27

	// Accepted characters regex.
	okRegex = `^[A-Z\r\x7F]+$`

	// Terminal formatting.
	newLine          = "\n\r"
	clearScreen      = "\033[H\033[2J"
	green            = "\x1b[1m\x1b[32m%s\x1b[0m"
	yellow           = "\x1b[1m\x1b[33m%s\x1b[0m"
	black            = "\033[1m%s\033[0m"
	greenBackground  = "\x1b[7m\x1b[32m%s\x1b[0m"
	yellowBackground = "\x1b[7m\x1b[33m%s\x1b[0m"
	greyBackground   = "\x1b[7m\x1b[90m%s\x1b[0m"
	flash            = "\x1b[30m\x1b[47m%s\x1b[0m"
	italics          = "\x1b[3m%s\x1b[0m"
)

type renderer struct {
	wordle   *wordle.Game
	keyboard *keyboard
	rounds   []*round

	printer      io.Writer
	regex        *regexp.Regexp
	currentRound int
	errorMsg     string
}

func New(w io.Writer, hardMode, offline bool) *renderer { //nolint: revive
	return &renderer{
		wordle:   wordle.NewGame(hardMode, offline),
		keyboard: NewKB(),
		rounds:   NewRounds(),
		printer:  w,
		regex:    regexp.MustCompile(okRegex),
	}
}

func NewTestTerminal(w io.Writer, word string) *renderer { //nolint: revive
	return &renderer{
		wordle:   wordle.NewTestWordle(false, word),
		keyboard: NewKB(),
		rounds:   NewRounds(),
		printer:  w,
		regex:    regexp.MustCompile(okRegex),
	}
}

func (r *renderer) Start() {
	buf := make([]byte, 1)
	r.render()

	for {
		ok, msg := r.wordle.Finish()
		if ok {
			r.errorMsg = fmt.Sprintf(italics, msg)
			r.render()

			return
		}

		_, err := os.Stdin.Read(buf)
		if err != nil {
			log.Fatalf("Error reading input: %v", err)
		}

		// Ctrl-C or Esc exits the game
		if buf[0] == ctrlC || buf[0] == esc {
			return
		}

		r.enter(buf[0])
	}
}

func (r *renderer) render() {
	fmt.Fprint(r.printer, clearScreen)
	fmt.Fprint(r.printer, title)
	fmt.Fprint(r.printer, newLine)

	for _, v := range r.rounds {
		fmt.Fprint(r.printer, v.string())
		fmt.Fprint(r.printer, newLine)
	}

	fmt.Fprint(r.printer, newLine)
	fmt.Fprint(r.printer, r.errorMsg)
	fmt.Fprint(r.printer, newLine)
	fmt.Fprint(r.printer, r.keyboard.string())
}

func (r *renderer) enter(b byte) {
	if !r.regex.MatchString(strings.ToUpper(string(b))) {
		return
	}

	r.showKBFlash(b)
	switch b {
	case backspace:
		r.rounds[r.currentRound].backspace()
	case enter:
		if r.rounds[r.currentRound].index < 5 {
			r.showError(errors.New("Not enough letters")) //nolint: stylecheck
			return
		}

		lastWord := strings.Join(r.rounds[r.currentRound].status, "")
		result, err := r.wordle.Try(lastWord)
		if err != nil {
			r.showError(err)
			return
		}

		r.showResult(result)
		r.keyboard.update(result, lastWord)
		r.currentRound++
	default:
		r.rounds[r.currentRound].add(string(b))
	}
	r.render()
}

func (r *renderer) showResult(res *wordle.Result) {
	for i, v := range r.rounds[r.currentRound].status {
		color := black
		switch (*res)[i] {
		case wordle.Correct:
			color = green
		case wordle.Present:
			color = yellow
		}

		r.rounds[r.currentRound].status[i] = "_"
		r.render()
		time.Sleep(150 * time.Millisecond)
		r.rounds[r.currentRound].status[i] = fmt.Sprintf(color, v)
	}
}

func (r *renderer) showError(err error) {
	go func() {
		r.errorMsg = fmt.Sprintf(italics, err.Error())
		defer func() {
			r.errorMsg = ""
			r.render()
		}()

		for i := range 6 {
			if i%2 == 0 {
				r.rounds[r.currentRound].animation = " "
			} else {
				r.rounds[r.currentRound].animation = ""
			}
			r.render()
			time.Sleep(50 * time.Millisecond)
		}
		time.Sleep(1500 * time.Millisecond)
	}()
}

func (r *renderer) showKBFlash(l byte) {
	var char string
	switch l {
	case backspace:
		char = "←"
	case enter:
		char = "↩︎"
	default:
		char = strings.ToUpper(string(l))
	}

	f := fmt.Sprintf(flash, char)
	temp := r.keyboard.am[char]
	r.keyboard.am[char] = f
	r.render()
	r.keyboard.am[char] = temp
	time.Sleep(25 * time.Millisecond)
	r.render()
}
