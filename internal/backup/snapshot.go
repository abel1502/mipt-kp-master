package backup

import (
	"time"
)

type Snapshot struct {
	// SavedAt is the time at which this container backup was taken
	SavedAt time.Time
	// Blobs is the list of all blobs included in this backup
	Blobs []Blob
	// LocalPath is the path to the snapshot's index file.
	// The composition of the saved blobs is saved there, but not the actual contents
	IndexPath string
}
