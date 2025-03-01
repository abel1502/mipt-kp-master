package backup

import (
	"encoding/hex"
	"fmt"
	"path"
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

func (f *FileBuf) RelPath() string {
	// TODO: Do I need this separation by the first byte?
	return path.Join("files", f.ID[:2], f.ID)
}

func (f *FileBuf) MD5() []byte {
	result, err := hex.DecodeString(f.ID)
	if err != nil {
		panic(fmt.Sprintf("invalid FileBuf ID: %q", f.ID))
	}

	return result
}
