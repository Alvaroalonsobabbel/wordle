package terminal

import (
	"fmt"
	"strings"
	"time"

	"github.com/Alvaroalonsobabbel/wordle/wordle"
)

const (
	moveToYX    = "\033[%d;%dH%s"
	roundPos    = "\033[%d;9H"
	roundOffset = 3
)

type round struct {
	index     int
	status    [5]string
	animation string
	wordle    *wordle.Status
	render    *render
}

func newRound(w *wordle.Status, r *render) *round { //nolint: revive
	return &round{
		render: r,
		wordle: w,
	}
}

func (r *round) word() string {
	return strings.Join(r.status[:], "")
}

func (r *round) print(round int) {
	p := fmt.Sprintf(roundPos, round+roundOffset)

	switch round < len(r.wordle.Results) {
	case true:
		for _, b := range r.wordle.Results[round] {
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
				p += fmt.Sprintf(color, string(k))
			}
		}
	default:
		p += r.animation
		for _, s := range r.status {
			if s == "" {
				s = "_"
			}
			p += fmt.Sprintf(emptyChar, s)
		}
	}

	r.render.string(p + " ") // Trailing space to delete the error animation tail.
}

func (r *round) shake() {
	go func() {
		for i := range 6 {
			if i%2 == 0 {
				r.animation = " "
			} else {
				r.animation = ""
			}

			r.print(r.wordle.Round)
			time.Sleep(50 * time.Millisecond)
		}
	}()
}

func (r *round) result() {
	defer func() {
		r.index = 0
		r.status = [5]string{}
	}()

	var (
		// results are displyed after wordle.Try increments the internal
		// round counter that's why we run it with wordle.Round - 1
		round = r.wordle.Round - 1
		row   = round + roundOffset
	)

	for i, res := range r.wordle.Results[round] {
		for k, v := range res {
			color := greyBackground
			switch v {
			case wordle.Correct:
				color = greenBackground
			case wordle.Present:
				color = yellowBackground
			}

			// 9 is the distance between the marging and the round.
			// 3 is the spaces each letter occupies.
			r.render.string(fmt.Sprintf(moveToYX, row, 9+(i*3), " _ "))
			time.Sleep(250 * time.Millisecond)
			r.render.string(fmt.Sprintf(moveToYX, row, 9+(i*3), fmt.Sprintf(color, string(k))))
		}
	}
}

func (r *round) add(s string) {
	defer r.print(r.wordle.Round)

	if r.index == 5 {
		return
	}

	r.status[r.index] = s
	r.index++
}

func (r *round) backspace() {
	defer r.print(r.wordle.Round)

	if r.index == 0 {
		return
	}

	r.index--
	r.status[r.index] = "_"
}
