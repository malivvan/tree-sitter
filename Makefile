.PHONY: build
default: build

check-update:
	@cd src && go run ./_gen check-updates

update-all:
	@cd src && go run ./_gen update-all

update-%:
	@cd src && go run ./_gen update $*

force-update-%:
	@cd src && go run ./_gen force-update $*

build:
	@mkdir -p lib
	@rm -f lib/ts.wasm
	@zig cc --target=wasm32-wasi-musl -mexec-model=reactor -I src/ src/lib.c \
			src/bash/parser.c src/bash/scanner.c \
			src/c/parser.c \
			src/cpp/parser.c src/cpp/scanner.c \
			src/csharp/parser.c src/csharp/scanner.c \
			src/css/parser.c src/css/scanner.c \
			src/cue/parser.c src/cue/scanner.c \
			src/dockerfile/parser.c src/dockerfile/scanner.c \
			src/elixir/parser.c src/elixir/scanner.c \
			src/elm/parser.c src/elm/scanner.c \
			src/golang/parser.c \
			src/groovy/parser.c src/groovy/scanner.c \
			src/hcl/parser.c src/hcl/scanner.c \
			src/html/parser.c src/html/scanner.c \
			src/java/parser.c \
			src/javascript/parser.c src/javascript/scanner.c \
			src/kotlin/parser.c src/kotlin/scanner.c \
			src/lua/parser.c src/lua/scanner.c \
			src/python/parser.c src/python/scanner.c \
			src/ruby/parser.c src/ruby/scanner.c \
			src/rust/parser.c src/rust/scanner.c \
			src/sql/parser.c src/sql/scanner.c \
			src/php/parser.c src/php/scanner.c \
			-o lib/ts.wasm -Oz -fPIC -Wl,--no-entry -Wl,-z -Wl,stack-size=65536 -Wl,--strip-debug \
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
			-Wl,--export=tree_sitter_bash \
			-Wl,--export=tree_sitter_c \
			-Wl,--export=tree_sitter_cpp \
			-Wl,--export=tree_sitter_c_sharp \
			-Wl,--export=tree_sitter_css \
			-Wl,--export=tree_sitter_cue \
			-Wl,--export=tree_sitter_dockerfile \
			-Wl,--export=tree_sitter_elixir \
			-Wl,--export=tree_sitter_elm \
			-Wl,--export=tree_sitter_go \
			-Wl,--export=tree_sitter_groovy \
			-Wl,--export=tree_sitter_hcl \
			-Wl,--export=tree_sitter_html \
			-Wl,--export=tree_sitter_java \
			-Wl,--export=tree_sitter_javascript \
			-Wl,--export=tree_sitter_kotlin \
			-Wl,--export=tree_sitter_lua \
			-Wl,--export=tree_sitter_python \
			-Wl,--export=tree_sitter_ruby \
			-Wl,--export=tree_sitter_rust \
			-Wl,--export=tree_sitter_sql \
			-Wl,--export=tree_sitter_php
	@go run ./lib/_gen

pack:
	@du -h lib/ts.wasm
	@brotli --best lib/ts.wasm -o lib/ts.wasm.br && rm lib/ts.wasm
	@du -h lib/ts.wasm.br

