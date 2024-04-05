package store

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
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

func (p PathKey) FullPath() string {
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

func (p PathKey) firstPathName() string {
	paths := strings.Split(p.PathName, "/")

	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)

	_, err := os.Stat(pathKey.FullPath())
	if err == fs.ErrNotExist {
		return false
	}
	return true
}

func (s *Store) Post(key string, content io.Reader) error {
	// Transform the key into a PathKey
	pathKey := s.PathTransformFunc(key)

	// Create the full path for storing the content
	pathAndFileName := pathKey.FullPath()

	// Create any necessary directories
	if err := os.MkdirAll(pathKey.PathName, os.ModePerm); err != nil {
		return err
	}

	// Create or open the file for writing
	f, err := os.Create(pathAndFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write the content to the file
	n, err := io.Copy(f, content)
	if err != nil {
		return err
	}

	log.Printf("written %d bytes to disk: %s", n, pathAndFileName)

	return nil
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)
	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.FileName)
	}()
	return os.RemoveAll(pathKey.firstPathName())
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

	return os.Open(PathKey.FullPath())
}

func (s *Store) WriteStream(key string, r io.Reader) error {
	PathKey := s.PathTransformFunc(key)
	if err := os.MkdirAll(PathKey.PathName, os.ModePerm); err != nil {
		return err
	}

	pathAndFileName := PathKey.FullPath()

	f, err := os.Create(pathAndFileName)
	if err != nil {
		return err
	}

	defer f.Close()
	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("written %d bytes to disk: %s", n, pathAndFileName)

	return nil
}
