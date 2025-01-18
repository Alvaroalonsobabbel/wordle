package terminal

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/Alvaroalonsobabbel/wordle/wordle"
	"github.com/atotto/clipboard"
)

const (
	// Relevant unicode characters to control the game.
	enter     = 13
	backspace = 127
	ctrlC     = 3

	// Accepted characters regex.
	okRegex = `^[A-Z]$`

	title            = "\033[H\033[2J\033[1m6 attempts to find a 5-letter word\n\033[0m"
	postGameMenu     = "\033[15;0H(s)hare (e)xit"
	italicFooter     = "\033[10;0H\x1b[3m%s\x1b[0m"
	greenBackground  = "\x1b[7m\x1b[32m %s \x1b[0m"
	yellowBackground = "\x1b[7m\x1b[33m %s \x1b[0m"
	greyBackground   = "\x1b[7m\x1b[90m %s \x1b[0m"
	emptyChar        = " %s "
)

type terminal struct {
	wordle   *wordle.Status
	keyboard *keyboard
	round    *round
	render   *render
	reader   io.Reader
}

func New(w *wordle.Status) *terminal { //nolint: revive
	r := newRender(os.Stdout)

	return &terminal{
		reader:   os.Stdin,
		wordle:   w,
		render:   r,
		round:    newRound(w, r),
		keyboard: newKeyboard(w, r),
	}
}

func (t *terminal) Start() {
	defer t.render.close()
	t.initialScreen()
	// game loop
	for !t.wordle.Finish() {
		buf, quit := t.read()
		if quit {
			return
		}

		t.processInput(buf[0])
	}

	t.render.string(fmt.Sprintf(italicFooter, t.finishingMsg()))
	t.render.string(postGameMenu)

	// post game loop
	for {
		buf, quit := t.read()
		if quit {
			return
		}

		switch buf[0] {
		case 's', 'S':
			clipboard.WriteAll(t.wordle.Share()) //nolint: errcheck
			t.render.err("Copied to Clipboard!")
		case 'e', 'E':
			return
		}
	}
}

func (t *terminal) processInput(b byte) {
	t.keyboard.flash(b)

	switch b {
	case backspace:
		t.round.backspace()
	case enter:
		if err := t.wordle.Try(t.round.word()); err != nil {
			t.render.err(err.Error())
			t.round.shake()
			return
		}

		t.round.renderResult()
		t.keyboard.print()
	default:
		c := strings.ToUpper(string(b))
		if regexp.MustCompile(okRegex).MatchString(c) {
			t.round.add(c)
		}
	}
}

func (t *terminal) read() ([]byte, bool) {
	buf := make([]byte, 1)
	if _, err := t.reader.Read(buf); err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	// Ctrl-C exits the game
	if buf[0] == ctrlC {
		return nil, true
	}

	return buf, false
}

func (t *terminal) initialScreen() {
	t.render.string(title)
	t.keyboard.print()

	for i := range 6 {
		t.round.print(i)
	}
}

func (t *terminal) finishingMsg() string {
	var (
		message       = t.wordle.Wordle
		finishMessage = map[int]string{
			1: "Genius",
			2: "Magnificent",
			3: "Impressive",
			4: "Splendid",
			5: "Great",
			6: "Phew!",
		}
	)
	if string(t.wordle.Discovered[:]) == t.wordle.Wordle {
		message = finishMessage[t.wordle.Round]
	}
	return message
}
