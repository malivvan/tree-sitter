package sitter

import (
	"context"
	"crypto/rand"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"io/fs"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

const Version = "v0.24.7"

type Context struct {
	ctx context.Context
	mod api.Module
	vfs fs.FS
	out io.Writer
}

type TreeSitter struct {
	ctx  context.Context
	mod  api.Module
	vfs  fs.FS
	out  io.Writer
	fns  map[string]api.Function
	lang map[string]api.Function
}

func New(vfs fs.FS, out io.Writer) (*TreeSitter, error) {
	err := initialize()
	if err != nil {
		return nil, err
	}
	ts := new(TreeSitter)
	ts.ctx = context.Background()
	cfg := wazero.NewModuleConfig().
		WithName("").
		WithSysNanosleep().
		WithSysNanotime().
		WithSysWalltime().
		WithRandSource(rand.Reader)
	if vfs != nil {
		cfg = cfg.WithFS(vfs)
	}
	if out != nil {
		cfg = cfg.WithStdout(out).WithStderr(out)
	}
	ts.mod, err = instance.runtime.InstantiateModule(ts.ctx, instance.compiled, cfg)
	if err != nil {
		return nil, err
	}
	ts.fns = make(map[string]api.Function)
	for _, name := range _functions {
		if fn := ts.mod.ExportedFunction(name); fn != nil {
			ts.fns[name] = fn
		} else {
			return nil, fmt.Errorf("missing function %q", name)
		}
	}
	ts.lang = make(map[string]api.Function)
	for _, name := range _languages {
		if fn := ts.mod.ExportedFunction("tree_sitter_" + name); fn != nil {
			ts.lang[name] = fn
		} else {
			return nil, fmt.Errorf("missing language %q", name)
		}
	}
	return ts, nil
}

func (ts *TreeSitter) call(name string, args ...uint64) ([]uint64, error) {
	return ts.callCtx(ts.ctx, name, args...)
}

func (ts *TreeSitter) callCtx(ctx context.Context, name string, args ...uint64) ([]uint64, error) {
	return ts.fns[name].Call(ctx, args...)
}

func (ts *TreeSitter) allocateString(str string) (ptr uint64, size uint64, free func(), err error) {
	strByte := []byte(str)
	strSize := uint64(len(strByte))
	strPtr, err := ts.call(_malloc, strSize)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("allocating string: %w", err)
	}

	if !ts.mod.Memory().Write(uint32(strPtr[0]), strByte) {
		return 0, 0, nil, fmt.Errorf("writing string: %w", err)
	}

	return strPtr[0], strSize, func() {
		ts.callCtx(context.Background(), _free, strPtr[0])
	}, nil
}

func (ts *TreeSitter) readString(ptr uint64) (string, error) {
	strSize, err := ts.call(_strlen, ptr)
	if err != nil {
		return "", fmt.Errorf("getting string length: %w", err)
	}
	strBytes, ok := ts.mod.Memory().Read(uint32(ptr), uint32(strSize[0]))
	if !ok {
		return "", errors.New("error reading string")
	}
	return string(strBytes), nil
}
