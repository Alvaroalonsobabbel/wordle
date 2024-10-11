package terminal

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/Alvaroalonsobabbel/wordle/pkg/status"
	"github.com/Alvaroalonsobabbel/wordle/pkg/wordle"
	"github.com/atotto/clipboard"
)

const (
	// Relevant unicode characters to control the game.
	enter     = 13
	backspace = 127
	ctrlC     = 3

	// Accepted characters regex.
	okRegex = `^[A-Z]$`
)

type terminal struct {
	wordle *wordle.Status
	screen *screen
	reader io.Reader
}

func New(hardMode bool, localStatus *wordle.Status) *terminal { //nolint: revive
	t := &terminal{reader: os.Stdin}

	wordle := wordle.NewGame(wordle.WithDalyWordle(), wordle.WithHardMode(hardMode))
	if localStatus != nil && localStatus.Wordle == wordle.Wordle {
		wordle = localStatus
	}

	t.wordle = wordle
	t.screen = newScreen(wordle)

	return t
}

func (t *terminal) Start() {
	defer func() {
		fmt.Fprint(t.screen.Writer, newLine)
		if err := status.Game().Save(t.wordle); err != nil {
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
			t.screen.renderMsg(fmt.Sprintf(italics, msg))

			break
		}

		buf, quit := t.read()
		if quit {
			return
		}

		t.processInput(buf[0])
	}

	t.postGame()
}

func (t *terminal) postGame() {
	t.screen.renderPostGame()

	for {
		buf, quit := t.read()
		if quit {
			return
		}

		switch buf[0] {
		case 's', 'S':
			clipboard.WriteAll(t.wordle.Share()) //nolint: errcheck
			t.screen.queueErr("Copied to Clipboard!")
		case 'e', 'E':
			return
		}
	}
}

func (t *terminal) processInput(b byte) {
	t.screen.renderKBFlash(b)

	switch b {
	case backspace:
		t.screen.rounds.backspace()
	case enter:
		if t.screen.rounds.all[t.wordle.Round].index < 5 {
			t.screen.queueErr("Not enough letters")
			t.screen.shakeRound()
			return
		}

		lastWord := strings.Join(t.screen.rounds.all[t.wordle.Round].status, "")
		if err := t.wordle.Try(lastWord); err != nil {
			t.screen.queueErr(err.Error())
			t.screen.shakeRound()
			return
		}

		t.screen.renderResult()
		t.screen.renderKB()
	default:
		if regexp.MustCompile(okRegex).MatchString(strings.ToUpper(string(b))) {
			t.screen.rounds.add(string(b))
		}
	}

	t.screen.renderRound()
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
