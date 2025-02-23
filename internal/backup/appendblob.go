package backup

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

func (*AppendBlob) Type() BlobType {
	return BlobTypeAppend
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
