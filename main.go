package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Alvaroalonsobabbel/wordle/terminal"
	"golang.org/x/term"
)

const (
	hideCursor = "\033[?25l"
	showCursor = "\033[?25h"
)

var (
	hardMode bool
	random   bool
)

func main() {
	evalOptions()

	restoreConsole := startRawConsole()
	defer restoreConsole()

	terminal.New(os.Stdout, hardMode, random).Start()
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
	flag.BoolVar(&random, "random", false, "Picks a random word form the word collection")
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
