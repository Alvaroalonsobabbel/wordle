package terminal

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	t.Run("errors are queued and removed from the queue after errDur", func(t *testing.T) {
		render := newRender(io.Discard)
		render.errDur = 10 * time.Millisecond

		for range 5 {
			render.err("123")
		}

		time.Sleep(5 + time.Millisecond)
		assert.Equal(t, 5, len(render.errQ))

		render.wg.Wait()
		assert.Equal(t, 0, len(render.errQ))
	})

	t.Run("prints formatted err to w and clears error after", func(t *testing.T) {
		buf := &bytes.Buffer{}
		render := newRender(buf)
		render.errDur = 10 * time.Millisecond

		render.err("123")
		render.wg.Wait()
		assert.Equal(t, "\x1b[3;28H\x1b[K\x1b[4;28H\x1b[K\x1b[5;28H\x1b[K\x1b[6;28H\x1b[K\x1b[7;28H\x1b[K\x1b[8;28H\x1b[K\x1b[3;28H\x1b[3m\x1b[30m\x1b[47m 123 \x1b[0m\x1b[3;28H\x1b[K\x1b[4;28H\x1b[K\x1b[5;28H\x1b[K\x1b[6;28H\x1b[K\x1b[7;28H\x1b[K\x1b[8;28H\x1b[K", buf.String())
	})

	t.Run("prints str to w", func(t *testing.T) {
		buf := &bytes.Buffer{}
		render := newRender(buf)

		render.string("123")
		render.wg.Wait()
		assert.Equal(t, "123", buf.String())
	})
}
