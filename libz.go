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
