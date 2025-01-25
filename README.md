
# tree-sitter [![Go Reference](https://pkg.go.dev/badge/github.com/malivvan/tree-sitter.svg)](https://pkg.go.dev/github.com/malivvan/tree-sitter) [![Release](https://img.shields.io/github/v/release/malivvan/tree-sitter.svg?sort=semver)](https://github.com/malivvan/tree-sitter/releases/latest) [![Go Report Card](https://goreportcard.com/badge/github.com/malivvan/tree-sitter)](https://goreportcard.com/report/github.com/malivvan/tree-sitter) [![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
Go module github.com/malivvan/tree-sitter is a cgo-free tree-sitter wrapper.<br>
It wraps a Wasm build of [tree-sitter](https://github.com/tree-sitter/tree-sitter), and uses [wazero](https://github.com/tetratelabs/wazero) as the runtime.<br>
Go and wazero are the _only_ direct dependencies.

> This is pre release software so expect bugs and potentially API breaking changes
> but each release will be tagged to avoid breaking people's code.

## Installation

```sh
# go 1.23.4+
go get github.com/malivvan/tree-sitter@latest
```


## Credits
- [github.com/ngavinsir/treesittergo](https://github.com/ngavinsir/treesittergo) for the general idea
- [github.com/smacker/go-tree-sitter](https://github.com/smacker/go-tree-sitter) for all the grammars and tooling

## License
This project is licensed under the [MIT License](LICENSE).