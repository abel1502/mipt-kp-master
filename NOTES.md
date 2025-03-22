# Notes on blobs

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
    - As it turns out, when a block blob is uploaded in a single operation,
      it will not have any blocks associated with its contents. This has to be
      handled separately.
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
- Container-wide snapshots aren't provided by Azure. We have to resort to
  making snapshots of individual blobs. The problems with this are:
  - If someone updates one of the blobs while we're in the process of creating
    the online snapshot, we cannot attribute a single timestamp to it. We may
    either ignore this issue or retry in such cases.
  - If someone deletes a blob after we've made an online snapshot but before
    we're done downloading everything, the blob's snapshots are automatically
    deleted alongside the original. We may, in theory, lock blobs for the
    duration of the backup, or we may also simply ignore this edge case, since
    a blob which was immediately deleted afterwards probably wasn't a valuable
    addition to the backup in the first place. Note that snapshots cannot be
    locked (leased).

# TODO

- Actually, with the above implemented, we might not need to store fragments
  by pointer anymore; That used to be a deduplication feature, but on second
  thought, the filesystem and MD5 hashes works much better for that purpose
- Write the text in Typst in this repository (compile locally), so that we have
  a diff of it.
- The next step would be to deduplicate beyond azure-level fragments.
- Do we need to keep previous snapshots?

- We minimize both storage and network traffic. First priority is storage,
  network is secondary. The "true" goal is optimizing costs (will make a nice
  point in the diploma). To analyze costs, maybe take not only azure, but
  another cloud provider as well. Ideally would be to come up with a formula
  with coefficients and minimize it smartly. Also point out that azure is just
  an example, and the same works in other systems.
- Metadata is also an overhead
- Client might want to defend against their own employees (malice, accidents)
- Work on deduplication
  - Mentor will send a book on the subject
  - Obvious idea: slice into chunks, compute hashes
  - Improvement: instead of fixed-size chunks, use an algebraic expression
    (polynome modulo / window hash), and slice when hits zero. "Rabin Codes",
    "Rabin-Carp rolling hash", "Rabin fingerprinting". May ensure an
    approximate chunk size in advance.
  - Window size is independent from chunk size.
  - Compute hash in the window, cut when hash starts with a given number of
    zero bits. In practice, common hashes are computed much faster than those
    suggested in Rabin-Carp stuff.
  - Benchmark different hashes (or find a ready benchmark).
  - Choice of chunk size isn't trivial (tradeoff between metadata overhead
    and deduplication efficiency).
  - How to combine deduplication and encryption? Azure (and other clouds)
    support asymmetric encryption natively. While making a backup, we don't
    have access to raw data (?), but deduplication is inefficient. Can we
    come up with some encryption that would be compatible with deduplication?
  - Compression --- development of the chunk size question. Chunked data may
    have a negative effect on compression. It's not trivial to pick algorithms
    that wouldn't interfere with each other.
  - Collisions.
  - Also smaller problems about adapting the techniques to the cloud.
- Just understanding the API isn't a scientific achievement yet.
- Put it that the specific vendor is just a sample (solid, popular, good api,
  etc.), but the same achievements are applicable to other cloud providers.
- Try to come up with an abstract cost formula for later optimization.
  (Dependency on the chunk size in dedup, for instance).
- Current results: implemented a working prototype for the backup tool.
- Maybe also plot something related to the backups as another result.
- No need to show code, discuss details or Azure API in particular.
- Understanding prior work is key to answering questions.
- Three points: cost analysis; traditional (alogrithm optimizations) and
  safety (combination with encryption). Should be a lot of articles in all
  of them

- Forget about conference, thankfully!
- Write a rough summary of technologies in the articles, rough language is
  alright. Around 1.5 - 2 weeks deadline (until 05.04.2025).
- Make a typst project in this repo. Move todo list there as well?
