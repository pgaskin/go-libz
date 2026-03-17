package libz

import (
	"encoding/binary"
	"math"
)

// Compress compresses src into dst, returning a slice pointing to the
// compressed data. If dst is nil, the maximum possible compressed length is
// allocated. Otherwise, if dst is too short, an error is returned.
func Compress(dst, src []byte, level Level) ([]byte, error) {
	if dst == nil {
		n, err := compressBound(len(src))
		if err != nil {
			return nil, err
		}
		dst = make([]byte, n)
	}

	var z *libz
	if len(dst)+len(src) < 128*1024 { // don't use the pool if it'll be allocating a large amount of memory (TODO: make configurable?)
		z = pool.Get().(*libz)
		defer pool.Put(z)
	} else {
		z = instantiate()
	}

	if uint64(len(src)+len(dst)) > math.MaxUint32 {
		return nil, Z_MEM_ERROR
	}
	if len(dst) == 0 {
		return nil, Z_BUF_ERROR
	}

	ptr, err := z.malloc(4 + len(src) + len(dst))
	if err != nil {
		return nil, err
	}
	defer z.free(ptr)

	binary.LittleEndian.PutUint32((*z.mod.Xmemory().Slice())[uint32(ptr):], uint32(len(dst)))
	copy((*z.mod.Xmemory().Slice())[uint32(ptr)+4:], src)

	if err := toError(z.mod.Xcompress2(int32(ptr+4+uint32(len(src))), int32(ptr), int32(ptr+4), int32(uint32(len(src))), int32(level))); err != nil {
		return nil, err
	}

	n := binary.LittleEndian.Uint32((*z.mod.Xmemory().Slice())[uint32(ptr):])
	copy(dst[:n], (*z.mod.Xmemory().Slice())[uint32(ptr)+4+uint32(len(src)):])

	return dst[:n], nil
}

func compressBound(n int) (int, error) {
	if uint64(n) > math.MaxUint32 {
		return 0, Z_MEM_ERROR
	}

	z := pool.Get().(*libz)
	defer pool.Put(z)

	res := z.mod.XcompressBound(int32(uint32(n)))
	return int(res), nil // sizeof(unsigned long) == 4
}

// Uncompress decompresses src into dst, returning a slice pointing to the
// decompressed data. If dst is too short, an error is returned.
func Uncompress(dst, src []byte) ([]byte, error) {
	var z *libz
	if len(dst)+len(src) < 128*1024 { // don't use the pool if it'll be allocating a large amount of memory (TODO: make configurable?)
		z = pool.Get().(*libz)
		defer pool.Put(z)
	} else {
		z = instantiate()
	}

	if uint64(len(src)+len(dst)) > math.MaxUint32 {
		return nil, Z_MEM_ERROR
	}

	ptr, err := z.malloc(4 + len(src) + len(dst))
	if err != nil {
		return nil, err
	}
	defer z.free(ptr)

	binary.LittleEndian.PutUint32((*z.mod.Xmemory().Slice())[uint32(ptr):], uint32(len(dst)))
	copy((*z.mod.Xmemory().Slice())[uint32(ptr)+4:], src)

	if err := toError(z.mod.Xuncompress(int32(ptr+4+uint32(len(src))), int32(ptr), int32(ptr+4), int32(uint32(len(src))))); err != nil {
		return nil, err
	}

	n := binary.LittleEndian.Uint32((*z.mod.Xmemory().Slice())[uint32(ptr):])
	copy(dst[:n], (*z.mod.Xmemory().Slice())[uint32(ptr)+4+uint32(len(src)):])

	return dst[:n], nil
}

// TODO: inflate/deflate reader/writer wrappers
