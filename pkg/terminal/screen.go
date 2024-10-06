package terminal

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Alvaroalonsobabbel/wordle/pkg/wordle"
)

const (
	// Static strings.
	title        = "\033[1m6 attempts to find a 5-letter word\n\033[0m"
	postGameMenu = "(s)hare (e)xit"

	// Terminal formatting.
	tab              = "\t"
	newLine          = "\n\r"
	emptyChar        = " %s "
	greenBackground  = "\x1b[7m\x1b[32m %s \x1b[0m"
	yellowBackground = "\x1b[7m\x1b[33m %s \x1b[0m"
	greyBackground   = "\x1b[7m\x1b[90m %s \x1b[0m"
	flash            = "\x1b[30m\x1b[47m %s \x1b[0m"
	italics          = "\x1b[3m%s\x1b[0m"
	moveToY          = "\033[%d;0H%s"
	moveToYX         = "\033[%d;%dH%s"
	clearRow         = "\033[K"
	clearRowDown     = "\033[J"
	clearScreen      = "\033[H\033[2J"

	// Screen position.
	roundOffset = 3
	msgPos      = 10
	kbPos       = 11
	postGamePos = 15
)

type screen struct {
	*wordle.Status
	*rounds
	*keyboard
	*errorQ
	io.Writer
	msg      string
	postGame string
}

func newScreen(wordle *wordle.Status) *screen {
	return newTestScreen(os.Stdout, wordle)
}

func newTestScreen(w io.Writer, wordle *wordle.Status) *screen {
	return &screen{
		wordle,
		newRounds(wordle),
		NewKB(wordle),
		newErrorQueue(w),
		w,
		"",
		"",
	}
}

func (s *screen) renderRound() {
	fmt.Fprintf(s.Writer, moveToY, s.Round+roundOffset, s.rounds.string(s.Round))
}

func (s *screen) renderKB() {
	fmt.Fprintf(s.Writer, moveToY, kbPos, s.keyboard.string())
}

func (s *screen) renderMsg() {
	fmt.Fprintf(s.Writer, moveToY, msgPos, s.msg)
}

func (s *screen) renderPostGame() {
	pg := postGameMenu
	if s.postGame != "" {
		pg = clearRowDown + s.postGame + newLine + postGameMenu
	}

	fmt.Fprintf(s.Writer, moveToY, postGamePos, pg)
}

func (s *screen) renderKBFlash(l byte) {
	var char string
	switch l {
	case backspace:
		char = "←"
	case enter:
		char = "↩︎"
	default:
		char = strings.ToUpper(string(l))
	}
	s.keyboard.flash = char
	s.renderKB()
	time.Sleep(25 * time.Millisecond)
	s.keyboard.flash = ""
	s.renderKB()
}

func (s *screen) renderErr(err error) {
	s.queueErr(err.Error())

	go func() {
		for i := range 6 {
			if i%2 == 0 {
				s.rounds.all[s.Round].animation = " "
			} else {
				s.rounds.all[s.Round].animation = ""
			}

			s.rewriteRow(s.Round+roundOffset, "")
			s.renderRound()
			time.Sleep(50 * time.Millisecond)
		}
	}()
}

func (s *screen) renderResult() {
	// this happens after wordle.Try increments the round counter
	// that's why we run it with wordle.Round - 1
	var (
		r                    = s.Round - 1
		lineOffset           = 9
		lineOffsetMultiplier = 3
	)

	for i, res := range s.Results[r] {
		for k, v := range res {
			color := greyBackground
			switch v {
			case wordle.Correct:
				color = greenBackground
			case wordle.Present:
				color = yellowBackground
			}

			fmt.Fprintf(s.Writer, moveToYX, r+roundOffset, (i*lineOffsetMultiplier)+lineOffset, " _ ")
			time.Sleep(250 * time.Millisecond)
			fmt.Fprintf(s.Writer, moveToYX, r+roundOffset, (i*lineOffsetMultiplier)+lineOffset, fmt.Sprintf(color, string(k)))
		}
	}
}

func (s *screen) renderAll() {
	var r []string
	for i := range 6 {
		r = append(r, s.rounds.string(i))
	}

	fmt.Fprint(s.Writer,
		clearScreen,
		title,
		newLine,
		strings.Join(r, newLine),
		newLine,
		newLine,
		s.msg,
		newLine,
		s.keyboard.string(),
		newLine,
		s.postGame,
		newLine,
	)
}

func (s *screen) rewriteRow(row int, content string) {
	emptyRow := strings.Repeat(" ", 27) + "\r"
	fmt.Fprintf(s.Writer, moveToY, row, emptyRow+content)
}
