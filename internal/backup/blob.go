package backup

import (
	"context"
	"fmt"
	"time"

	azblob "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	azcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/abel1502/mipt-kp-m-test/internal/azure"
)

type Blob interface {
	Type() azcontainer.BlobType
	Common() *CommonBlob
	ShallowClone() Blob
}

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
	// TODO: Store too?
	// Properties are the blob properties
	// Properties container.BlobProperties
}

func DownloadBlob(
	ctx context.Context,
	client *azcontainer.Client,
	blobInfo azure.BlobInfo,
	blobType azcontainer.BlobType,
	oldBlob Blob,
) (Blob, error) {
	switch blobType {
	case azblob.BlobTypeAppendBlob:
		oldBlob, ok := oldBlob.(*AppendBlob)
		if !ok {
			return nil, fmt.Errorf("invalid old blob type: want AppendBlob, got %T", oldBlob)
		}
		blob, err := DownloadAppendBlob(ctx, client, blobInfo.Name, blobInfo.Snapshot, oldBlob)
		return blob, err

	case azblob.BlobTypeBlockBlob:
		oldBlob, ok := oldBlob.(*BlockBlob)
		if !ok {
			return nil, fmt.Errorf("invalid old blob type: want BlockBlob, got %T", oldBlob)
		}
		blob, err := DownloadBlockBlob(ctx, client, blobInfo.Name, blobInfo.Snapshot, oldBlob)
		return blob, err

	case azblob.BlobTypePageBlob:
		oldBlob, ok := oldBlob.(*PageBlob)
		if !ok {
			return nil, fmt.Errorf("invalid old blob type: want PageBlob, got %T", oldBlob)
		}
		blob, err := DownloadPageBlob(ctx, client, blobInfo.Name, blobInfo.Snapshot, oldBlob)
		return blob, err
	}

	panic(fmt.Sprintf("invalid blob type: %v", blobType))
}
