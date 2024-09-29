package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	originalKey := "momsbestpicture"

	expectedFilename := "6804429f74181a63c50c3d81d733a12f14a353ff"
	expectedPathname := "68044/29f74/181a6/3c50c/3d81d/733a1/2f14a/353ff"

	pathKey := CASPathTransformFunc(originalKey)

	assert.Equal(t, expectedFilename, pathKey.Filename)
	assert.Equal(t, expectedPathname, pathKey.Pathname)
}

func TestStore(t *testing.T) {
	opts := StoreOps{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("some jpeg bytes"))
	if err := s.writeStream("test-store", data); err != nil {
		t.Error(err)
	}
}
