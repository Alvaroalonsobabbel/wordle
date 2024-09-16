package wordle

import "net/http"

type ConfigSetter func(*config)

type config struct {
	wordle       string
	wordleNumber int
	hardMode     bool
}

func WithHardMode(h bool) ConfigSetter {
	return func(c *config) {
		c.hardMode = h
	}
}

func WithDalyWordle() ConfigSetter {
	return func(c *config) {
		c.wordle, c.wordleNumber = fetchTodaysWordle(&http.Client{})
	}
}

func WithCustomWord(word string) ConfigSetter {
	return func(c *config) {
		c.wordle = word
	}
}
