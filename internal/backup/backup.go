package backup

import (
	"context"
	"time"

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

	oldBlobLookup := make(map[string]*Blob)
	if len(r.Revisions) > 0 {
		lastRevision := &r.Revisions[len(r.Revisions)-1]

		for _, blob := range lastRevision.Blobs {
			// TODO: This may be wrong -- if so, do &lastRevision.Blobs[i] instead
			oldBlobLookup[blob.Common().Name] = &blob
		}
	}

	// TODO
	blobs := make([]Blob, 0, len(onlineSnapshot.Blobs))

	for _, blob := range onlineSnapshot.Blobs {
		// TODO: Also compare LastModified against TakenAt
		blobs = append(blobs, backupBlob(blob, oldBlobLookup[blob.Name]))
	}

	r.Revisions = append(r.Revisions, Snapshot{
		SavedAt: onlineSnapshot.TakenAt,
		Blobs:   blobs,
	})

	return nil
}

func backupBlob(newBlob azure.BlobInfo, oldBlob *Blob) Blob {
	// TODO
	return nil
}

/*
const metadataFile string = "metadata.json"

// Open opens an existing backup directory
func Open(root string) (*BackupDir, error) {
	mdPath := path.Join(root, metadataFile)
	mdFile, err := os.Open(mdPath)
	if err != nil {
		return nil, err
	}
	defer mdFile.Close()

	b := &BackupDir{
		Root: root,
	}

	err = json.NewDecoder(mdFile).Decode(&b.Metadata)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// New creates a BackupDir for a given container.
// Note that by default the backup is empty (equivalent to T-infinity),
// and you need to call PullChanges to populate it
func New(containerURL string, root string) (*BackupDir, error) {
	// TODO
	return nil, nil
}

// Flush saves the metadata to disk.
func (b *BackupDir) Flush() error {
	mdPath := path.Join(b.Root, metadataFile)
	mdFile, err := os.Create(mdPath)
	if err != nil {
		return err
	}
	defer mdFile.Close()

	return json.NewEncoder(mdFile).Encode(b.Metadata)
}

// TODO: Maybe remove?
func (b *BackupDir) Close() {
	_ = b.Flush()
}

// NewIncrement creates a new incremental backup based on this one.
// Note that by default the new backup is empty, and you need to
// call PullChanges to populate it
func (b *BackupDir) NewIncrement(root string) (*BackupDir, error) {
	// TODO
	return nil, nil
}

// Clone creates a new backup identical to this one.
func (b *BackupDir) Clone(root string) (*BackupDir, error) {
	// TODO
	return nil, nil
}

// PullChanges incorporates any changes from the container
// made since the currently saved timestamp. Note that this
// updates the existing backup; if you wish to preserve the
// original, either clone it or use NewIncrement.
func (b *BackupDir) PullChanges() error {
	// TODO
	return nil
}

// MergeHistory turns an incremental backup into a full one
// by incorporating all the changes since the last full backup
func (b *BackupDir) MergeHistory() error {
	// TODO
	return nil
}
*/
