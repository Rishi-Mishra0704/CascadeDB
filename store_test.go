package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "mombestpicture"
	PathKey := CasPathTransformFunc(key)
	fmt.Println(PathKey)
	expectedOriginalKey := "cf5d4b01c4d9438c22c56c832f83bd3e8c6304f9"
	expectedPath := "cf5d4/b01c4/d9438/c22c5/6c832/f83bd/3e8c6/304f9"
	if PathKey.PathName != expectedPath {
		t.Errorf("have %s want %s", PathKey.PathName, expectedPath)
	}
	if PathKey.FileName != expectedPath {
		t.Errorf("have %s want %s", PathKey.PathName, expectedOriginalKey)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CasPathTransformFunc,
	}
	s := NewStore(opts)

	key := "momspecials"

	data := []byte("some jpg bytes")

	if err := s.WriteStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)
	if string(b) != string(data) {
		t.Errorf("want %s have %s", data, b)
	}
}
