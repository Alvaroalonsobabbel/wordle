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

const VERSION = "v0.4.7"

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

	wordle := wordle.NewGame(wordle.WithDalyWordle(), wordle.WithHardMode(hardMode))
	if status != nil && status.Wordle == wordle.Wordle {
		wordle = status
	}

	terminal.New(wordle).Start()
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
