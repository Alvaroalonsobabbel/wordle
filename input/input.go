package input

import (
	"bufio"
	"fmt"
	"io"
)

const ctrlC = 3

type Screener interface {
	SendRune(rune)
	Continue() bool
}

type Input struct {
	reader io.Reader
	screen Screener
}

func New(r io.Reader, s Screener) *Input {
	return &Input{r, s}
}

func (i *Input) Start() {
	buf := bufio.NewReader(i.reader)

	for i.screen.Continue() {
		r, _, err := buf.ReadRune()
		if err != nil {
			if err == io.EOF {
				continue
			}
			fmt.Printf("unable to read rune %v", err)
			return
		}
		if r == ctrlC {
			return
		}

		i.screen.SendRune(r)
		buf.Reset(i.reader)
	}
}
