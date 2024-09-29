package terminal

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/Alvaroalonsobabbel/wordle/pkg/wordle"
)

const VERSION = "v0.2.0"

const (
	// Relevant unicode characters to control the game.
	enter     = 13
	backspace = 127
	ctrlC     = 3
	esc       = 27

	// Accepted characters regex.
	okRegex = `^[A-Z\r\x7F]+$`
)

type terminal struct {
	wordle *wordle.Status
	screen *screen
	regex  *regexp.Regexp
	status *status
	reader io.Reader

	buf []byte
}

func New(hardMode bool) *terminal { //nolint: revive
	wordle := wordle.NewGame(wordle.WithDalyWordle(), wordle.WithHardMode(hardMode))
	return &terminal{
		reader: os.Stdin,
		wordle: wordle,
		screen: newScreen(wordle),
		status: newStatus(wordle),
		regex:  regexp.MustCompile(okRegex),
		buf:    make([]byte, 1),
	}
}

func NewTestTerminal(w io.Writer, r io.Reader) *terminal { //nolint: revive
	wordle := wordle.NewGame(wordle.WithCustomWord("CHORE"))
	return &terminal{
		reader: r,
		wordle: wordle,
		screen: newTestScreen(w, wordle),
		regex:  regexp.MustCompile(okRegex),
		buf:    make([]byte, 1),
	}
}

func (t *terminal) Start() {
	if err := t.status.loadGame(); err != nil {
		fmt.Println(err)

		return
	}

	defer func() {
		t.screen.renderAll()
		if err := t.status.saveGame(); err != nil {
			fmt.Println(err)
		}
	}()

	t.game()
}

func (t *terminal) game() {
	t.screen.renderAll()

	for {
		ok, msg := t.wordle.Finish()
		if ok {
			t.screen.msg = fmt.Sprintf(italics, msg)
			t.screen.renderMsg()
			t.postGame()

			break
		}

		if quit := t.read(); quit {
			break
		}

		t.processInput(t.buf[0])
	}
}

func (t *terminal) postGame() {
	t.screen.renderPostGame()
	for {
		if quit := t.read(); quit {
			return
		}

		switch t.buf[0] {
		case 's', 'S':
			t.screen.postGame = t.wordle.Share()
			t.screen.renderPostGame()
		case 'e', 'E':
			return
		}
	}
}

func (t *terminal) processInput(b byte) {
	if !t.regex.MatchString(strings.ToUpper(string(b))) {
		return
	}

	t.screen.renderKBFlash(b)
	switch b {
	case backspace:
		t.screen.rounds.backspace()
	case enter:
		if t.screen.rounds.all[t.wordle.Round].index < 5 {
			t.screen.renderErr(errors.New("Not enough letters")) //nolint: stylecheck
			return
		}

		lastWord := strings.Join(t.screen.rounds.all[t.wordle.Round].status, "")
		if err := t.wordle.Try(lastWord); err != nil {
			t.screen.renderErr(err)
			return
		}

		t.screen.renderResult()
		t.screen.renderKB()
	default:
		t.screen.rounds.add(string(b))
	}

	t.screen.renderRound()
}

func (t *terminal) read() bool {
	_, err := t.reader.Read(t.buf)
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	// Ctrl-C or Esc exits the game
	if t.buf[0] == ctrlC || t.buf[0] == esc {
		return true
	}

	return false
}
