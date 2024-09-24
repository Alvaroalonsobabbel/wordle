package terminal

import (
	"fmt"
	"strings"

	"github.com/Alvaroalonsobabbel/wordle/pkg/wordle"
)

type round struct {
	index     int
	status    []string
	animation string
}

type rounds struct {
	w   *wordle.Status
	all []round
}

func emptyRound() round { //nolint: revive
	return round{status: strings.Split(strings.Repeat("_", 5), "")}
}

func newRounds(w *wordle.Status) *rounds { //nolint: revive
	var r []round
	for range 6 {
		r = append(r, emptyRound())
	}
	return &rounds{
		w:   w,
		all: r,
	}
}

func (r *rounds) string() string {
	var p string
	for i := range 6 {
		if i < len(r.w.Results) {
			var str []string
			for _, b := range r.w.Results[i] {
				for k, v := range b {
					switch v {
					case wordle.Correct:
						str = append(str, fmt.Sprintf(green, string(k)))
					case wordle.Present:
						str = append(str, fmt.Sprintf(yellow, string(k)))
					case wordle.Absent:
						str = append(str, fmt.Sprintf(black, string(k)))
					}
				}
			}
			p += "\t" + strings.Join(str, " ") + newLine
		} else {
			p += "\t" + r.all[i].animation + strings.Join(r.all[i].status, " ") + newLine
		}
	}

	return p
}

func (r *rounds) add(s string) {
	letter := strings.ToUpper(s)

	if r.all[r.w.Round].index == 5 {
		return
	}

	r.all[r.w.Round].status[r.all[r.w.Round].index] = letter
	r.all[r.w.Round].index++
}

func (r *rounds) backspace() {
	if r.all[r.w.Round].index == 0 {
		return
	}

	r.all[r.w.Round].index--
	r.all[r.w.Round].status[r.all[r.w.Round].index] = "_"
}
