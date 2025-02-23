package backup

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

func (*PageBlob) Type() BlobType {
	return BlobTypePage
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
