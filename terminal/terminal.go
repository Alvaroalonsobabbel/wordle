package terminal

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/Alvaroalonsobabbel/wordle/status"
	"github.com/Alvaroalonsobabbel/wordle/wordle"
	"github.com/atotto/clipboard"
	"golang.org/x/term"
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
	hideCursor       = "\033[?25l"
	showCursor       = "\033[13;0H\n\r\033[?25h"
	emptyChar        = " %s "
)

var finishMessage = []string{"Genius", "Magnificent", "Impressive", "Splendid", "Great", "Phew!"}

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
	restoreConsole := startRawConsole()

	defer func() {
		t.render.close()
		restoreConsole()
		if err := status.Game().Save(t.wordle); err != nil {
			fmt.Println(err)
		}
	}()

	t.initialScreen()
	t.game()
}

func (t *terminal) game() {
	for !t.wordle.Finish() {
		buf, quit := t.read()
		if quit {
			return
		}

		t.processInput(buf[0])
	}

	t.render.string(fmt.Sprintf(italicFooter, t.finishingMsg()))
	t.postGame()
}

func (t *terminal) postGame() {
	t.render.string(postGameMenu)

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
		if t.round.index < 5 {
			t.render.err("Not enough letters")
			t.round.shake()
			return
		}

		lastWord := strings.Join(t.round.status, "")
		if err := t.wordle.Try(lastWord); err != nil {
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
	message := t.wordle.Wordle
	if string(t.wordle.Discovered[:]) == t.wordle.Wordle {
		message = finishMessage[t.wordle.Round-1]
	}
	return message
}

func startRawConsole() func() {
	fmt.Print(hideCursor)
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Error setting terminal to raw mode: %v", err)
	}

	return func() {
		if err := term.Restore(int(os.Stdin.Fd()), oldState); err != nil {
			log.Fatalf("unable to retore the terminal original state: %v", err)
		}
		fmt.Print(showCursor)
	}
}
