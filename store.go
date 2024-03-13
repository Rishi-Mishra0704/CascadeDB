package main

import (
	"bytes"
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
		FileName: hashStr,
	}
}

type PathTransformFunc func(string) PathKey

type PathKey struct {
	PathName string
	FileName string
}

func (p PathKey) FullPth() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.FileName)
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

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	PathKey := s.PathTransformFunc(key)

	return os.Open(PathKey.FullPth())
}

func (s *Store) WriteStream(key string, r io.Reader) error {
	PathKey := s.PathTransformFunc(key)
	if err := os.MkdirAll(PathKey.PathName, os.ModePerm); err != nil {
		return err
	}

	pathAndFileName := PathKey.FullPth()

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
