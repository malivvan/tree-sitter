package sitter

// malloc:                ts.mod.ExportedFunction("malloc"),
// free:                  ts.mod.ExportedFunction("free"),
// strlen:                ts.mod.ExportedFunction("strlen"),
// parserNew:             ts.mod.ExportedFunction("ts_parser_new"),
// parserParseString:     ts.mod.ExportedFunction("ts_parser_parse_string"),
// parserSetLanguage:     ts.mod.ExportedFunction("ts_parser_set_language"),
// parserDelete:          ts.mod.ExportedFunction("ts_parser_delete"),
// queryNew:              ts.mod.ExportedFunction("ts_query_new"),
// queryCursorNew:        ts.mod.ExportedFunction("ts_query_cursor_new"),
// queryCusorExec:        ts.mod.ExportedFunction("ts_query_cursor_exec"),
// queryCursorNextMatch:  ts.mod.ExportedFunction("ts_query_cursor_next_match"),
// queryCaptureNameForID: ts.mod.ExportedFunction("ts_query_capture_name_for_id"),
// //	languageName:          ts.mod.ExportedFunction("ts_language_name"),
// languageVersion:     ts.mod.ExportedFunction("ts_language_version"),
// treeRootNode:        ts.mod.ExportedFunction("ts_tree_root_node"),
// nodeString:          ts.mod.ExportedFunction("ts_node_string"),
// nodeChildCount:      ts.mod.ExportedFunction("ts_node_child_count"),
// nodeNamedChildCount: ts.mod.ExportedFunction("ts_node_named_child_count"),
// nodeChild:           ts.mod.ExportedFunction("ts_node_child"),
// nodeNamedChild:      ts.mod.ExportedFunction("ts_node_named_child"),
// nodeType:            ts.mod.ExportedFunction("ts_node_type"),
// nodeStartByte:       ts.mod.ExportedFunction("ts_node_start_byte"),
// nodeEndByte:         ts.mod.ExportedFunction("ts_node_end_byte"),
// nodeIsError:         ts.mod.ExportedFunction("ts_node_is_error"),
const (
	_malloc                = "malloc"
	_free                  = "free"
	_strlen                = "strlen"
	_parserNew             = "ts_parser_new"
	_parserParseString     = "ts_parser_parse_string"
	_parserDelete          = "ts_parser_delete"
	_parserSetLanguage     = "ts_parser_set_language"
	_languageVersion       = "ts_language_version"
	_treeRootNode          = "ts_tree_root_node"
	_queryNew              = "ts_query_new"
	_queryCursorNew        = "ts_query_cursor_new"
	_queryCursorExec       = "ts_query_cursor_exec"
	_queryCursorNextMatch  = "ts_query_cursor_next_match"
	_queryCaptureNameForID = "ts_query_capture_name_for_id"
	_nodeString            = "ts_node_string"
	_nodeChildCount        = "ts_node_child_count"
	_nodeNamedChildCount   = "ts_node_named_child_count"
	_nodeChild             = "ts_node_child"
	_nodeNamedChild        = "ts_node_named_child"
	_nodeType              = "ts_node_type"
	_nodeStartByte         = "ts_node_start_byte"
	_nodeEndByte           = "ts_node_end_byte"
	_nodeIsError           = "ts_node_is_error"
)

var _languages = []string{
	"bash",
	"c",
	"cpp",
	"c_sharp",
	"css",
	"cue",
	"dockerfile",
	"elixir",
	"elm",
	"go",
	"groovy",
	"hcl",
	"html",
	"java",
	"javascript",
	"kotlin",
	"lua",
}
var _functions = [23]string{
	_malloc,
	_free,
	_strlen,
	_parserNew,
	_parserParseString,
	_parserDelete,
	_parserSetLanguage,
	_languageVersion,
	_treeRootNode,
	_queryNew,
	_queryCursorNew,
	_queryCursorExec,
	_queryCursorNextMatch,
	_queryCaptureNameForID,
	_nodeString,
	_nodeChildCount,
	_nodeNamedChildCount,
	_nodeChild,
	_nodeNamedChild,
	_nodeType,
	_nodeEndByte,
	_nodeStartByte,
	_nodeIsError,
}
