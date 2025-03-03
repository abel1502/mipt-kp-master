package backup

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	azcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

type Snapshot struct {
	// SavedAt is the time at which this container backup was taken
	SavedAt time.Time `json:"saved_at"`
	// LocalPath is the path to the snapshot's index file.
	// The composition of the saved blobs is saved there, but not the actual contents
	IndexPath string `json:"-"`
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
	indexFile, err := os.Create(s.IndexPath)
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
	indexFile, err := os.Open(s.IndexPath)
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
