package main

import (
	"fmt"
	"os"

	"github.com/ipld/go-ipld-prime"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	wasmadl "github.com/willscott/wasm-adl/host"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: example rot13.wasm string\n")
		os.Exit(1)
	}

	builder := wasmadl.New(os.Args[1])
	baseStr := basicnode.NewString(os.Args[2])

	reified, err := builder.Reify(ipld.LinkContext{}, baseStr, nil)
	if err != nil {
		panic(err)
	}
	outStr, err := reified.AsString()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Reified: %s\n", outStr)
}
