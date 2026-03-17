package libz

import (
	"math"
	"sync"

	"github.com/pgaskin/go-libz/internal/libz_wasm"
)

// TODO: cache exported functions
// TODO: reuse stack slice

var pool sync.Pool

type libz struct {
	mod *libz_wasm.Module
}

func instantiate() *libz {
	return &libz{
		mod: libz_wasm.New(),
	}
}

func init() {
	pool.New = func() any {
		return instantiate()
	}
}

func (z *libz) malloc(n int) (uint32, error) {
	if int64(n) >= math.MaxUint32 {
		return 0, Z_MEM_ERROR
	}
	res := z.mod.Xmalloc(int32(uint32(n)))
	if res == 0 {
		return 0, Z_MEM_ERROR
	}
	return uint32(res), nil
}

func (z *libz) free(ptr uint32) {
	if ptr != 0 {
		z.mod.Xfree(int32(ptr))
	}
}

func toError(res int32) error {
	r := ErrorCode(res)
	if r == Z_OK {
		return nil
	}
	return &Error{rc: r}
}
