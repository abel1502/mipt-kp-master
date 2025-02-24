package backup

import (
	"context"

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

	// TODO
	_ = client
	panic("not implemented")

	return nil, nil
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
