package terminal

import (
	"io"
	"testing"

	"github.com/Alvaroalonsobabbel/wordle/pkg/wordle"
	"github.com/stretchr/testify/assert"
)

type mockFile struct {
	data []byte
}

func (m *mockFile) Read(p []byte) (n int, err error) {
	if len(m.data) == 0 {
		return 0, io.EOF
	}
	n = copy(p, m.data)
	m.data = m.data[n:]
	return n, nil
}

func (m *mockFile) Write(p []byte) (n int, err error) {
	m.data = append(m.data, p...)

	return len(p), nil
}

func (m *mockFile) Close() error {
	return nil
}

type mockOpener struct {
	f *mockFile
}

func (m *mockOpener) file(int) (io.ReadWriteCloser, error) {
	return m.f, nil
}

func TestLoadGame(t *testing.T) {
	tests := []struct {
		name   string
		word   string
		data   string
		wantPN int
		wantHM bool
	}{
		{
			name:   "when daily wordle matches stored status, stored status gets loaded",
			word:   "BRAIN",
			data:   `{"round":4,"puzzle_number":1197,"wordle":"BRAIN","hard_mode":true,"results":[],"discovered":[],"hints":[]}`,
			wantPN: 1197,
			wantHM: true,
		},
		{
			name:   "when daily wordle does not matches stored status, new wordle gets loaded",
			word:   "CHAIR",
			data:   `{"round":4,"puzzle_number":1197,"wordle":"BRAIN","hard_mode":true,"results":[],"discovered":[],"hints":[]}`,
			wantPN: 0,
			wantHM: false,
		},
		{
			name:   "when status file is empty, new wordle gets loaded",
			word:   "CHAIR",
			data:   ``,
			wantPN: 0,
			wantHM: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockFile := &mockFile{data: []byte(test.data)}
			wordle := wordle.NewGame(wordle.WithCustomWord(test.word))
			status := &status{wordle: wordle, open: &mockOpener{f: mockFile}}

			err := status.loadGame()
			assert.NoError(t, err)
			assert.Equal(t, test.word, status.wordle.Wordle)
			assert.Equal(t, test.wantPN, status.wordle.PuzzleNumber)
			assert.Equal(t, test.wantHM, status.wordle.HardMode)
		})
	}
}

func TestSaveGame(t *testing.T) {
	mockFile := &mockFile{}
	wordle := wordle.NewGame(wordle.WithCustomWord("CHAIR"), wordle.WithHardMode(true))
	status := &status{wordle: wordle, open: &mockOpener{f: mockFile}}

	err := status.saveGame()
	assert.NoError(t, err)
	want := `{"round":0,"puzzle_number":0,"wordle":"CHAIR","hard_mode":true,"results":null,"discovered":[0,0,0,0,0],"hints":null}
`
	assert.Equal(t, want, string(mockFile.data))
}
