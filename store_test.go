package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "mombestpicture"
	pathname := CasPathTransformFunc(key)
	fmt.Println(pathname)
	expectedPath := "cf5d4/b01c4/d9438/c22c5/6c832/f83bd/3e8c6/304f9"
	if pathname != expectedPath {
		t.Errorf("have %s want %s", pathname, expectedPath)
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
