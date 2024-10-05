package terminal

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestErrorQueue(t *testing.T) {
	t.Run("errors are deleted from the queue after the timer expires", func(t *testing.T) {
		errQ := newErrorQueue(io.Discard)
		errQ.expTime = 10 * time.Millisecond

		errQ.queueErr("123")
		assert.Equal(t, 1, len(errQ.queue))
		assert.Equal(t, "123", errQ.queue[0].message)
		time.Sleep(15 * time.Millisecond)
		assert.Equal(t, 0, len(errQ.queue))
	})

	t.Run("displayErr prints the error queue after every update", func(t *testing.T) {
		var buf bytes.Buffer
		errQ := newErrorQueue(&buf)
		errQ.expTime = 10 * time.Millisecond

		errQ.queueErr("123")
		want := "\x1b[3;22H\x1b[K\x1b[4;22H\x1b[K\x1b[3;22H\x1b[3m\x1b[30m\x1b[47m 123 \x1b[0m"
		assert.Equal(t, want, buf.String())
		buf.Reset()
		time.Sleep(15 * time.Millisecond)
		want = "\x1b[3;22H\x1b[K"
		assert.Equal(t, want, buf.String())
	})
}
