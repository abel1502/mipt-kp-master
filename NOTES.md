- Current implementation needs to be reworked a lot, since I had many wrong
  assumptions about azure blob storage initially.
- We could use event grids to track extensively what changes occur to what
  files. Though it appears that just the modification timestamps should be
  enough.
- Instead of temporary snapshots, we could lease blocks, effectively preventing
  their updates for the duration of our operations on them.
- Do we back up metadata too? It is included in the last-modified timestamp,
  apparently.
- Blobs are still means to correspond to simple files, the only nuance is the
  support for more efficient update mechanisms. Any blob may be updated by
  simply overwriting its contents, but besides that:
  - Block blobs consist of named chunks; one can update new chunks independently
    and then rewrite the old blob with a new one, consisting of the specified
    chunks in the specified order. This is very friendly towards incremental
    backups.
  - Page blobs must be 512-bytes aligned in size, but they support random access
    rewrites of these 512-byte blocks. For incremental backups, we may,
    in theory, only store the information of what blocks have been added,
    updated and deleted (and with what data, unless deleted). Azure supports
    comparing page ranges against earlier snapshots, so it might be useful
    to maintain a server-side snapshot for these kinds of blocks.
  - Append blobs act as restricted block blobs, in that they only allow
    concatenating new data to their end. This is even more friendly towards
    incremental backups than conventional block blobs.
    The `x-ms-blob-committed-block-count` header can help us identify the added
    data, though we might just as use the blob size.
- A simple overwrite can always be identified using the creation date (right?).
- Blob type cannot be changed without overwriting the blob (?), so we can
  be certain that, if the blob creation date (?) stays the same, only a known
  set of operations could have been applied to it.
- We get MD5 and CRC32 checksums (on explicit demand?) for whole blocks and
  for block chunks when requesting to read them. This can help us identify
  and locate changes in large volumes of data very quickly, provided we can
  skip reading the full response body (or maybe even request that it isn't sent
  at all?).
- Many (all?) requests support the `If-Modified-Since` header, which may be
  useful to optimize some of the incremental backups. There's also ETag matching
- Blobs retain some information about being copied via Azure's inbuilt means.
  This may be useful for incremental backups.


