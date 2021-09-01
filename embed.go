package ioembed

import (
	"embed"
	"io/fs"
)

func GetFiles(f embed.FS) (map[string][]byte, error) {
	files := map[string][]byte{}
	err := fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
		if d.Type() == fs.ModeDir {
			return nil
		}
		content, err := f.ReadFile(path)
		if err != nil {
			return err
		}
		files[path] = content
		return nil
	})
	return files, err
}
