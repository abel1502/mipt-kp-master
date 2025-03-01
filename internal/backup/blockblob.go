package backup

import (
	"context"

	azblob "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
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

	blockList, err := client.GetBlockList(ctx, blockblob.BlockListTypeCommitted, nil)
	if err != nil {
		return nil, err
	}

	commonBlob, err := downloadCommon(ctx, client.BlobClient(), name)
	if err != nil {
		return nil, err
	}

	blob := &BlockBlob{
		CommonBlob: *commonBlob,
		Fragments:  make([]*BlockBlobFragment, 0, len(blockList.CommittedBlocks)),
	}

	knownFragments := make(map[string]*BlockBlobFragment)
	if prev != nil {
		// Note: fragments from any preceding revisions of the blob do not matter
		// in this context, as they weren't accessible to the update operations for this blob
		for _, fragment := range prev.Fragments {
			knownFragments[fragment.ID] = fragment
		}
	}

	offset := uint64(0)

	for _, block := range blockList.CommittedBlocks {
		fragment, ok := knownFragments[*block.Name]
		if !ok {
			fragment = &BlockBlobFragment{
				ID:      *block.Name,
				Content: make([]byte, *block.Size),
			}

			_, err := client.DownloadBuffer(ctx, fragment.Content, &azblob.DownloadBufferOptions{
				Range: azblob.HTTPRange{
					Offset: int64(offset),
					Count:  int64(len(fragment.Content)),
				},
			})
			if err != nil {
				return nil, err
			}
		}

		blob.Fragments = append(blob.Fragments, fragment)

		offset += uint64(len(fragment.Content))
	}

	return blob, nil
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
