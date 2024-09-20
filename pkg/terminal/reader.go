package terminal

import (
	"log"
	"os"
)

type reader struct {
	buf []byte
}

func newReader() *reader {
	return &reader{buf: make([]byte, 1)}
}

func (r *reader) read() {
	_, err := os.Stdin.Read(r.buf)
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	// Ctrl-C or Esc exits the game
	if r.buf[0] == ctrlC || r.buf[0] == esc {
		os.Exit(0)
	}
}
