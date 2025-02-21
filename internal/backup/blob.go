package backup

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

// TODO: Split into different files?

type Blob interface {
	Type() BlobType
	Common() *CommonBlob
}

type BlobType string

const (
	BlobTypeAppend BlobType = "AppendBlob"
	BlobTypeBlock  BlobType = "BlockBlob"
	BlobTypePage   BlobType = "PageBlob"
)

type CommonBlob struct {
	// Name is the name of the blob
	Name string
	// Timestamps stores various timestamps related to the blob
	Timestamps struct {
		// TODO: More?
		CreatedAt   time.Time
		SavedAt     time.Time
		LastUpdated time.Time
	}
	// ContentMD5 is the MD5 hash of the blob content
	ContentMD5 []byte
	// ETag is a tag that is updated when the blob is modified in any way
	ETag string
	// ContentSize is the total size of the blob
	ContentSize uint64
	// Metadata is the blob metadata
	Metadata map[string]string
	// Properties are the blob properties
	Properties container.BlobProperties // TODO: Remove?
}

type BlockBlob struct {
	CommonBlob
	// Fragments is the list of blocks that make up the blob
	Fragments []*BlockBlobFragment
}

func (*BlockBlob) Type() BlobType {
	return BlobTypeBlock
}

func (b *BlockBlob) Common() *CommonBlob {
	return &b.CommonBlob
}

type BlockBlobFragment struct {
	// ID is the base64-encoded block ID
	ID string
	// Content is the block data
	Content []byte
}

type AppendBlob struct {
	CommonBlob
	// Fragments is the single-linked list of this blob's fragments
	Fragments *AppendBlobFragment
}

func (*AppendBlob) Type() BlobType {
	return BlobTypeAppend
}

func (a *AppendBlob) Common() *CommonBlob {
	return &a.CommonBlob
}

type AppendBlobFragment struct {
	// LastChunk is the last chunk of this blob
	LastChunk []byte
	// Previous is the preceding fragment
	Previous *AppendBlobFragment
}

type PageBlob struct {
	CommonBlob
	// Fragments is the list of this blob's pages
	Fragments []*PageBlobFragment
}

func (*PageBlob) Type() BlobType {
	return BlobTypePage
}

func (p *PageBlob) Common() *CommonBlob {
	return &p.CommonBlob
}

// TODO: Support clearing parts in the middle of a page range?
// We don't always want to keep the entire history, though.
// For that scenario, I'll need
type PageBlobFragment struct {
	// Offset is the fragment offset (512-bytes-aligned)
	Offset uint64
	// Content is the fragment data (512-bytes-aligned in size)
	Content []byte
	// ContentMD5 is the MD5 hash of the fragment
	ContentMD5 []byte
}
