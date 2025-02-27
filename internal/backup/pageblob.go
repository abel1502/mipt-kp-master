package backup

import (
	"context"
	"io"
	"slices"

	azblob "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	azcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/pageblob"
)

type PageBlob struct {
	CommonBlob
	// Fragments is the list of this blob's pages
	Fragments []*PageBlobFragment
}

// TODO: Support clearing parts in the middle of a page range?
// We don't always want to keep the entire history, though.
// For that scenario, I'll need
type PageBlobFragment struct {
	// Offset is the fragment offset (512-bytes-aligned)
	Offset uint64
	// Content is the fragment data (512-bytes-aligned in size)
	Content []byte
	// ContentMD5 is the MD5 hash of the fragment
	ContentMD5 []byte
}

func DownloadPageBlob(
	ctx context.Context,
	contClient *azcontainer.Client,
	name string,
	snapshot string,
	prev *PageBlob,
) (*PageBlob, error) {
	client, err := contClient.NewPageBlobClient(name).WithSnapshot(snapshot)
	if err != nil {
		return nil, err
	}

	pages, err := listPages(ctx, client)
	if err != nil {
		return nil, err
	}

	commonBlob, err := downloadCommon(ctx, client.BlobClient())
	if err != nil {
		return nil, err
	}

	blob := &PageBlob{
		CommonBlob: *commonBlob,
		Fragments:  make([]*PageBlobFragment, 0, len(pages)),
	}

	knownFragments := make(map[uint64]*PageBlobFragment)
	if prev != nil {
		// TODO: unlike blob blocks, here we might be interested in fragments from previous
		// versions of the blob. Also, perhaps look up based on MD5 instead of the offset?
		for _, fragment := range prev.Fragments {
			knownFragments[fragment.Offset] = fragment
		}
	}

	for _, page := range pages {
		stream, err := client.DownloadStream(ctx, &azblob.DownloadStreamOptions{
			Range: azblob.HTTPRange{
				Offset: int64(page.Offset),
				Count:  int64(page.Size),
			},
			// TODO: This, apparently, only works for <4MB ranges!
			// Why do they even need a *bool?
			RangeGetContentMD5: func(b bool) *bool { return &b }(true),
		})
		if err != nil {
			return nil, err
		}
		defer stream.Body.Close()

		fragment, ok := knownFragments[page.Offset]
		ok = ok && len(fragment.Content) == int(*stream.ContentLength)
		ok = ok && slices.Equal(fragment.ContentMD5, stream.ContentMD5)
		if !ok {
			fragment = &PageBlobFragment{
				Offset:     page.Offset,
				Content:    make([]byte, *stream.ContentLength),
				ContentMD5: stream.ContentMD5,
			}

			_, err := io.ReadFull(stream.Body, fragment.Content)
			if err != nil {
				return nil, err
			}
		}

		blob.Fragments = append(blob.Fragments, fragment)
	}

	return blob, nil
}

type pageInfo struct {
	Offset uint64
	Size   uint64
}

func listPages(ctx context.Context, client *pageblob.Client) ([]pageInfo, error) {
	pagePager := client.NewGetPageRangesPager(nil)
	result := make([]pageInfo, 0, 8)

	for pagePager.More() {
		pagePage, err := pagePager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, page := range pagePage.PageRange {
			result = append(result, pageInfo{
				Offset: uint64(*page.Start),
				Size:   uint64(*page.End) - uint64(*page.Start),
			})
		}
	}

	return result, nil
}

func (*PageBlob) Type() azcontainer.BlobType {
	return azcontainer.BlobTypePageBlob
}

func (p *PageBlob) Common() *CommonBlob {
	return &p.CommonBlob
}

func (p *PageBlob) ShallowClone() Blob {
	return &PageBlob{
		CommonBlob: p.CommonBlob,
		Fragments:  p.Fragments,
	}
}

var _ Blob = (*PageBlob)(nil)
