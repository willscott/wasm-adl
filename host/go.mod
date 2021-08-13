module github.com/willscott/wasm-adl/host

go 1.16

require github.com/ipld/go-ipld-prime v0.11.0

require (
	github.com/ipfs/go-cid v0.0.4
	github.com/willscott/gasm v0.1.0
	github.com/willscott/wasm-adl v0.1.0
)

replace github.com/willscott/gasm => ../../gasm

replace github.com/willscott/wasm-adl => ../
