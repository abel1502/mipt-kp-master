package backup

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"slices"
	"time"

	azblob "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	azcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/abel1502/mipt-kp-m-test/internal/azure"
)

type Repository struct {
	// ContainerURL is the URL of the container to back up
	ContainerURL string `json:"container_url"`
	// LocalPath is the path to the repository's root directory on the local filesystem.
	// FileBufs are stored in the "files" subdirectory (binary files);
	// Snapshots are stored in the "snapshots" subdirectory (json files);
	// Repository-wide metadata is stored in an "info.json" file.
	LocalPath string `json:"local_path"`
	// Revisions are the container snapshots in the chronological order.
	// Note that different revisions in a repository might share some
	// of the blob content pieces.
	Revisions []Snapshot `json:"-"`
}

func NewRepository(containerURL string, localPath string) (*Repository, error) {
	err := os.MkdirAll(localPath, 0755)
	if err != nil {
		return nil, err
	}

	result := &Repository{
		ContainerURL: containerURL,
		LocalPath:    localPath,
		Revisions:    nil,
	}

	err = result.save()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func OpenRepository(localPath string) (*Repository, error) {
	result := &Repository{
		LocalPath: localPath,
	}

	err := result.load()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repository) save() error {
	metadataFile, err := os.Create(path.Join(r.LocalPath, "info.json"))
	if err != nil {
		return err
	}
	defer metadataFile.Close()

	encoder := json.NewEncoder(metadataFile)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(r)
	if err != nil {
		return err
	}

	for _, snapshot := range r.Revisions {
		err = snapshot.save()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) load() error {
	metadataFile, err := os.Open(path.Join(r.LocalPath, "info.json"))
	if err != nil {
		return err
	}
	defer metadataFile.Close()

	decoder := json.NewDecoder(metadataFile)
	err = decoder.Decode(r)
	if err != nil {
		return err
	}

	snapshotDirs, err := os.ReadDir(path.Join(r.LocalPath, "snapshots"))
	if err != nil {
		return err
	}

	r.Revisions = nil
	for _, snapshotIndex := range snapshotDirs {
		snapshot := Snapshot{
			IndexPath: path.Join(r.LocalPath, "snapshots", snapshotIndex.Name()),
		}

		err = snapshot.load()
		if err != nil {
			log.Printf("Warning: Failed to load snapshot %q: %v", snapshot.IndexPath, err)
			continue
		}

		r.Revisions = append(r.Revisions, snapshot)
	}

	return nil
}

func (r *Repository) Close() error {
	return r.save()
}

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
		IndexPath: snapshotPath,
		Blobs:     make(BlobList, 0, len(onlineSnapshot.Blobs)),
	}

	for _, blobInfo := range onlineSnapshot.Blobs {
		// TODO: Also compare LastModified against TakenAt
		// Note: If the (online) blob snapshot was modified after
		// the (online) container snapshot was started,
		// we should abort the process and try again.
		// Also note that, if the blob is deleted before we've
		// finished backing it up, the snapshot is deleted too.

		newBlob, err := r.backupBlob(ctx, client, blobInfo, oldBlobLookup[blobInfo.Name])
		if err != nil {
			return err
		}
		snapshot.Blobs = append(snapshot.Blobs, newBlob)
	}

	r.Revisions = append(r.Revisions, snapshot)

	success = true

	return nil
}

func (r *Repository) backupBlob(
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
		blob, err := DownloadBlob(ctx, r, client, newBlobInfo, *newBlobProps.BlobType, nil)
		return blob, err
	}

	// 2. The blob is unchanged since last time.
	// TODO: This doesn't account for changed metadata! Perhaps also see LastModified?
	if slices.Equal(oldBlob.Common().ContentMD5, newBlobProps.ContentMD5) {
		blob := oldBlob.ShallowClone()
		return blob, nil
	}

	// 3. The blob is updated in a known way
	blob, err := DownloadBlob(ctx, r, client, newBlobInfo, *newBlobProps.BlobType, oldBlob)
	return blob, err
}

func (r *Repository) DownloadBlobRangeAsFileBuf(
	ctx context.Context,
	client *azblob.Client,
	offset uint64,
	size uint64,
) (*FileBuf, error) {
	if size > 4*1024*1024 {
		// TODO: Perhaps correct somehow? Worst case, we'll need to download and recompute everything, which is unpleasant
		log.Printf("Warning: Attempt to download large blob range (%v bytes). ContentMD5 might not work in these scenarios.", size)
	}

	stream, err := client.DownloadStream(ctx, &azblob.DownloadStreamOptions{
		Range: azblob.HTTPRange{
			Offset: int64(offset),
			Count:  int64(size),
		},
		// Why do they even need a *bool?
		RangeGetContentMD5: func(b bool) *bool { return &b }(true),
	})
	if err != nil {
		return nil, err
	}
	defer stream.Body.Close()

	contentMD5 := stream.ContentMD5

	if uint64(*stream.ContentLength) != size {
		panic(fmt.Sprintf("unexpected returned size: want %v, got %v", size, *stream.ContentLength))
	}

	fb := NewFileBuf(contentMD5, size)
	file, err := os.OpenFile(fb.Path(r.LocalPath), os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if errors.Is(err, os.ErrExist) {
		return fb, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, stream.Body)
	if err != nil {
		return nil, err
	}

	return fb, nil
}

// TODO: Manual garbage collection for FileBufs!
