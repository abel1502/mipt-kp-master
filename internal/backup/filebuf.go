package backup

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileBuf struct {
	ID   string
	Size uint64 // TODO: Perhaps implicit?
}

func NewFileBuf(contentMD5 []byte, size uint64) *FileBuf {
	return &FileBuf{
		ID:   hex.EncodeToString(contentMD5),
		Size: size,
	}
}

func (f *FileBuf) Path(base string) string {
	// TODO: Do I need this separation by the first byte?
	return filepath.Join(base, "files", f.ID[:2], f.ID)
}

func (f *FileBuf) MD5() []byte {
	result, err := hex.DecodeString(f.ID)
	if err != nil {
		panic(fmt.Sprintf("invalid FileBuf ID: %q", f.ID))
	}

	return result
}

func (f *FileBuf) LazyReader(base string) io.ReadCloser {
	return &lazyFileReader{Path: f.Path(base), File: nil}
}

type lazyFileReader struct {
	Path string
	File *os.File
	Done bool
}

func (r *lazyFileReader) Read(p []byte) (int, error) {
	if r.Done {
		return 0, io.EOF
	}

	if r.File == nil {
		f, err := os.Open(r.Path)
		if err != nil {
			return 0, err
		}
		r.File = f
	}

	n, err := r.File.Read(p)

	if err == io.EOF {
		r.File.Close()
		r.File = nil
		r.Done = true
	}

	return n, err
}

func (r *lazyFileReader) Close() error {
	if r.File == nil {
		return os.ErrClosed
	}

	defer func() {
		r.File = nil
		r.Done = true
	}()

	return r.File.Close()
}

type multiReadCloser struct {
	io.Reader
	closers []io.Closer
}

func (r *multiReadCloser) Close() error {
	for _, closer := range r.closers {
		err := closer.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func ChainReader(readClosers ...io.ReadCloser) io.ReadCloser {
	readers := make([]io.Reader, 0, len(readClosers))
	closers := make([]io.Closer, 0, len(readClosers))

	for _, readCloser := range readClosers {
		readers = append(readers, readCloser)
		closers = append(closers, readCloser)
	}

	return &multiReadCloser{
		io.MultiReader(readers...),
		closers,
	}
}
