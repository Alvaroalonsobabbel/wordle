package wordle

import "net/http"

type ConfigSetter func(*Game)

func WithHardMode(h bool) ConfigSetter {
	return func(g *Game) {
		g.hardMode = h
	}
}

func WithDalyWordle() ConfigSetter {
	return func(g *Game) {
		g.wordle, g.wordleNumber = fetchTodaysWordle(&http.Client{})
	}
}

func WithCustomWord(word string) ConfigSetter {
	return func(g *Game) {
		g.wordle = word
	}
}
