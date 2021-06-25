package ioembed_test

import (
	"embed"
	"errors"
	"github.com/audrenbdb/ioembed"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed test.txt
var testFile embed.FS


func TestCopyFS_Integration(t *testing.T) {
	//integration tests done 25/06/2021 16:26 UTC
	t.Skip()
	getFsBytes := ioembed.NewGetFSBytesFunc()
	writeToFile := ioembed.NewWriteToFileFunc()
	copyFS := ioembed.NewCopyFSFunc(getFsBytes, writeToFile)

	err := copyFS(testFile, "test.txt", "output")
	assert.NoError(t, err)
	assert.FileExists(t, "output/test.txt")
}

func TestCopyFS(t *testing.T) {
	tests := []struct{
		description string

		getFSBytes ioembed.GetFSBytes
		writeToFile ioembed.WriteToFile

		fs embed.FS
		fileName string
		dirPath string

		outErr error
	}{
		{
			description: "Extracting bytes from embed FS should fail",

			getFSBytes: func(fs embed.FS, fileName string) ([]byte, error) {
				return nil, errors.New("bytes from file error")
			},
			outErr: errors.New("bytes from file error"),
		},
		{
			description: "Copying file should fail",
			getFSBytes: func(fs embed.FS, fileName string) ([]byte, error) {
				return nil, nil
			},
			writeToFile: func(bytes []byte, dirPath string, fileName string) error {
				return errors.New("copying file error")
			},
			outErr: errors.New("copying file error"),
		},
	}

	for _, test := range tests {
		copyFS := ioembed.NewCopyFSFunc(test.getFSBytes, test.writeToFile)
		err := copyFS(test.fs, test.fileName, test.dirPath)
		assert.Equal(t, test.outErr, err)
	}
}