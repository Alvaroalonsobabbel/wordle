package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Alvaroalonsobabbel/wordle/status"
	"github.com/Alvaroalonsobabbel/wordle/terminal"
	"github.com/Alvaroalonsobabbel/wordle/wordle"
	"golang.org/x/term"
)

const VERSION = "v0.4.8"

const (
	hardModeFlag     = "hard"
	versionFlag      = "version"
	removeStatusFlag = "rmstatus"

	hideCursor = "\033[?25l"
	showCursor = "\033[13;0H\n\r\033[?25h"
)

var hardMode bool

func main() {
	evalOptions()

	s, err := status.Game().Load()
	if err != nil {
		log.Fatal(err)
	}

	wordle := wordle.NewGame(wordle.WithDalyWordle(), wordle.WithHardMode(hardMode))
	if s != nil && s.Wordle == wordle.Wordle {
		wordle = s
	}

	restore := startRawConsole()
	defer restore()
	terminal.New(wordle).Start()

	if err := status.Game().Save(wordle); err != nil {
		fmt.Println(err)
	}
}

func evalOptions() {
	flag.BoolVar(&hardMode, hardModeFlag, false, "Sets the Game to Hard Mode")
	flag.BoolFunc(versionFlag, "Prints version", version)
	flag.BoolFunc(removeStatusFlag, "Deletes the status file", status.Remove)
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
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

func version(string) error {
	fmt.Println(VERSION)
	os.Exit(0)

	return nil
}
