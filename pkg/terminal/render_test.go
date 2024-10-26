package terminal

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	t.Run("prints formatted err to w and clears error after", func(t *testing.T) {
		buf := &bytes.Buffer{}
		render := newRender(buf)
		render.errDur = 10 * time.Millisecond

		render.err("123")
		render.wg.Wait()

		want := "\x1b[3;28H\x1b[K\x1b[4;28H\x1b[K\x1b[3;28H\x1b[3m\x1b[30m\x1b[47m 123 \x1b[0m\x1b[3;28H\x1b[K"
		assert.Equal(t, want, buf.String())
	})

	t.Run("prints str to w", func(t *testing.T) {
		buf := &bytes.Buffer{}
		render := newRender(buf)

		render.string("123")
		render.wg.Wait()
		assert.Equal(t, "123", buf.String())
	})
}
