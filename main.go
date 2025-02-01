package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Alvaroalonsobabbel/wordle/status"
	"github.com/Alvaroalonsobabbel/wordle/terminal"
	"github.com/Alvaroalonsobabbel/wordle/wordle"
)

const VERSION = "v0.4.9"

const (
	hardModeFlag     = "hard"
	versionFlag      = "version"
	removeStatusFlag = "rmstatus"
)

var hardMode bool

func main() {
	evalOptions()

	status, err := status.Game().Load()
	if err != nil {
		log.Fatal(err)
	}

	terminal.New(wordle.NewGame(hardMode, wordle.WithSavedWordle(status))).Start()
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
