package server

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
	if PathKey.FileName != expectedOriginalKey {
		t.Errorf("have %s want %s", PathKey.FileName, expectedOriginalKey)
	}
}

func TestDeleteFunc(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CasPathTransformFunc,
	}
	s := NewStore(opts)

	key := "test"

	data := []byte("some jpg in bytes")

	if _, err := s.WriteStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

func TestStore(t *testing.T) {
	s := newStore()
	defer teardown(t, s)
	for i := 0; i <= 50; i++ {

		key := fmt.Sprintf("test%d", i)

		data := generateRandomBytes(100)

		if _, err := s.WriteStream(key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); !ok {
			t.Errorf("expected to have key %s", key)
		}

		r, err := s.Read(key)
		if err != nil {
			t.Error(err)
		}

		b, _ := io.ReadAll(r)
		if string(b) != string(data) {
			t.Errorf("want %s have %s", data, b)
		}

		if err := s.Delete(key); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); ok {
			t.Errorf("expected to NOT have key %s", key)
		}

	}
}
func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CasPathTransformFunc,
	}
	return NewStore(opts)
}

func teardown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}

func generateRandomBytes(n int) []byte {
	b := make([]byte, n)
	return b
}
