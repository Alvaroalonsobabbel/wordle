package wordle

import (
	"log"
	"net/http"
)

type ConfigSetter func(*Status)

func WithHardMode(h bool) ConfigSetter {
	return func(g *Status) {
		g.HardMode = h
	}
}

func WithDalyWordle() ConfigSetter {
	w, n, err := fetchTodaysWordle(http.DefaultClient)
	if err != nil {
		log.Fatal(err)
	}
	return func(g *Status) {
		g.Wordle, g.PuzzleNumber = w, n
	}
}

func WithCustomWord(w string) ConfigSetter {
	return func(g *Status) {
		g.Wordle = w
	}
}
