package terminal

import (
	"strings"
)

type round struct {
	index     int
	status    []string
	animation string
}

func NewRound() *round { //nolint: revive
	return &round{status: strings.Split(strings.Repeat("_", 5), "")}
}

func NewRounds() []*round { //nolint: revive
	var rounds []*round
	for range 6 {
		rounds = append(rounds, NewRound())
	}

	return rounds
}

func (r *round) string() string {
	return "\t" + r.animation + strings.Join(r.status, " ")
}

func (r *round) add(s string) {
	letter := strings.ToUpper(s)

	if r.index == 5 {
		return
	}

	r.status[r.index] = letter
	r.index++
}

func (r *round) backspace() {
	if r.index == 0 {
		return
	}

	r.index--
	r.status[r.index] = "_"
}
