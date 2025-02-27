package backup

type BufFile struct {
	// ID is the MD5 hash of the file's data. The file name is derived from this
	// TODO: maybe string?
	ID []byte
}

// TODO: Constructor
// TODO: Get path
// TODO: Get MD5
// TODO: Get size!

// TODO: Directory with such files. Maybe structured by first byte. Should have methods to either
// download or reference blobs. Also, stop-the-world garbage collection, perhaps? Or ref counts...
