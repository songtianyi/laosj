package storage

import (
	"github.com/songtianyi/rrframework/utils"
	"os"
)

type LocalDiskStorage struct {
	Dir string
}

func NewLocalDiskStorage(dir string) StorageWrapper {
	s := &LocalDiskStorage{
		Dir: dir,
	}
	return s
}

func (s *LocalDiskStorage) Save(data []byte) error {
	// random name
	path := s.Dir + rrutils.NewV4().String()

	//open a file for writing
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(data); err != nil {
		return err
	}
	return nil

}
