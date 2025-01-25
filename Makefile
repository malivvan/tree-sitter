.PHONY: build
default: build

check-update:
	@cd src && go run ./_gen check-updates

update-all:
	@cd src && go run ./_gen update-all

build:
	@mkdir -p lib
	@rm -f lib/ts.wasm
	@zig cc --target=wasm32-wasi-musl -mexec-model=reactor -I src/ src/lib.c \
			src/c/parser.c \
			src/cpp/parser.c src/cpp/scanner.c \
			-o lib/ts.wasm -Os -fPIC -Wl,--no-entry -Wl,-z -Wl,stack-size=65536 -Wl,--strip-debug \
			-Wl,--import-symbols \
			-Wl,--export=malloc \
			-Wl,--export=free \
			-Wl,--export=strlen \
			-Wl,--export=ts_parser_new \
			-Wl,--export=ts_parser_parse_string \
			-Wl,--export=ts_parser_set_language \
			-Wl,--export=ts_parser_delete \
			-Wl,--export=ts_query_new \
			-Wl,--export=ts_query_cursor_new \
			-Wl,--export=ts_query_cursor_exec \
			-Wl,--export=ts_query_cursor_next_match \
			-Wl,--export=ts_query_capture_name_for_id \
			-Wl,--export=ts_language_name \
			-Wl,--export=ts_language_version \
			-Wl,--export=ts_tree_root_node \
			-Wl,--export=ts_node_string \
			-Wl,--export=ts_node_child_count \
			-Wl,--export=ts_node_named_child_count \
			-Wl,--export=ts_node_child \
			-Wl,--export=ts_node_named_child \
			-Wl,--export=ts_node_type \
			-Wl,--export=ts_node_start_byte \
			-Wl,--export=ts_node_end_byte \
			-Wl,--export=ts_node_is_error \
			-Wl,--export=tree_sitter_cpp \
			-Wl,--export=tree_sitter_c
	@go run ./lib/_gen




