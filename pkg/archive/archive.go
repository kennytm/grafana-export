// Copyright 2020 PingCAP, Inc. Licensed under Apache-2.0.

package archive

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type Archive interface {
	io.Closer
	WriteToFile(name string, writer func(io.Writer) error) error
}

type Directory struct {
	path string
}

func NewDirectory(path string) (dir Directory, err error) {
	dir.path = path
	err = os.MkdirAll(path, 0755)
	return
}

func (d Directory) WriteToFile(name string, writer func(io.Writer) error) (err error) {
	var f *os.File
	f, err = os.Create(filepath.Join(d.path, name))
	if err != nil {
		return err
	}
	defer func() {
		errClose := f.Close()
		if err == nil {
			err = errClose
		}
	}()

	return writer(f)
}

func (d Directory) Close() error {
	return nil
}

type ZipFile struct {
	mu     sync.Mutex
	writer *zip.Writer
}

func NewZipFile(path string) (*ZipFile, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return &ZipFile{writer: zip.NewWriter(file)}, nil
}

func (zf *ZipFile) WriteToFile(name string, writer func(io.Writer) error) error {
	zf.mu.Lock()
	defer zf.mu.Unlock()

	f, err := zf.writer.Create(name)
	if err != nil {
		return err
	}
	return writer(f)
}

func (zf *ZipFile) Close() error {
	zf.mu.Lock()
	defer zf.mu.Unlock()
	return zf.writer.Close()
}
