package libz

import (
	"errors"
	"math"
)

// Compress compresses src into dst, returning a slice pointing to the
// compressed data. If dst is nil, the maximum possible compressed length is
// allocated. Otherwise, if dst is too short, an error is returned.
// TODO: func Compress(dst, src []byte, level int) ([]byte, error)

// Uncompress decompresses src into dst, returning a slice pointing to the
// decompressed data. If dst is too short, an error is returned.
func Uncompress(dst, src []byte) ([]byte, error) {
	var z *libz
	if len(dst)+len(src) < 128*1024 { // don't use the pool if it'll be allocating a large amount of memory (TODO: make configurable?)
		var err error
		z, err = poolInstantiate()
		if err != nil {
			return nil, err
		}
		defer pool.Put(z)
	} else {
		var err error
		z, err = instantiate()
		if err != nil {
			return nil, err
		}
	}

	uncompress := z.getfn("uncompress")
	if uncompress == nil {
		return nil, errors.New("libz: bad wasm binary: missing uncompress")
	}

	if len(src)+len(dst) > math.MaxUint32 {
		return nil, Z_MEM_ERROR
	}

	ptr, err := z.malloc(4 + len(src) + len(dst))
	if err != nil {
		return nil, err
	}
	defer z.free(ptr)

	z.mod.Memory().WriteUint32Le(uint32(ptr), uint32(len(dst)))
	z.mod.Memory().Write(uint32(ptr+4), src)

	res, err := uncompress.Call(z.ctx, ptr+4+uint64(len(src)), ptr, ptr+4, uint64(len(src)))
	if err != nil {
		return nil, err
	}
	if err := toError(res); err != nil {
		return nil, err
	}

	n, ok := z.mod.Memory().ReadUint32Le(uint32(ptr))
	if !ok || n > uint32(len(dst)) {
		panic("wtf")
	}

	v, ok := z.mod.Memory().Read(uint32(ptr)+4+uint32(len(src)), n)
	if !ok {
		panic("wtf")
	}
	copy(dst, v)

	return dst[:n], nil
}

// TODO: inflate/deflate reader/writer wrappers
