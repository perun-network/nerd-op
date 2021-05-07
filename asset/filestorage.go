// SPDX-License-Identifier: Apache-2.0

package asset

import (
	"fmt"
	"io"
	"io/fs"
	"math/big"
	"os"
)

type FileStorage struct {
	dir       fs.FS
	extension string
}

func NewFileStorage(dirpath string) (*FileStorage, error) {
	dir := os.DirFS(dirpath)
	dirf, err := dir.Open(".")
	if err != nil {
		return nil, fmt.Errorf("opening directory '%s': %w", dirpath, err)
	}
	defer dirf.Close()
	if dirstat, err := dirf.Stat(); err != nil {
		return nil, fmt.Errorf("reading stats '%s': %w", dirpath, err)
	} else if !dirstat.IsDir() {
		return nil, fmt.Errorf("file '%s' is not a directory", dirpath)
	}

	return &FileStorage{dir: dir}, nil
}

func (s *FileStorage) SetExtension(ext string) {
	s.extension = ext
}

func (s *FileStorage) dotExt() string {
	if s.extension == "" {
		return ""
	}
	return "." + s.extension
}

func (s *FileStorage) Get(id *big.Int) ([]byte, error) {
	fname := id.Text(10) + s.dotExt()
	f, err := s.dir.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("opening file '%s': %w", fname, err)
	}
	defer f.Close()

	return io.ReadAll(f)
}
