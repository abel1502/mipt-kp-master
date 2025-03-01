package backup

import (
	"context"

	azblob "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	azcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

type AppendBlob struct {
	CommonBlob
	// Fragments is the single-linked list of this blob's fragments
	Fragments *AppendBlobFragment
}

type AppendBlobFragment struct {
	// LastChunk is the last chunk of this blob
	LastChunk []byte
	// Previous is the preceding fragment
	Previous *AppendBlobFragment
}

func DownloadAppendBlob(
	ctx context.Context,
	contClient *azcontainer.Client,
	name string,
	snapshot string,
	prev *AppendBlob,
) (*AppendBlob, error) {
	client, err := contClient.NewAppendBlobClient(name).WithSnapshot(snapshot)
	if err != nil {
		return nil, err
	}

	common, err := downloadCommon(ctx, client.BlobClient(), name)
	if err != nil {
		return nil, err
	}

	offset := uint64(0)
	size := common.ContentSize

	if prev != nil {
		offset = prev.Common().ContentSize
		size -= offset
	}

	blob := &AppendBlob{
		CommonBlob: *common,
		Fragments: &AppendBlobFragment{
			LastChunk: make([]byte, size),
			Previous:  nil,
		},
	}

	_, err = client.DownloadBuffer(ctx, blob.Fragments.LastChunk, &azblob.DownloadBufferOptions{
		Range: azblob.HTTPRange{
			Offset: int64(offset),
			Count:  int64(size),
		},
	})
	if err != nil {
		return nil, err
	}

	if prev != nil {
		blob.Fragments.Previous = prev.Fragments
	}

	return blob, nil
}

func (*AppendBlob) Type() azcontainer.BlobType {
	return azcontainer.BlobTypeAppendBlob
}

func (a *AppendBlob) Common() *CommonBlob {
	return &a.CommonBlob
}

func (a *AppendBlob) ShallowClone() Blob {
	return &AppendBlob{
		CommonBlob: a.CommonBlob,
		Fragments:  a.Fragments,
	}
}

var _ Blob = (*AppendBlob)(nil)
