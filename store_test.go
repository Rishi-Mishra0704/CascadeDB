package main

import (
	"bytes"
	"fmt"
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
	if PathKey.Orignal != expectedPath {
		t.Errorf("have %s want %s", PathKey.PathName, expectedOriginalKey)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CasPathTransformFunc,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("some jpg bytes"))

	if err := s.WriteStream("mypicture", data); err != nil {
		t.Error(err)
	}

}
