package sitter

import (
	"bytes"
	"context"
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"io"
	"os"
	"path/filepath"

	_ "embed"
	"math/bits"
	"runtime"
	"sync"
)

//go:embed lib/ts.wasm.br
var binaryBR []byte

var instance struct {
	runtime  wazero.Runtime
	compiled wazero.CompiledModule
	err      error
	once     sync.Once
}

func initialize() error {
	instance.once.Do(compileTreeSitter)
	return instance.err
}

func compileTreeSitter() {
	ctx := context.Background()
	var cfg wazero.RuntimeConfig
	if compilerSupported() {
		cfg = wazero.NewRuntimeConfigCompiler()
	} else {
		cfg = wazero.NewRuntimeConfigInterpreter()
	}
	if bits.UintSize < 64 {
		cfg = cfg.WithMemoryLimitPages(128) // 8MB
	} else {
		cfg = cfg.WithMemoryLimitPages(1024) // 64MB
	}
	cfg = cfg.WithCoreFeatures(api.CoreFeaturesV2)
	cacheDir, err := os.UserCacheDir()
	if err == nil {
		cacheDir = filepath.Join(cacheDir, "github.com", "malivvan", "tree-sitter")
		cache, err := wazero.NewCompilationCacheWithDir(cacheDir)
		if err == nil {
			cfg = cfg.WithCompilationCache(cache)
		}
	}

	instance.runtime = wazero.NewRuntimeWithConfig(ctx, cfg)
	_, instance.err = wasi_snapshot_preview1.Instantiate(ctx, instance.runtime)
	if instance.err != nil {
		return
	}
	if binaryBR == nil {
		instance.err = fmt.Errorf("tree-sitter wasm binary not found")
		return
	}
	binary, err := decompressBinary()
	if err != nil {
		instance.err = fmt.Errorf("decompressing tree-sitter wasm binary: %w", err)
		return
	}
	instance.compiled, instance.err = instance.runtime.CompileModule(ctx, binary)
}

func compilerSupported() bool {
	switch runtime.GOOS {
	case "linux", "android",
		"windows", "darwin",
		"freebsd", "netbsd", "dragonfly",
		"solaris", "illumos":
		break
	default:
		return false
	}
	switch runtime.GOARCH {
	case "amd64":
		return true
	case "arm64":
		return true
	default:
		return false
	}
}

func decompressBinary() ([]byte, error) {
	return io.ReadAll(brotli.NewReader(bytes.NewReader(binaryBR)))
}
