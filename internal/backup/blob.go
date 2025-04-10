package backup

import (
	"context"
	"fmt"
	"io"
	"time"

	azblob "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	azcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/abel1502/mipt-kp-m-test/internal/azure"
)

type Blob interface {
	Type() azcontainer.BlobType
	Common() *CommonBlob
	ShallowClone() Blob
	Export(ctx context.Context, repo *Repository) io.ReadCloser
	// TODO: Save/Load metadata to disk; restore references to fragments?
}

type CommonBlob struct {
	// Name is the name of the blob
	Name string `json:"name"`
	// Timestamps stores various timestamps related to the blob
	Timestamps struct {
		// TODO: More?
		CreatedAt   time.Time `json:"created_at"`
		SavedAt     time.Time `json:"saved_at"`
		LastUpdated time.Time `json:"last_updated"`
	} `json:"timestamps"`
	// ContentMD5 is the MD5 hash of the blob content
	ContentMD5 []byte `json:"content_md5"`
	// ETag is a tag that is updated when the blob is modified in any way
	ETag string `json:"etag"`
	// ContentSize is the total size of the blob
	ContentSize uint64 `json:"content_size"`
	// Metadata is the blob metadata
	Metadata map[string]*string `json:"metadata"`
	// TODO: Also include?
	// Tags are the blob tags
	// Tags []string `json:"tags"`
	// TODO: Even more properties?
}

func DownloadBlob(
	ctx context.Context,
	repo *Repository,
	client *azcontainer.Client,
	blobInfo azure.BlobInfo,
	blobType azcontainer.BlobType,
	oldBlob Blob,
) (Blob, error) {
	switch blobType {
	case azblob.BlobTypeAppendBlob:
		oldBlobTyped, ok := oldBlob.(*AppendBlob)
		// Note: if oldBlob was nil, the cast apparently fails, so we account for it here
		if oldBlob == nil {
			ok = true
		}
		if !ok {
			return nil, fmt.Errorf("invalid old blob type: want *backup.AppendBlob, got %T", oldBlob)
		}
		blob, err := DownloadAppendBlob(ctx, repo, client, blobInfo.Name, blobInfo.Snapshot, oldBlobTyped)
		return blob, err

	case azblob.BlobTypeBlockBlob:
		oldBlobTyped, ok := oldBlob.(*BlockBlob)
		// Note: if oldBlob was nil, the cast apparently fails, so we account for it here
		if oldBlob == nil {
			ok = true
		}
		if !ok {
			return nil, fmt.Errorf("invalid old blob type: want *backup.BlockBlob, got %T", oldBlob)
		}
		blob, err := DownloadBlockBlob(ctx, repo, client, blobInfo.Name, blobInfo.Snapshot, oldBlobTyped)
		return blob, err

	case azblob.BlobTypePageBlob:
		oldBlobTyped, ok := oldBlob.(*PageBlob)
		// Note: if oldBlob was nil, the cast apparently fails, so we account for it here
		if oldBlob == nil {
			ok = true
		}
		if !ok {
			return nil, fmt.Errorf("invalid old blob type: want *backup.PageBlob, got %T", oldBlob)
		}
		blob, err := DownloadPageBlob(ctx, repo, client, blobInfo.Name, blobInfo.Snapshot, oldBlobTyped)
		return blob, err
	}

	panic(fmt.Sprintf("invalid blob type: %v", blobType))
}

func downloadCommon(
	ctx context.Context,
	client *azblob.Client,
	name string,
) (*CommonBlob, error) {
	props, err := client.GetProperties(ctx, nil)
	if err != nil {
		return nil, err
	}

	result := &CommonBlob{
		Name:       name,
		ContentMD5: props.ContentMD5,
		ETag:       string(*props.ETag),
		Metadata:   props.Metadata,
	}

	result.Timestamps.CreatedAt = *props.CreationTime
	result.Timestamps.LastUpdated = *props.LastModified
	result.Timestamps.SavedAt = time.Now() // TODO: ?

	stream, err := client.DownloadStream(ctx, nil)
	if err != nil {
		return nil, err
	}
	// We don't actually care about the body at all,
	// might as well close it right away to not waste the traffic
	stream.Body.Close()

	result.ContentSize = uint64(*stream.ContentLength)

	return result, nil
}
