package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path"
)

type PathTransformFunc func(string) PathKey

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		Filename: key,
		Pathname: key,
	}
}

var CASPathTransformFunc = func(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashString := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashString) / blockSize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashString[from:to]
	}

	return PathKey{
		Filename: hashString,
		Pathname: path.Join(paths...),
	}
}

type PathKey struct {
	Filename string
	Pathname string
}

func (p PathKey) FullPath() string {
	return path.Join(p.Pathname, p.Filename)
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

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	return os.Open(pathKey.FullPath())
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)

	if err := os.MkdirAll(pathKey.Pathname, os.ModePerm); err != nil {
		return err
	}

	pathAndFilename := pathKey.FullPath()

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
