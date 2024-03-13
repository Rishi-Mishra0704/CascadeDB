package main

import (
	"io"
	"log"
	"os"
)

type PathTransformFunc func(string) string

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

var DefaultTransformFunc = func(key string) string {
	return key
}

func (s *Store) WriteStream(key string, r io.Reader) error {
	pathname := s.PathTransformFunc(key)
	if err := os.MkdirAll(pathname, os.ModePerm); err != nil {
		return err
	}

	filename := "somefilename"

	pathAndFileName := pathname + "/" + filename

	f, err := os.Create(pathAndFileName)
	if err != nil {
		return err
	}
	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("written %d bytes to disk: %s", n, pathAndFileName)

	return nil
}
