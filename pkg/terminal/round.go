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

func newRounds(w *wordle.Status) *rounds { //nolint: revive
	var r []round
	for range 6 {
		r = append(r, round{status: strings.Split(strings.Repeat("_", 5), "")})
	}
	return &rounds{
		w:   w,
		all: r,
	}
}

func (r *rounds) string(round int) string {
	p := tab

	if round < len(r.w.Results) {
		var str []string
		for _, b := range r.w.Results[round] {
			for k, v := range b {
				var color string
				switch v {
				case wordle.Correct:
					color = greenBackground
				case wordle.Present:
					color = yellowBackground
				case wordle.Absent:
					color = greyBackground
				}
				str = append(str, fmt.Sprintf(color, string(k)))
			}
		}
		p += strings.Join(str, "")
	} else if r.w.Round < 6 {
		p += r.all[round].animation
		for _, s := range r.all[round].status {
			p += fmt.Sprintf(emptyChar, s)
		}
	}

	return p
}

func (r *rounds) add(s string) {
	if r.all[r.w.Round].index == 5 {
		return
	}

	r.all[r.w.Round].status[r.all[r.w.Round].index] = strings.ToUpper(s)
	r.all[r.w.Round].index++
}

func (r *rounds) backspace() {
	if r.all[r.w.Round].index == 0 {
		return
	}

	r.all[r.w.Round].index--
	r.all[r.w.Round].status[r.all[r.w.Round].index] = "_"
}
