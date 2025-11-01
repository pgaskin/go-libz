package libz

// Compress compresses src into dst, returning a slice pointing to the
// compressed data. If dst is nil, the maximum possible compressed length is
// allocated. Otherwise, if dst is too short, an error is returned.
// TODO: func Compress(dst, src []byte, level int) ([]byte, error)

// Decompress decompresses src into dst, returning a slice pointing to the
// decompressed data. If dst is too short, an error is returned.
// TODO: func Decompress(dst, src []byte) ([]byte, error)

// TODO: inflate/deflate reader/writer wrappers
