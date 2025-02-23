package backup

import (
	"context"
	"slices"
	"time"

	azblob "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	azcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/abel1502/mipt-kp-m-test/internal/azure"
)

type Snapshot struct {
	// SavedAt is the time at which this container backup was taken
	SavedAt time.Time
	// Blobs is the list of all blobs included in this backup
	Blobs []Blob
}

type Repository struct {
	// ContainerURL is the URL of the container to back up
	ContainerURL string
	// Revisions are the container snapshots in the chronological order.
	// Note that different revisions in a repository might share some
	// of the blob content pieces.
	Revisions []Snapshot
}

func NewRepository(containerURL string) *Repository {
	return &Repository{
		ContainerURL: containerURL,
		Revisions:    nil,
	}
}

func (r *Repository) TakeSnapshot(ctx context.Context) error {
	client, err := azure.OpenClient(r.ContainerURL)
	if err != nil {
		return err
	}

	onlineSnapshot, err := azure.TakeSnapshot(ctx, client)
	defer onlineSnapshot.Delete(ctx) // TODO: Other context?

	oldBlobLookup := make(map[string]Blob)
	if len(r.Revisions) > 0 {
		lastRevision := &r.Revisions[len(r.Revisions)-1]

		for _, blob := range lastRevision.Blobs {
			// TODO: This may be wrong -- if so, do &lastRevision.Blobs[i] instead
			oldBlobLookup[blob.Common().Name] = blob
		}
	}

	// TODO
	blobs := make([]Blob, 0, len(onlineSnapshot.Blobs))

	for _, blob := range onlineSnapshot.Blobs {
		// TODO: Also compare LastModified against TakenAt
		newBlob, err := backupBlob(ctx, client, blob, oldBlobLookup[blob.Name])
		if err != nil {
			return err
		}
		blobs = append(blobs, newBlob)
	}

	r.Revisions = append(r.Revisions, Snapshot{
		SavedAt: onlineSnapshot.TakenAt,
		Blobs:   blobs,
	})

	return nil
}

// Note: oldBlob may be nil, meaning this blob hasn't been backed up in the last revision
func backupBlob(ctx context.Context, client *azcontainer.Client, newBlob azure.BlobInfo, oldBlob Blob) (Blob, error) {
	newBlobClient, err := client.NewBlobClient(newBlob.Name).WithSnapshot(newBlob.Snapshot)
	if err != nil {
		return nil, err
	}

	newProps, err := newBlobClient.GetProperties(ctx, nil)
	if err != nil {
		return nil, err
	}

	// There are three options here:

	// 1. The blob is created fresh.
	//    The old one is either overwritten (different creation time)
	//    or didn't exist (nil)
	if oldBlob == nil || (oldBlob.Common().Timestamps.CreatedAt != *newProps.CreationTime) {
		blob, err := blobFullDownload(ctx, newBlobClient)
		return blob, err
	}

	// 2. The blob is unchanged since last time.
	// TODO: This doesn't account for changed metadata! Perhaps also see LastModified?
	if slices.Equal(oldBlob.Common().ContentMD5, newProps.ContentMD5) {
		blob := oldBlob.ShallowClone()
		return blob, nil
	}

	// 3. The blob is updated in a known way
	blob, err := blobIncrementalDownload(ctx, newBlobClient, oldBlob)
	return blob, err
}

func blobFullDownload(ctx context.Context, blobClient *azblob.Client) (Blob, error) {
	// TODO
	panic("not implemented")

	return nil, nil
}

func blobIncrementalDownload(ctx context.Context, blobClient *azblob.Client, oldBlob Blob) (Blob, error) {
	// TODO
	panic("not implemented")

	return nil, nil
}
