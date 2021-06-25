package ioembed

import (
	"embed"
	"io/ioutil"
	"os"
	"path"
)

//CopyFS copy embed FS to given destination
type CopyFS = func(fs embed.FS, fileName string, dirPath string) error
type GetFSBytes = func(fs embed.FS, fileName string) ([]byte, error)
type WriteToFile = func(bytes []byte, dirPath string, fileName string) error

func NewCopyFSFunc(getFSBytes GetFSBytes, writeToFile WriteToFile) CopyFS {
	return func(fs embed.FS, fileName string, dirPath string) error {
		bytes, err := getFSBytes(fs, fileName)
		if err != nil {
			return err
		}
		return writeToFile(bytes, dirPath, fileName)
	}
}

func NewGetFSBytesFunc() GetFSBytes {
	return func(fs embed.FS, fileName string) ([]byte, error) {
		return fs.ReadFile(fileName)
	}
}

func NewWriteToFileFunc() WriteToFile {
	return func(bytes []byte, dirPath string, fileName string) error {
		fp := path.Join(dirPath, fileName)
		return ioutil.WriteFile(fp, bytes, os.ModePerm)
	}
}
