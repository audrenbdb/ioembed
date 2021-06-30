package ioembed

import (
	"embed"
	"io/ioutil"
	"os"
	"path"
)

//CopyFS copy embed FS to given destination
type CopyFS = func(fs embed.FS, fileName string, dirPath string) error
type getFSBytes = func(fs embed.FS, fileName string) ([]byte, error)
type writeToFile = func(bytes []byte, dirPath string, fileName string) error

//NewCopyFSFunc with default libraries
func NewCopyFSFunc() CopyFS {
	getFSBytes := newGetFSBytesFunc()
	writeToFile := newWriteToFileFunc()
	return newCopyFSFunc(getFSBytes, writeToFile)
}

func newCopyFSFunc(getFSBytes getFSBytes, writeToFile writeToFile) CopyFS {
	return func(fs embed.FS, fileName string, dirPath string) error {
		bytes, err := getFSBytes(fs, fileName)
		if err != nil {
			return err
		}
		return writeToFile(bytes, dirPath, fileName)
	}
}

func newGetFSBytesFunc() getFSBytes {
	return func(fs embed.FS, fileName string) ([]byte, error) {
		return fs.ReadFile(fileName)
	}
}

func newWriteToFileFunc() writeToFile {
	return func(bytes []byte, dirPath string, fileName string) error {
		fp := path.Join(dirPath, fileName)
		return ioutil.WriteFile(fp, bytes, os.ModePerm)
	}
}
