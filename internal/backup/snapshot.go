package backup

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	azcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/gobwas/glob"
)

type Snapshot struct {
	// SavedAt is the time at which this container backup was taken
	SavedAt time.Time `json:"saved_at"`
	// IndexFile is the path to the snapshot's index file.
	// The composition of the saved blobs is saved there, but not the actual contents
	IndexFile string `json:"-"`
	// Blobs is the list of all blobs included in this backup
	Blobs BlobList `json:"blobs"`
}

type BlobList []Blob

// TODO: Unsure about pointers
var _ json.Marshaler = (BlobList)(nil)
var _ json.Unmarshaler = (*BlobList)(nil)

type annotatedBlob struct {
	Type azcontainer.BlobType `json:"type"`
	Blob json.RawMessage      `json:"blob"`
}

func (l BlobList) MarshalJSON() ([]byte, error) {
	raw := make([]annotatedBlob, 0, len(l))
	for _, blob := range l {
		rawBlob, err := json.Marshal(blob)
		if err != nil {
			return nil, err
		}

		raw = append(raw, annotatedBlob{
			Type: blob.Type(),
			Blob: rawBlob,
		})
	}

	return json.Marshal(raw)
}

func (l *BlobList) UnmarshalJSON(data []byte) error {
	raw := make([]annotatedBlob, 0)
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	*l = make(BlobList, 0, len(raw))
	for _, rawBlob := range raw {
		switch rawBlob.Type {
		case azcontainer.BlobTypeAppendBlob:
			var blob AppendBlob
			err = json.Unmarshal(rawBlob.Blob, &blob)
			*l = append(*l, &blob)

		case azcontainer.BlobTypeBlockBlob:
			var blob BlockBlob
			err = json.Unmarshal(rawBlob.Blob, &blob)
			*l = append(*l, &blob)

		case azcontainer.BlobTypePageBlob:
			var blob PageBlob
			err = json.Unmarshal(rawBlob.Blob, &blob)
			*l = append(*l, &blob)

		default:
			return fmt.Errorf("unknown blob type: %s", rawBlob.Type)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Snapshot) save() error {
	indexFile, err := os.Create(s.IndexFile)
	if err != nil {
		return err
	}
	defer indexFile.Close()

	encoder := json.NewEncoder(indexFile)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(s)
	if err != nil {
		return err
	}

	// TODO: Maybe do something with filebufs?

	return nil
}

func (s *Snapshot) load() error {
	indexFile, err := os.Open(s.IndexFile)
	if err != nil {
		return err
	}
	defer indexFile.Close()

	decoder := json.NewDecoder(indexFile)
	err = decoder.Decode(s)
	if err != nil {
		return err
	}

	// TODO: Maybe do something with filebufs?

	return nil
}

// TODO: Export files into regular FS by a glob
func (s *Snapshot) ExportByGlob(
	ctx context.Context,
	repo *Repository,
	targets glob.Glob,
	destination string,
	flat bool,
) error {
	for _, blob := range s.Blobs {
		blobName := blob.Common().Name

		if !targets.Match(blobName) {
			continue
		}

		var dstPath string
		if flat {
			dstPath = filepath.Join(destination, path.Base(blobName))
		} else {
			dstPath = filepath.Join(destination, filepath.FromSlash(blobName))
		}

		blobReader := blob.Export(ctx, repo)
		defer blobReader.Close()

		outWriter, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer outWriter.Close()

		_, err = io.Copy(outWriter, blobReader)
		if err != nil {
			return err
		}

		log.Printf("Exported %q to %q", blobName, dstPath)
	}

	return nil
}
