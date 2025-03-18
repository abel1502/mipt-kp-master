package backup

import (
	"context"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	azcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/abel1502/mipt-kp-m-test/internal/azure"
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
	Content *FileBuf
}

func DownloadBlockBlob(
	ctx context.Context,
	repo *Repository,
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

	if len(blockList.CommittedBlocks) == 0 && commonBlob.ContentSize != 0 {
		// This is possible if the blob was uploaded in a single request.
		// It is roughly equivalent to the blob consisting of a single block.

		blockList.CommittedBlocks = []*blockblob.Block{{
			Name: azure.Addressof(""),
			Size: azure.Addressof(int64(commonBlob.ContentSize)),
		}}
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
			fb, err := repo.DownloadBlobRangeAsFileBuf(ctx, client.BlobClient(), offset, uint64(*block.Size))
			if err != nil {
				return nil, err
			}

			fragment = &BlockBlobFragment{
				ID:      *block.Name,
				Content: fb,
			}
		}

		blob.Fragments = append(blob.Fragments, fragment)

		offset += fragment.Content.Size
	}

	return blob, nil
}

func (*BlockBlob) Type() azcontainer.BlobType {
	return azcontainer.BlobTypeBlockBlob
}

func (b *BlockBlob) Common() *CommonBlob {
	return &b.CommonBlob
}

func (b *BlockBlob) Export(ctx context.Context, repo *Repository) io.ReadCloser {
	fragments := make([]io.ReadCloser, 0, len(b.Fragments))

	for _, fragment := range b.Fragments {
		fragments = append(fragments, fragment.Content.LazyReader(repo.LocalPath))
	}

	return ChainReader(fragments...)
}

func (b *BlockBlob) ShallowClone() Blob {
	return &BlockBlob{
		CommonBlob: b.CommonBlob,
		Fragments:  b.Fragments,
	}
}

var _ Blob = (*BlockBlob)(nil)
