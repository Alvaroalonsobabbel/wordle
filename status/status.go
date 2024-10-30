package status

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Alvaroalonsobabbel/wordle/wordle"
)

const (
	statusFile = ".wordle"
	read       = os.O_RDONLY
	write      = os.O_CREATE | os.O_RDWR | os.O_TRUNC
)

type fileOpener interface {
	file(int) (io.ReadWriteCloser, error)
}

type status struct {
	open fileOpener
}

func Game() *status { //nolint: revive
	return &status{
		open: &opener{},
	}
}

func (s *status) Load() (*wordle.Status, error) {
	file, err := s.open.file(read)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}
	defer file.Close()

	status := &wordle.Status{}

	if err := json.NewDecoder(file).Decode(status); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil
		}
		return nil, fmt.Errorf("error decoding wordle status into file: %v", err)
	}

	return status, nil
}

func (s *status) Save(status *wordle.Status) error {
	file, err := s.open.file(write)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(status); err != nil {
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

	file, err := os.OpenFile(filepath.Join(homeDir, statusFile), mode, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		return nil, fmt.Errorf("error opening status file: %v", err)
	}

	return file, nil
}

func Remove(string) error {
	// Remove takes a string and it's not used since it
	// has to match the fn signature for flag.BoolFunc
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %v", err)
	}

	err = os.Remove(filepath.Join(homeDir, statusFile))
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error deleting file: %v", err)
	}

	fmt.Println("Status file removed.")
	os.Exit(0)

	return nil
}
