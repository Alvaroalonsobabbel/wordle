package terminal

import (
	"fmt"
	"io"
	"time"
)

const (
	errTimeoff = 1500 * time.Millisecond
	errOffset  = 3
)

type err struct {
	message string
	ts      int64
	rm      bool
}

type errorQ struct {
	queue   []err
	errCh   chan err
	timeOff time.Duration
	w       io.Writer
}

func newErrorQueue(w io.Writer) *errorQ {
	e := &errorQ{
		errCh:   make(chan err),
		timeOff: errTimeoff,
		w:       w,
	}
	go e.queueMgr()

	return e
}

func (e *errorQ) queueMgr() {
	for log := range e.errCh {
		switch log.rm {
		case true:
			for i, l := range e.queue {
				if log.ts == l.ts {
					e.queue = append(e.queue[:i], e.queue[i+1:]...)
					continue
				}
			}
		default:
			e.queue = append([]err{log}, e.queue...)
			go e.timeErr(log)
		}

		e.displayErr()
	}
}

func (e *errorQ) timeErr(log err) {
	time.Sleep(e.timeOff)
	log.rm = true
	e.errCh <- log
}

func (e *errorQ) queueErr(message string) {
	e.errCh <- err{message: message, ts: time.Now().UnixNano()}
}

func (e *errorQ) displayErr() {
	for i := range len(e.queue) + 1 {
		fmt.Fprintf(e.w, "\033[%d;22H\033[K", i+errOffset)
	}

	for i, log := range e.queue {
		fmt.Fprintf(e.w, "\033[%d;22H\x1b[3m\x1b[30m\x1b[47m %s \x1b[0m", i+errOffset, log.message)
	}
}
