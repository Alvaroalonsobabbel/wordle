package terminal

import (
	"fmt"
	"io"
	"sync"
	"time"
)

const (
	errDuration = 1500 * time.Millisecond
	errOffset   = 3
)

type err struct {
	msg string
	ts  int64
}

type render struct {
	errCh chan err
	strCh chan string
	dur   time.Duration
	w     io.Writer
	wg    sync.WaitGroup
}

func newRender(w io.Writer) (*render, func()) {
	r := &render{
		errCh: make(chan err),
		strCh: make(chan string),
		dur:   errDuration,
		w:     w,
	}
	go r.errMgr()
	go r.strMgr()

	return r, func() {
		r.wg.Wait()
		close(r.errCh)
		close(r.strCh)
	}
}

func (r *render) strMgr() {
	for s := range r.strCh {
		fmt.Fprint(r.w, s)
		r.wg.Done()
	}
}

func (r *render) errMgr() {
	var q []err

	for log := range r.errCh {
		switch log.ts < time.Now().UnixNano()-r.dur.Nanoseconds() {
		case true:
			q = q[:len(q)-1]
		default:
			q = append([]err{log}, q...)
			go r.wait(log)
		}

		r.printQ(q)
		r.wg.Done()
	}
}

func (r *render) wait(log err) {
	r.wg.Add(1)
	time.Sleep(r.dur)
	r.errCh <- log
}

func (r *render) err(s string) {
	r.wg.Add(1)
	r.errCh <- err{msg: s, ts: time.Now().UnixNano()}
}

func (r *render) string(s string) {
	r.wg.Add(1)
	r.strCh <- s
}

func (r *render) printQ(q []err) {
	for i := range len(q) + 1 {
		fmt.Fprintf(r.w, "\033[%d;28H\033[K", i+errOffset)
	}

	for i, log := range q {
		fmt.Fprintf(r.w, "\033[%d;28H\x1b[3m\x1b[30m\x1b[47m %s \x1b[0m", i+errOffset, log.msg)
	}
}
