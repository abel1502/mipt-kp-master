package backup

import (
	"context"

	azcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

type BlockBlob struct {
	CommonBlob
	// Fragments is the list of blocks that make up the blob
	Fragments []*BlockBlobFragment
}

type BlockBlobFragment struct {
	// ID is the base64-encoded block ID
	ID string
	// Content is the block data
	Content []byte
}

func DownloadBlockBlob(
	ctx context.Context,
	contClient *azcontainer.Client,
	name string,
	snapshot string,
	prev *BlockBlob,
) (*BlockBlob, error) {
	client, err := contClient.NewBlockBlobClient(name).WithSnapshot(snapshot)
	if err != nil {
		return nil, err
	}

	// TODO
	_ = client
	panic("not implemented")

	return nil, nil
}

func (*BlockBlob) Type() azcontainer.BlobType {
	return azcontainer.BlobTypeBlockBlob
}

func (b *BlockBlob) Common() *CommonBlob {
	return &b.CommonBlob
}

func (b *BlockBlob) ShallowClone() Blob {
	return &BlockBlob{
		CommonBlob: b.CommonBlob,
		Fragments:  b.Fragments,
	}
}

var _ Blob = (*BlockBlob)(nil)
