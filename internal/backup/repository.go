package backup

import (
	"context"
	"os"
	"path"
	"slices"
	"time"

	azcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/abel1502/mipt-kp-m-test/internal/azure"
)

type Repository struct {
	// ContainerURL is the URL of the container to back up
	ContainerURL string
	// Revisions are the container snapshots in the chronological order.
	// Note that different revisions in a repository might share some
	// of the blob content pieces.
	Revisions []Snapshot
	// LocalPath is the path to the repository's root directory on the local filesystem.
	// FileBufs are stored in the "files" subdirectory (binary files);
	// Snapshots are stored in the "snapshots" subdirectory (json files);
	// Repository-wide metadata is stored in an "info.json" file.
	LocalPath string
}

func NewRepository(containerURL string, localPath string) *Repository {
	return &Repository{
		ContainerURL: containerURL,
		Revisions:    nil,
		LocalPath:    localPath,
	}
}

// TODO: Open repository

func (r *Repository) TakeSnapshot(ctx context.Context) error {
	success := false

	client, err := azure.OpenClient(r.ContainerURL)
	if err != nil {
		return err
	}

	onlineSnapshot, err := azure.TakeSnapshot(ctx, client)
	if err != nil {
		return err
	}
	defer onlineSnapshot.Delete(ctx) // TODO: Other context?

	oldBlobLookup := make(map[string]Blob)
	if len(r.Revisions) > 0 {
		lastRevision := &r.Revisions[len(r.Revisions)-1]

		for _, blob := range lastRevision.Blobs {
			oldBlobLookup[blob.Common().Name] = blob
		}
	}

	snapshotPath := path.Join(
		r.LocalPath,
		"snapshots",
		onlineSnapshot.TakenAt.Format(time.RFC3339)+".json",
	)
	err = os.MkdirAll(path.Dir(snapshotPath), 0755)
	if err != nil {
		return err
	}
	// TODO: Maybe not remove the incomplete snapshot directory? Though then it would be loaded next time.
	defer func() {
		if !success {
			_ = os.Remove(snapshotPath)
		}
	}()

	snapshot := Snapshot{
		SavedAt:   onlineSnapshot.TakenAt,
		Blobs:     make([]Blob, 0, len(onlineSnapshot.Blobs)),
		IndexPath: snapshotPath,
	}

	for _, blobInfo := range onlineSnapshot.Blobs {
		// TODO: Also compare LastModified against TakenAt
		// Note: If the (online) blob snapshot was modified after
		// the (online) container snapshot was started,
		// we should abort the process and try again.
		// Also note that, if the blob is deleted before we've
		// finished backing it up, the snapshot is deleted too.

		// TODO: Also pass repository object to enable storing FileBufs directly to fs
		newBlob, err := backupBlob(ctx, client, blobInfo, oldBlobLookup[blobInfo.Name])
		if err != nil {
			return err
		}
		snapshot.Blobs = append(snapshot.Blobs, newBlob)
	}

	r.Revisions = append(r.Revisions, snapshot)

	success = true

	return nil
}

func backupBlob(
	ctx context.Context,
	client *azcontainer.Client,
	newBlobInfo azure.BlobInfo,
	oldBlob Blob,
) (Blob, error) {
	newBlobClient, err := client.NewBlobClient(newBlobInfo.Name).WithSnapshot(newBlobInfo.Snapshot)
	if err != nil {
		return nil, err
	}

	newBlobProps, err := newBlobClient.GetProperties(ctx, nil)
	if err != nil {
		return nil, err
	}

	// There are three options here:

	// 1. The blob is created fresh.
	//    The old one is either overwritten (different creation time)
	//    or didn't exist (nil)
	if oldBlob == nil || (oldBlob.Common().Timestamps.CreatedAt != *newBlobProps.CreationTime) {
		blob, err := DownloadBlob(ctx, client, newBlobInfo, *newBlobProps.BlobType, nil)
		return blob, err
	}

	// 2. The blob is unchanged since last time.
	// TODO: This doesn't account for changed metadata! Perhaps also see LastModified?
	if slices.Equal(oldBlob.Common().ContentMD5, newBlobProps.ContentMD5) {
		blob := oldBlob.ShallowClone()
		return blob, nil
	}

	// 3. The blob is updated in a known way
	blob, err := DownloadBlob(ctx, client, newBlobInfo, *newBlobProps.BlobType, oldBlob)
	return blob, err
}

// TODO: Lookup/download FileBufs
