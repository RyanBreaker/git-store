package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const key = "momsspecials"

func TestDefaultPathTransformFunc(t *testing.T) {
	pathKey := DefaultPathTransformFunc("foobar")
	assert.Equal(t, "foobar", pathKey.Filename)
	assert.Equal(t, "foobar", pathKey.Pathname)
}

func TestCASPathTransformFunc(t *testing.T) {
	originalKey := "momsbestpicture"

	expectedFilename := "6804429f74181a63c50c3d81d733a12f14a353ff"
	expectedPathname := "68044/29f74/181a6/3c50c/3d81d/733a1/2f14a/353ff"

	pathKey := CASPathTransformFunc(originalKey)

	assert.Equal(t, expectedFilename, pathKey.Filename)
	assert.Equal(t, expectedPathname, pathKey.Pathname)
}

func TestStore_Read(t *testing.T) {
	opts := StoreOps{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)

	data := []byte("some jpeg bytes")
	err := s.writeStream(key, bytes.NewReader(data))
	assert.Nil(t, err)

	//r, err := s.Read(key)
	//assert.Nil(t, err)
	//
	//b, err := io.ReadAll(r)
	//assert.Equal(t, data, b)
	//
	//teardown(t)
}

func TestStore_Delete(t *testing.T) {
	opts := StoreOps{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)

	data := []byte("some jpeg bytes")
	err := s.writeStream(key, bytes.NewReader(data))
	assert.Nil(t, err)

	err = s.Delete(key)
	assert.Nil(t, err)

	teardown(t)
}

func TestStore_Has(t *testing.T) {
	opts := StoreOps{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)

	assert.False(t, s.Has(key))

	data := []byte("some jpeg bytes")
	err := s.writeStream(key, bytes.NewReader(data))
	assert.Nil(t, err)

	assert.True(t, s.Has(key))

	teardown(t)
}

func TestPathKey_RootPath(t *testing.T) {
	pathKey := CASPathTransformFunc(key)
	assert.Equal(t, "ff254", pathKey.RootPath())
}

func teardown(t *testing.T) {
	_ = os.RemoveAll(DefaultRootName)
	_, err := os.Stat(DefaultRootName)
	assert.NotNil(t, err)
}
