package azure

import (
	"context"
	"iter"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

// IterBlobs provides an iterator over all blobs in a container
// TODO: private?
func IterBlobs(ctx context.Context, client *container.Client) iter.Seq2[*container.BlobItem, error] {
	return func(yield func(*container.BlobItem, error) bool) {
		blobPager := client.NewListBlobsFlatPager(nil)

		for blobPager.More() {
			blobPage, err := blobPager.NextPage(ctx)
			if err != nil {
				yield(nil, err)
				return
			}

			for _, blob := range blobPage.Segment.BlobItems {
				if !yield(blob, nil) {
					return
				}
			}
		}
	}
}

// ContainerSnapshot encapsulates snapshots of all blobs in a container
type ContainerSnapshot struct {
	Client  *container.Client
	TakenAt time.Time
	Blobs   []BlobInfo
}

// BlobInfo stores the information sufficient to reference a blob snapshot within a known container
type BlobInfo struct {
	Name         string
	Snapshot     string
	LastModified time.Time // TODO: May be checked against TakenAt to make sure the whole snapshot was atomic
}

// TakeSnapshot takes snapshots of all blobs in a container
func TakeSnapshot(ctx context.Context, client *container.Client) (*ContainerSnapshot, error) {
	success := false

	result := &ContainerSnapshot{
		Client:  client,
		TakenAt: time.Now(),
		Blobs:   nil,
	}

	// Clean up the snapshots if an error occurs
	defer func() {
		if !success {
			result.Delete(ctx)
		}
	}()

	for blob, err := range IterBlobs(ctx, client) {
		if err != nil {
			return nil, err
		}

		if blob.Snapshot != nil || *blob.Deleted {
			continue
		}

		blobClient := client.NewBlobClient(*blob.Name)
		snapshotResp, err := blobClient.CreateSnapshot(ctx, nil)
		if err != nil {
			return nil, err
		}

		result.Blobs = append(result.Blobs, BlobInfo{
			Name:         *blob.Name,
			Snapshot:     *snapshotResp.Snapshot,
			LastModified: *snapshotResp.LastModified,
		})
	}

	success = true
	return result, nil
}

// Delete cleans up the snapshots from the server
func (c *ContainerSnapshot) Delete(ctx context.Context) {
	for _, blob := range c.Blobs {
		snapshotClient, err := c.Client.NewBlobClient(blob.Name).WithSnapshot(blob.Snapshot)
		if err != nil {
			continue
		}
		_, _ = snapshotClient.Delete(ctx, nil)
	}
}

func OpenClient(containerURL string) (*container.Client, error) {
	defaultCred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	containerClient, err := container.NewClient(
		containerURL,
		defaultCred,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return containerClient, nil
}
