package storage

import (
	"fmt"
	"io"
	"os"
)

type LocalStorage struct {
}

func (l *LocalStorage) UploadFile(fileName string, r io.Reader) (err error) {
	f, err := os.Create(fmt.Sprintf("storage/%s", fileName))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)

	if err != nil {
		return err
	}

	return nil

}
