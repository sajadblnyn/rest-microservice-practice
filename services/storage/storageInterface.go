package storage

import "io"

type StorageInterFace interface {
	UploadFile(fileName string, r io.Reader) (err error)
}

func NewStorage() StorageInterFace {
	return &LocalStorage{}
}
