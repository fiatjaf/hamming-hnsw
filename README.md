# hamming-hnsw

Package `hamming-hnsw` implements Hierarchical Navigable Small World graphs in Go. You
can read up about how they work [here](https://www.youtube.com/watch?v=77QH0Y2PYKg).

It's a variant of https://github.com/coder/hnsw that uses _Hamming distance_ instead of
cosine distance. It can only be used with binary vectors, not normal float32 vectors.

Read https://github.com/coder/hnsw/blob/ff889c91944e4850627a4d81ff80d22f72355d2d/README.md
for the rest of the information you need to use this.
