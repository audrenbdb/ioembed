package ioembed

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"embed"
	"io"
	"io/fs"
)

func GetGZippedTarFiles(fs embed.FS) (map[string][]byte, error) {
	f, err := getFirstFileFromEmbedFS(fs)
	if err != nil {
		return nil, err
	}
	r, err := newGzipTarReader(f)
	if err != nil {
		return nil, err
	}
	files := map[string][]byte{}
	err = getEachTarFile(files, r)
	return files, err
}

func getFirstFileFromEmbedFS(e embed.FS) (io.Reader, error) {
	var r io.Reader
	err := fs.WalkDir(e, ".", func(path string, d fs.DirEntry, err error) error {
		if d.Type() == fs.ModeDir || r != nil{
			return nil
		}
		b, err := e.ReadFile(path)
		if err != nil {
			return err
		}
		r = bytes.NewBuffer(b)
		return nil
	})
	return r, err
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
