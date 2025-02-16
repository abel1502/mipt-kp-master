package azure

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"

type Container struct {
	client *container.Client
	// Note: container snapshots do not exist; this field simply means that every blob in the container is snapshotted
	isSnapshot bool
	blobItems  []container.BlobItem
}

// TODO: A constructor

func (c *Container) IsSnapshot() bool {
	return c.isSnapshot
}

func (c *Container) CreateSnapshot() (*Container, error) {
	// TODO
	return nil, nil
}
