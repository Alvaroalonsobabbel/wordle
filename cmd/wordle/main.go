package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Alvaroalonsobabbel/wordle/pkg/status"
	"github.com/Alvaroalonsobabbel/wordle/pkg/terminal"
	"golang.org/x/term"
)

const VERSION = "v0.4.4"

const (
	hideCursor = "\033[?25l"
	showCursor = "\033[13;0H\n\r\033[?25h"

	// Flags.
	hardModeFlag     = "hard"
	versionFlag      = "version"
	removeStatusFlag = "rmstatus"
)

var hardMode bool

func main() {
	evalOptions()

	status, err := status.Game().Load()
	if err != nil {
		fmt.Println(err)

		return
	}

	restoreConsole := startRawConsole()
	defer restoreConsole()

	terminal, closer := terminal.New(hardMode, status)
	defer closer()

	terminal.Start()
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

func evalOptions() {
	flag.BoolVar(&hardMode, hardModeFlag, false, "Sets the Game to Hard Mode")
	flag.BoolFunc(versionFlag, "Prints version", version)
	flag.BoolFunc(removeStatusFlag, "Deletes the status file", status.Remove)
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func version(string) error {
	fmt.Println(VERSION)
	os.Exit(0)

	return nil
}
