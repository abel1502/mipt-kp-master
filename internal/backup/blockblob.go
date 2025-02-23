package backup

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

func (*BlockBlob) Type() BlobType {
	return BlobTypeBlock
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
