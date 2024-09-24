package wordle

import "net/http"

type ConfigSetter func(*Status)

func WithHardMode(h bool) ConfigSetter {
	return func(g *Status) {
		g.HardMode = h
	}
}

func WithDalyWordle() ConfigSetter {
	return func(g *Status) {
		g.Wordle, g.PuzzleNumber = fetchTodaysWordle(&http.Client{})
	}
}

func WithCustomWord(word string) ConfigSetter {
	return func(g *Status) {
		g.Wordle = word
	}
}
