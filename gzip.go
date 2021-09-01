package ioembed

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"embed"
	"io"
)

func GetGZippedTarFiles(fs embed.FS, name string) (map[string][]byte, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r, err := newGzipTarReader(f)
	if err != nil {
		return nil, err
	}
	files := map[string][]byte{}
	err = getEachTarFile(files, r)
	return files, err
}

func newGzipTarReader(f io.Reader) (*tar.Reader, error) {
	r, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	return tar.NewReader(r), nil
}

func getEachTarFile(files map[string][]byte, tr *tar.Reader) error {
	h, err := tr.Next()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	if isTarFileNotDir(h) {
		err := addTarFileToFiles(files, h, tr)
		if err != nil {
			return err
		}
	}
	return getEachTarFile(files, tr)
}

func addTarFileToFiles(files map[string][]byte, h *tar.Header, r *tar.Reader) error {
	var b bytes.Buffer
	_, err := io.Copy(&b, r)
	if err != nil {
		return err
	}
	files[h.Name] = b.Bytes()
	return nil
}

func isTarFileNotDir(h *tar.Header) bool {
	return h.Typeflag != tar.TypeDir
}
