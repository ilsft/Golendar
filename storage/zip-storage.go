package storage

import (
	"archive/zip"
	"errors"

	"io"
	"os"
)

type ZipStorage struct {
	*Storage
}

func NewZipStorage(filename string) *ZipStorage {
	return &ZipStorage{
		&Storage{filename: filename},
	}
}

func (z *ZipStorage) Save(data []byte) error {
	f, err := os.Create(z.GetFilename())
	if err != nil {
		return err
	}
	defer f.Close()

	zw := zip.NewWriter(f)

	w, err := zw.Create("data")
	if err != nil {
		zw.Close()
		return err
	}

	_, err = w.Write(data)
	if err != nil {
		zw.Close()
		return err
	}

	err = zw.Close()
	if err != nil {
		return err
	}

	return nil
}

func (z *ZipStorage) Load() ([]byte, error) {
	r, err := zip.OpenReader(z.GetFilename())
	if err != nil {
		return nil, err
	}
	defer r.Close()

	if len(r.File) == 0 {
		return nil, errors.New("архив пуст")
	}

	file := r.File[0]
	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return io.ReadAll(rc)
}
