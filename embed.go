package ioembed

import (
	"archive/zip"
	"bytes"
	"embed"
)


func GetZippedFiles(fs embed.FS, zipName string) (map[string][]byte, error) {
	zf, err := fs.Open(zipName)
	if err != nil {
		return nil, err
	}
	defer zf.Close()

	r, err := zip.OpenReader(zipName)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return getEachZippedFile(r)
}

func getEachZippedFile(r *zip.ReadCloser) (map[string][]byte, error) {
	files := map[string][]byte{}
	for _, f := range r.File {
		file, err := f.Open()
		if err != nil {
			return nil, err
		}
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(file)
		if err != nil {
			return nil, err
		}
		files[f.Name] = buf.Bytes()
		file.Close()
	}
	return files, nil
}

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
