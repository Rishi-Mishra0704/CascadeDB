package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
)

func CasPathTransformFunc(key string) string {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize

		paths[i] = hashStr[from:to]

	}
	return strings.Join(paths, "/")
}

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

	buf := new(bytes.Buffer)

	io.Copy(buf, r)

	filenameBytes := md5.Sum(buf.Bytes())

	filename := hex.EncodeToString(filenameBytes[:])

	pathAndFileName := pathname + "/" + filename

	f, err := os.Create(pathAndFileName)
	if err != nil {
		return err
	}
	n, err := io.Copy(f, buf)
	if err != nil {
		return err
	}

	log.Printf("written %d bytes to disk: %s", n, pathAndFileName)

	return nil
}
