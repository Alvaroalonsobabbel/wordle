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

type render struct {
	errQ   []string
	errCh  chan string
	strCh  chan string
	errDur time.Duration
	w      io.Writer
	wg     sync.WaitGroup
}

func newRender(w io.Writer) *render {
	r := &render{
		errCh:  make(chan string),
		strCh:  make(chan string),
		errDur: errDuration,
		w:      w,
	}
	go r.errMgr()
	go r.strMgr()

	return r
}

func (r *render) err(s string) {
	r.wg.Add(1)
	r.errCh <- s
}

func (r *render) string(s string) {
	r.wg.Add(1)
	r.strCh <- s
}

func (r *render) strMgr() {
	for s := range r.strCh {
		fmt.Fprint(r.w, s)
		r.wg.Done()
	}
}

func (r *render) errMgr() {
	for log := range r.errCh {
		r.errQ = append([]string{log}, r.errQ...)
		time.AfterFunc(r.errDur, r.rmLastErr)
		r.printErrQ()
	}
}

func (r *render) rmLastErr() {
	defer r.wg.Done()

	r.errQ = r.errQ[:len(r.errQ)-1]
	r.printErrQ()
}

func (r *render) printErrQ() {
	for i := range 6 {
		fmt.Fprintf(r.w, "\033[%d;28H\033[K", i+errOffset)
	}

	for i, log := range r.errQ {
		if i < 6 {
			fmt.Fprintf(r.w, "\033[%d;28H\x1b[3m\x1b[30m\x1b[47m %s \x1b[0m", i+errOffset, log)
		}
	}
}

func (r *render) close() {
	r.wg.Wait()
	close(r.errCh)
	close(r.strCh)
}
