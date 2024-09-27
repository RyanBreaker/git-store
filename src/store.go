package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type PathTransformFunc func(string) string

var DefaultPathTransformFunc = func(key string) string {
	return key
}

type StoreOps struct {
	PathTransformFunc
}

type Store struct {
	StoreOps
}

func NewStore(opts StoreOps) *Store {
	return &Store{
		opts,
	}
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathname := s.PathTransformFunc(key)

	if err := os.MkdirAll(pathname, os.ModePerm); err != nil {
		return err
	}

	filename := "somefilename"

	pathAndFilename := fmt.Sprintf("%s/%s", pathname, filename)

	f, err := os.Create(pathAndFilename)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("Wrote %d bytes: %s\n", n, pathAndFilename)

	return nil
}
