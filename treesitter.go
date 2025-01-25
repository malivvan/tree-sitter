package sitter

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

const Version = "v0.24.7"

//go:embed lib/ts.wasm
var tsWasm []byte

type TreeSitter struct {
	m api.Module

	malloc api.Function
	free   api.Function
	strlen api.Function

	parserNew         api.Function
	parserParseString api.Function
	parserDelete      api.Function
	parserSetLanguage api.Function

	//languageName    api.Function
	languageVersion api.Function

	treeRootNode api.Function

	queryNew              api.Function
	queryCursorNew        api.Function
	queryCusorExec        api.Function
	queryCursorNextMatch  api.Function
	queryCaptureNameForID api.Function

	nodeString          api.Function
	nodeChildCount      api.Function
	nodeNamedChildCount api.Function
	nodeChild           api.Function
	nodeNamedChild      api.Function
	nodeType            api.Function
	nodeEndByte         api.Function
	nodeStartByte       api.Function
	nodeIsError         api.Function

	languageC   api.Function
	languageCpp api.Function
}

func New(ctx context.Context) (TreeSitter, error) {
	r := wazero.NewRuntime(ctx)

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	compiled, err := r.CompileModule(ctx, tsWasm)
	if err != nil {
		return TreeSitter{}, fmt.Errorf("compiling wasm module: %w", err)
	}

	//env := r.NewHostModuleBuilder("env")

	//_, err = env.Instantiate(ctx)
	//if err != nil {
	//	return TreeSitter{}, fmt.Errorf("instantiating emscripten module: %w", err)
	//}

	mod, err := r.InstantiateModule(ctx, compiled, wazero.NewModuleConfig())
	if err != nil {
		return TreeSitter{}, fmt.Errorf("instantiating module: %w", err)
	}

	return TreeSitter{
		m:                     mod,
		malloc:                mod.ExportedFunction("malloc"),
		free:                  mod.ExportedFunction("free"),
		strlen:                mod.ExportedFunction("strlen"),
		parserNew:             mod.ExportedFunction("ts_parser_new"),
		parserParseString:     mod.ExportedFunction("ts_parser_parse_string"),
		parserSetLanguage:     mod.ExportedFunction("ts_parser_set_language"),
		parserDelete:          mod.ExportedFunction("ts_parser_delete"),
		queryNew:              mod.ExportedFunction("ts_query_new"),
		queryCursorNew:        mod.ExportedFunction("ts_query_cursor_new"),
		queryCusorExec:        mod.ExportedFunction("ts_query_cursor_exec"),
		queryCursorNextMatch:  mod.ExportedFunction("ts_query_cursor_next_match"),
		queryCaptureNameForID: mod.ExportedFunction("ts_query_capture_name_for_id"),
		//	languageName:          mod.ExportedFunction("ts_language_name"),
		languageVersion:     mod.ExportedFunction("ts_language_version"),
		treeRootNode:        mod.ExportedFunction("ts_tree_root_node"),
		nodeString:          mod.ExportedFunction("ts_node_string"),
		nodeChildCount:      mod.ExportedFunction("ts_node_child_count"),
		nodeNamedChildCount: mod.ExportedFunction("ts_node_named_child_count"),
		nodeChild:           mod.ExportedFunction("ts_node_child"),
		nodeNamedChild:      mod.ExportedFunction("ts_node_named_child"),
		nodeType:            mod.ExportedFunction("ts_node_type"),
		nodeStartByte:       mod.ExportedFunction("ts_node_start_byte"),
		nodeEndByte:         mod.ExportedFunction("ts_node_end_byte"),
		nodeIsError:         mod.ExportedFunction("ts_node_is_error"),
		languageC:           mod.ExportedFunction("tree_sitter_c"),
		languageCpp:         mod.ExportedFunction("tree_sitter_cpp"),
	}, nil
}

func (t TreeSitter) allocateString(
	ctx context.Context,
	str string,
) (ptr uint64, size uint64, free func(), err error) {
	strByte := []byte(str)
	strSize := uint64(len(strByte))
	strPtr, err := t.malloc.Call(ctx, strSize)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("allocating string: %w", err)
	}

	if !t.m.Memory().Write(uint32(strPtr[0]), strByte) {
		return 0, 0, nil, fmt.Errorf("writing string: %w", err)
	}

	return strPtr[0], strSize, func() {
		t.free.Call(context.Background(), strPtr[0])
	}, nil
}

func (t TreeSitter) readString(ctx context.Context, ptr uint64) (string, error) {
	strSize, err := t.strlen.Call(ctx, ptr)
	if err != nil {
		return "", fmt.Errorf("getting string length: %w", err)
	}
	strBytes, ok := t.m.Memory().Read(uint32(ptr), uint32(strSize[0]))
	if !ok {
		return "", errors.New("error reading string")
	}
	return string(strBytes), nil
}
