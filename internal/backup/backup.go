package backup

import (
	"time"
)

type BlobData struct {
	// Name is the name of the blob
	Name string
	// Type is the blob type (append/block/page)
	Type BlobType
	// CreatedAt is the time at which the original blob was created.
	// Used for detecting full overwrites
	CreatedAt time.Time
	// SavedAt is the time at which this blob backup was taken
	SavedAt time.Time
	// ValidAt is the time of last check which confirmed the blob remained
	// since SavedAt
	ValidAt time.Time
	// DeletedAt is the time at which a blob was deleted, or nil if it is still up
	DeletedAt *time.Time
	// Metadata is the blob metadata
	Metadata map[string]string
	// ContentMD5 is the MD5 hash of the blob content
	ContentMD5 []byte
	// Content is the blob's content, possibly stored incrementally
	Content BackedUpContent
}

type BackedUpContent interface {
	ExplicitData(blob *BlobData) []byte
}

type ContainerData struct {
	// SavedAt is the time at which this container backup was taken
	SavedAt time.Time
	// Blobs is the list of all blobs included in this backup
	Blobs []BlobData
	// Parent is the previous backup, upon which this one may rely
	// to save storage space using incremental representations of changes
	Parent *ContainerData
}

type BlobType string

const (
	BlobTypeAppendBlob BlobType = "AppendBlob"
	BlobTypeBlockBlob  BlobType = "BlockBlob"
	BlobTypePageBlob   BlobType = "PageBlob"
)

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
