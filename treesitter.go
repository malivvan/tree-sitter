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
	//
	//malloc api.Function
	//free   api.Function
	//strlen api.Function
	//
	//parserNew         api.Function
	//parserParseString api.Function
	//parserDelete      api.Function
	//parserSetLanguage api.Function
	//
	////languageName    api.Function
	//languageVersion api.Function
	//
	//treeRootNode api.Function
	//
	//queryNew              api.Function
	//queryCursorNew        api.Function
	//queryCusorExec        api.Function
	//queryCursorNextMatch  api.Function
	//queryCaptureNameForID api.Function
	//
	//nodeString          api.Function
	//nodeChildCount      api.Function
	//nodeNamedChildCount api.Function
	//nodeChild           api.Function
	//nodeNamedChild      api.Function
	//nodeType            api.Function
	//nodeEndByte         api.Function
	//nodeStartByte       api.Function
	//nodeIsError         api.Function
	//
	//languageC   api.Function
	//languageCpp api.Function
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
	//
	//return &TreeSitter{
	//	mod:                   ts.mod,
	//	malloc:                ts.mod.ExportedFunction("malloc"),
	//	free:                  ts.mod.ExportedFunction("free"),
	//	strlen:                ts.mod.ExportedFunction("strlen"),
	//	parserNew:             ts.mod.ExportedFunction("ts_parser_new"),
	//	parserParseString:     ts.mod.ExportedFunction("ts_parser_parse_string"),
	//	parserSetLanguage:     ts.mod.ExportedFunction("ts_parser_set_language"),
	//	parserDelete:          ts.mod.ExportedFunction("ts_parser_delete"),
	//	queryNew:              ts.mod.ExportedFunction("ts_query_new"),
	//	queryCursorNew:        ts.mod.ExportedFunction("ts_query_cursor_new"),
	//	queryCusorExec:        ts.mod.ExportedFunction("ts_query_cursor_exec"),
	//	queryCursorNextMatch:  ts.mod.ExportedFunction("ts_query_cursor_next_match"),
	//	queryCaptureNameForID: ts.mod.ExportedFunction("ts_query_capture_name_for_id"),
	//	//	languageName:          ts.mod.ExportedFunction("ts_language_name"),
	//	languageVersion:     ts.mod.ExportedFunction("ts_language_version"),
	//	treeRootNode:        ts.mod.ExportedFunction("ts_tree_root_node"),
	//	nodeString:          ts.mod.ExportedFunction("ts_node_string"),
	//	nodeChildCount:      ts.mod.ExportedFunction("ts_node_child_count"),
	//	nodeNamedChildCount: ts.mod.ExportedFunction("ts_node_named_child_count"),
	//	nodeChild:           ts.mod.ExportedFunction("ts_node_child"),
	//	nodeNamedChild:      ts.mod.ExportedFunction("ts_node_named_child"),
	//	nodeType:            ts.mod.ExportedFunction("ts_node_type"),
	//	nodeStartByte:       ts.mod.ExportedFunction("ts_node_start_byte"),
	//	nodeEndByte:         ts.mod.ExportedFunction("ts_node_end_byte"),
	//	nodeIsError:         ts.mod.ExportedFunction("ts_node_is_error"),
	//	languageC:           ts.mod.ExportedFunction("tree_sitter_c"),
	//	languageCpp:         ts.mod.ExportedFunction("tree_sitter_cpp"),
	//}, nil
	return ts, nil
}

func (ts *TreeSitter) call(name string, args ...uint64) ([]uint64, error) {
	return ts.callCtx(ts.ctx, name, args...)
}

func (ts *TreeSitter) callCtx(ctx context.Context, name string, args ...uint64) ([]uint64, error) {
	return ts.fns[name].Call(ctx, args...)
}

func (t *TreeSitter) allocateString(str string) (ptr uint64, size uint64, free func(), err error) {
	strByte := []byte(str)
	strSize := uint64(len(strByte))
	strPtr, err := t.call(_malloc, strSize)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("allocating string: %w", err)
	}

	if !t.mod.Memory().Write(uint32(strPtr[0]), strByte) {
		return 0, 0, nil, fmt.Errorf("writing string: %w", err)
	}

	return strPtr[0], strSize, func() {
		t.callCtx(context.Background(), _free, strPtr[0])
	}, nil
}

func (t *TreeSitter) readString(ptr uint64) (string, error) {
	strSize, err := t.call(_strlen, ptr)
	if err != nil {
		return "", fmt.Errorf("getting string length: %w", err)
	}
	strBytes, ok := t.mod.Memory().Read(uint32(ptr), uint32(strSize[0]))
	if !ok {
		return "", errors.New("error reading string")
	}
	return string(strBytes), nil
}
