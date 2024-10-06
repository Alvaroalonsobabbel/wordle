package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Alvaroalonsobabbel/wordle/pkg/terminal"
	"golang.org/x/term"
)

const (
	hideCursor = "\033[?25l"
	showCursor = "\033[?25h"
)

var hardMode bool

func main() {
	evalOptions()

	restoreConsole := startRawConsole()
	defer restoreConsole()

	terminal.New(hardMode).Start()
}

func startRawConsole() func() {
	fmt.Print(hideCursor)
	oldState, err := term.MakeRaw(int(os.Stdin.Fd())) //nolint: gosec
	if err != nil {
		log.Fatalf("Error setting terminal to raw mode: %v", err)
	}

	return func() {
		if err := term.Restore(int(os.Stdin.Fd()), oldState); err != nil { //nolint: gosec
			log.Fatalf("unable to retore the terminal original state: %v", err)
		}
		fmt.Print(showCursor)
	}
}

func evalOptions() {
	flag.BoolVar(&hardMode, "hard", false, "Sets the Game to Hard Mode")
	flag.BoolFunc("version", "Prints version", version)
	flag.BoolFunc("rmstatus", "Deletes the status file", removeStatus)
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func version(string) error {
	fmt.Println(terminal.VERSION)
	os.Exit(0)

	return nil
}

func removeStatus(string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %v", err)
	}

	err = os.Remove(filepath.Join(homeDir, terminal.StatusFile))
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error deleting file: %v", err)
	}

	fmt.Println("Status file removed.")
	os.Exit(0)

	return nil
}
