## Benchmarks comparing various hash functions that can be used in CDC
(Content-defined chunking)

- Rabin hash
- Adler32
- Cyclic hash
- MD5
- SHA256
- HighwayHash


Note: actually, hash-specific hardware acceleration appears to be quite rare.
The most I could find is the SHA instruction set for Intel, and it was only
introduced in 2024 (so not available on my machine). Vectored instructions
also aren't particularly useful since chunking needs to operate in 1-byte
increments. I suspect the performance gain in established hashes might not
provide a meaningful impact here.

