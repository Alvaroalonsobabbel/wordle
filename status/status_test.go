package status

import (
	"io"
	"testing"

	"github.com/Alvaroalonsobabbel/wordle/wordle"
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
	t.Run("when status file has content, a new wordle.Status struct is returned", func(t *testing.T) {
		mockFile := &mockFile{data: []byte(`{"round":4,"puzzle_number":1197,"wordle":"BRAIN","hard_mode":true,"results":[],"discovered":[],"hints":[]}`)}
		status := &status{open: &mockOpener{f: mockFile}}

		wordle, err := status.Load()
		assert.NoError(t, err)
		assert.Equal(t, 1197, wordle.PuzzleNumber)
		assert.True(t, wordle.HardMode)
	})

	t.Run("when status file is empty, nil wordle.Status is returned", func(t *testing.T) {
		mockFile := &mockFile{data: []byte(``)}
		status := &status{open: &mockOpener{f: mockFile}}
		wordle, err := status.Load()
		assert.NoError(t, err)
		assert.Nil(t, wordle)
	})
}

func TestSaveGame(t *testing.T) {
	mockFile := &mockFile{}
	wordle := wordle.NewGame(wordle.WithCustomWord("CHAIR"), wordle.WithHardMode(true))
	status := &status{open: &mockOpener{f: mockFile}}

	err := status.Save(wordle)
	assert.NoError(t, err)
	want := `{"round":0,"puzzle_number":0,"wordle":"CHAIR","hard_mode":true,"results":null,"discovered":[0,0,0,0,0],"hints":null,"used":null}
`
	assert.Equal(t, want, string(mockFile.data))
}
