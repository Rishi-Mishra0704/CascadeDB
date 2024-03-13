package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func CasPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize

		paths[i] = hashStr[from:to]

	}

	return PathKey{
		PathName: strings.Join(paths, "/"),
		Orignal:  hashStr,
	}
}

type PathTransformFunc func(string) PathKey

type PathKey struct {
	PathName string
	Orignal  string
}

func (p PathKey) FileName() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.Orignal)
}

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
	PathKey := s.PathTransformFunc(key)
	if err := os.MkdirAll(PathKey.PathName, os.ModePerm); err != nil {
		return err
	}

	pathAndFileName := PathKey.FileName()

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
