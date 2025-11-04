package libz

import (
	"context"
	"errors"
	"runtime"
	"sync"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// Options for initializing wasm.
//
// To use the default build:
//
//	import _ "github.com/pgaskin/go-libz/embed"
var (
	Binary        []byte
	RuntimeConfig wazero.RuntimeConfig
)

var instance struct {
	runtime  wazero.Runtime
	compiled wazero.CompiledModule
	err      error
	once     sync.Once
}

// Initialize compiles the wasm binary.
//
// This is called implicitly when used for the first time.
func Initialize() error {
	instance.once.Do(initialize)
	return instance.err
}

var (
	errNoBinary      = errors.New("libz: no wasm binary available")
	errMissingSymbol = errors.New("libz: wasm binary missing required symbol")
	errBadReturn     = errors.New("libz: bad return value from wasm func")
)

func initialize() {
	ctx := context.Background()

	cfg := RuntimeConfig
	if cfg == nil {
		cfg = wazero.NewRuntimeConfig()
	}
	cfg = cfg.WithCoreFeatures(api.CoreFeaturesV2)

	instance.runtime = wazero.NewRuntimeWithConfig(ctx, cfg)

	bin := Binary
	if bin == nil {
		instance.err = errNoBinary
		return
	}

	instance.compiled, instance.err = instance.runtime.CompileModule(ctx, bin)
}

type libz struct {
	ctx context.Context
	mod api.Module
}

// TODO: cache exported functions
// TODO: reuse stack slice

// pool is a pool of instantiated sqlite modules.
var pool sync.Pool

// poolInstantiate gets a instantiated module from [presspool], or instantiates
// a new one.
func poolInstantiate() (*libz, error) {
	if x := pool.Get(); x != nil {
		return x.(*libz), nil
	}
	return instantiate()
}

func instantiate() (*libz, error) {
	if err := Initialize(); err != nil {
		return nil, err
	}

	z := new(libz)
	z.ctx = context.Background()

	// we don't need wasi_snapshot_preview1 since we optimized out the gzip stuff which requires filesystem functions

	var err error
	z.mod, err = instance.runtime.InstantiateModule(z.ctx, instance.compiled, wazero.NewModuleConfig().WithName(""))
	if err != nil {
		return nil, err
	}
	if z.getfn("zError") == nil {
		return nil, errMissingSymbol
	}
	if z.getfn("malloc") == nil {
		return nil, errMissingSymbol
	}
	if z.getfn("free") == nil {
		return nil, errMissingSymbol
	}
	runtime.SetFinalizer(z, func(z *libz) {
		z.mod.Close(z.ctx)
	})
	return z, nil
}

func (z *libz) getfn(name string) api.Function {
	return z.mod.ExportedFunction(name)
}

func (z *libz) malloc(n int) (uint64, error) {
	malloc := z.getfn("malloc")
	if malloc == nil {
		return 0, errMissingSymbol
	}
	res, err := malloc.Call(z.ctx, uint64(n))
	if err != nil {
		return 0, err
	}
	if len(res) != 1 {
		return 0, errBadReturn
	}
	if res[0] == 0 {
		return 0, Z_MEM_ERROR
	}
	return res[0], nil
}

func (z *libz) free(ptr uint64) error {
	free := z.getfn("free")
	if free == nil {
		return errMissingSymbol
	}
	free.Call(z.ctx, ptr)
	return nil
}

func toError(res []uint64) error {
	if len(res) != 1 {
		return errBadReturn
	}
	x := uint32(res[0])
	r := ErrorCode(x)
	if r == Z_OK {
		return nil
	}
	return &Error{rc: r}
}
