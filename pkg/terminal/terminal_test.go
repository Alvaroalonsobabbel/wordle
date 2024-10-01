package terminal

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockReader struct {
	data []byte
}

func (m *mockReader) Read(p []byte) (int, error) {
	n := copy(p, m.data)
	m.data = m.data[n:]
	return len(p), nil
}

func TestRead(t *testing.T) {
	mockReader := &mockReader{}
	terminal := NewTestTerminal(io.Discard, mockReader)

	tests := []struct {
		name     string
		input    []byte
		wantExit bool
	}{
		{"ctrl-c exits the game", []byte{ctrlC}, true},
		{"'A' does not exits the game", []byte{'A'}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockReader.data = test.input
			exit := terminal.read()
			if test.wantExit {
				assert.True(t, exit)
			} else {
				assert.False(t, exit)
			}
		})
	}
}
