package ioembed

import (
	"embed"
)


func GetFiles(fs embed.FS, fileNames ...string) (map[string][]byte, error) {
	return getFiles(fs, fileNames...)
}

func getFiles(reader fileReader, fileNames ...string) (map[string][]byte, error) {
	files := map[string][]byte{}
	for _, f := range fileNames {
		b, err := reader.ReadFile(f)
		if err != nil {
			return nil, err
		}
		files[f] = b
	}
	return files, nil
}

type fileReader interface {
	ReadFile(fileName string) ([]byte, error)
}
