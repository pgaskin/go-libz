package libz

import (
	"context"
	"errors"
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
		instance.err = errors.New("libz: no wasm binary available")
		return
	}

	instance.compiled, instance.err = instance.runtime.CompileModule(ctx, bin)
}

type libz struct {
	ctx context.Context
	mod api.Module
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
		return nil, errors.New("libz: bad wasm binary")
	}
	return z, nil
}

func (z *libz) getfn(name string) api.Function {
	return z.mod.ExportedFunction(name)
}

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
