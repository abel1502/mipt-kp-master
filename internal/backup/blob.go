package backup

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

// TODO: Split into different files?

type Blob interface {
	Type() BlobType
	Common() *CommonBlob
	ShallowClone() Blob
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
