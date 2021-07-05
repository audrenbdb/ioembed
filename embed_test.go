package ioembed

import (
	"embed"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed test.txt
var testFile embed.FS

func TestGetFiles(t *testing.T) {
	by := make([]byte, 5)
	fileNames := []string{"test.txt"}
	tests := []struct {
		inReader  fileReader
		fileNames []string
		outFiles  map[string][]byte
		outErr    error
	}{
		{
			inReader:  &mockFs{bytes: by, err: errors.New("failed to read")},
			fileNames: fileNames,
			outFiles:  nil,
			outErr:    errors.New("failed to read"),
		},
		{
			inReader:  &mockFs{bytes: by, err: nil},
			fileNames: fileNames,
			outFiles: map[string][]byte{"test.txt":by},
		},
	}
	for _, test := range tests {
		files, err := getFiles(test.inReader, test.fileNames...)
		assert.Equal(t, test.outErr, err)
		assert.Equal(t, test.outFiles, files)
	}
}

type mockFs struct {
	bytes []byte
	err   error
}

func (m *mockFs) ReadFile(fileName string) ([]byte, error) {
	return m.bytes, m.err
}
