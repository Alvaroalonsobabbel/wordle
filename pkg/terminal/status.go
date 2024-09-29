package terminal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Alvaroalonsobabbel/wordle/pkg/wordle"
)

const (
	StatusFile = ".wordle"
	read       = os.O_RDONLY
	write      = os.O_CREATE | os.O_RDWR | os.O_TRUNC
)

type fileOpener interface {
	file(int) (io.ReadWriteCloser, error)
}

type status struct {
	wordle *wordle.Status
	open   fileOpener
}

func newStatus(w *wordle.Status) *status {
	return &status{
		wordle: w,
		open:   &opener{},
	}
}

func (s *status) loadGame() error {
	file, err := s.open.file(read)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}
	defer file.Close()

	w := &wordle.Status{}

	if err := json.NewDecoder(file).Decode(w); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return fmt.Errorf("error decoding wordle status into file: %v", err)
	}

	if s.wordle.Wordle == w.Wordle {
		*s.wordle = *w
	}

	return nil
}

func (s *status) saveGame() error {
	file, err := s.open.file(write)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(s.wordle); err != nil {
		return fmt.Errorf("error encoding wordle status into file: %v", err)
	}

	return nil
}

type opener struct{}

func (opener) file(mode int) (io.ReadWriteCloser, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting home directory: %v", err)
	}

	file, err := os.OpenFile(filepath.Join(homeDir, StatusFile), mode, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		return nil, fmt.Errorf("error opening status file: %v", err)
	}

	return file, nil
}
